package database


import (
  grpc "google.golang.org/grpc"
  	context "context"
  _ "io/ioutil"
  "time"
  client "github.com/proxima-one/proxima-db-client-go/pkg/client"
  _ "fmt"
  "math/rand"
  "encoding/json"
  "os"
)

var DefaultDatabaseConfig = make(map[string]interface{})



func  (db *ProximaDatabase) NewDefaultTable(name string) (*ProximaTable, error) {
  return NewProximaTable(db, name, db.id, db.sleep), nil
}

func NewProximaDatabase(name, id, version string, client client.ProximaServiceClient, clients []interface{}, sleepInterval time.Duration,
  compressionInterval time.Duration,
  batchingInterval time.Duration) (*ProximaDatabase, error) {

  db := &ProximaDatabase{name: name, id: id, version: version, client: client, clients: clients, tables: make(map[string]*ProximaTable), sleep: sleepInterval, compression: compressionInterval, batching: batchingInterval}
  return db, nil
}

func CheckLatest(checkType string, config map[string]interface{}) (map[string]interface{}, error) {
  currentType := "."
  currentVersion := "0.0.0"
  //currentName := "."
  returnValue := make(map[string]interface{})
//fmt.Println(config)
  for newType, nConfig := range config {
    if nConfig == nil {
      continue
    }
    var newConfig map[string]interface{} = nConfig.(map[string]interface{})
    if newConfig[checkType] == nil  {
      continue
    }
    newVersion := newConfig[checkType].(string)
    if currentVersion <= newVersion {
      currentType = newType
      currentVersion = newVersion
      //currentName = newConfig["name"].(string)
    }
  }
  returnValue["type"] = currentType
  returnValue["config"] = config
  return returnValue, nil
}

func (db *ProximaDatabase) UpdateClients(newClients []interface{}) {
  db.clients = append(db.clients, newClients)
}

func GetClients(config map[string]interface{}) ([]interface{}, error) {
  if config["client"] == nil {
    return make([]interface{}, 0), nil
  }
  clients := config["clients"].([]interface{})
  return clients, nil
}

func getEnv(key, defaultVal string) (string) {
  val := os.Getenv(key)
  if val != "" {
    return val
  }
  return defaultVal
}

func DefaultProximaServiceClient(dbIP, dbPort string) (client.ProximaServiceClient, error)   {
  address := dbIP + ":" + dbPort
  maxMsgSize := 1024*1024*1024
  conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithDefaultCallOptions(
      grpc.MaxCallRecvMsgSize(maxMsgSize),
      grpc.MaxCallSendMsgSize(maxMsgSize)))
  if err != nil {
    return nil, err
  }
  return client.NewProximaServiceClient(conn), nil
}

func GetClient(clients []interface{}) (client.ProximaServiceClient, error) {
  ip := getEnv("DB_ADDRESS" , "0.0.0.0")
  port :=  getEnv("DB_PORT", "50051")
  client, err := DefaultProximaServiceClient(ip, port)
  r := rand.New(rand.NewSource(99))
  for (err != nil) && (len(clients) > 0) {
    i := r.Intn(len(clients))
    clientConfig := clients[i].(map[string]interface{})
    clients = append(clients[:i], clients[i+1:]...)
    port = clientConfig["port"].(string)
    ip = clientConfig["ip"].(string)
    client, err = DefaultProximaServiceClient(ip, port)
  }
  if err != nil {
    return nil, err
  }
  return client, nil
}

func LoadProximaDatabase(config map[string]interface{}) (*ProximaDatabase, error) {
  clients, _ := GetClients(config)
  client, _ := GetClient(clients)
  name := config["name"].(string)
  id := config["id"].(string)

  var intervalConfig map[string]interface{} = config["config"].(map[string]interface{})
  //check config
  sleepStr := intervalConfig["sleep"].(string)
  compressionStr := intervalConfig["compression"].(string)
  batchingStr := intervalConfig["batching"].(string)
  sleep, _ := time.ParseDuration(sleepStr)
  compression, _ := time.ParseDuration(compressionStr)
  batching, _ := time.ParseDuration(batchingStr)

  version := config["version"].(string)

  db, _ := NewProximaDatabase(name, id,  version, client, clients, sleep, compression, batching)
  if config["tables"] != nil {
    var tables []interface{} = config["tables"].([]interface{})
    for _, tableConfig := range tables {
          var loadConfig map[string]interface{} = tableConfig.(map[string]interface{})
          db.LoadTable("local", loadConfig)
    }
  }
  db.Update()
  return db, nil
}

type ProximaDatabase struct {
  client client.ProximaServiceClient
  name string
  id string
  tables map[string]*ProximaTable
  version string
  clients []interface{}
  sleep time.Duration
  compression time.Duration
  batching time.Duration
}

func (db *ProximaDatabase) PushNetworkDatabaseConfig(method string) (bool, error){
  return true, nil
}

func (db *ProximaDatabase) PushNetworkDatabase(method string) (bool, error){
  return true, nil
}

func (db *ProximaDatabase) PullNetworkDatabaseConfig(method string) (map[string]interface{}, error) {
  return make(map[string]interface{}), nil
}

func (db *ProximaDatabase) GetNetworkDatabaseConfig(method string) (map[string]interface{}, error) {
  return make(map[string]interface{}), nil
}

func (db *ProximaDatabase) GetAllDatabaseConfig(methodType string) (map[string]interface{}, error) {
  config := make(map[string]interface{})
  config["local"], _ = db.GetLocalDatabaseConfig()
  config["current"], _ = db.GetCurrentDatabaseConfig()
  config["node"], _ = db.GetNetworkDatabaseConfig("node")
  if methodType == "global" {
    config["network"], _ = db.GetNetworkDatabaseConfig("global")
  }
  return config, nil
}

func (db *ProximaDatabase) GetLatestDatabaseConfig(methodType string) (map[string]interface{}, error) {
  config, _ := db.GetAllDatabaseConfig(methodType)
  return CheckLatest("version", config)
}

func (db *ProximaDatabase) Sync() (bool, error) {
  db.Update()
  config, _ := db.GetLatestDatabaseConfig("global")
  syncType := config["type"].(string)
  syncConfig := config[syncType].(map[string]interface{})
  db.SetCurrentDatabaseConfig(syncConfig, true)
  //newTables

  for _, table := range db.tables {
    go table.Sync(syncType, config)
  }
  db.Update()
  return true, nil
}

func (db *ProximaDatabase) Update() (bool, error) {
  newConfig, _ := db.GetLatestDatabaseConfig("node")
  syncType := newConfig["type"]
  if syncType != nil && newConfig[syncType.(string)] != nil {
    syncConfig := newConfig[syncType.(string)].(map[string]interface{})
    db.SetCurrentDatabaseConfig(syncConfig, true)
    db.SetLocalDatabaseConfig();
    db.PushNetworkDatabaseConfig("node");
    db.PushNetworkDatabase("node");
    return true, nil
  }
  return false, nil
}

func (db *ProximaDatabase) SetCurrentDatabaseConfig(newConfig map[string]interface{}, includeTables bool) (bool, error) {
  //fmt.Println(newConfig)
  db.name = newConfig["name"].(string)
  db.id = newConfig["id"].(string)
  db.version = newConfig["version"].(string)
  intervalConfig := newConfig["config"].(map[string]interface{})
  //check config
  sleep := intervalConfig["sleep"].(string)
  compression := intervalConfig["compression"].(string)
  batching := intervalConfig["batching"].(string)
  db.sleep, _ = time.ParseDuration(sleep)
  db.compression, _ = time.ParseDuration(compression)
  db.batching, _ = time.ParseDuration(batching)

  if includeTables {
  newTables := newConfig["tables"].([]interface{})
  newTablesMap := make(map[string]string)
   for _, t := range newTables {
     tableConfig := t.(map[string]interface{})
     if tableConfig["name"] == nil {
       continue
     }
     name := tableConfig["name"].(string)
     newTablesMap[name] = name
      go db.LoadTable("", tableConfig)
    }

    for _, table := range db.tables {
      name := table.name
      if newTablesMap[name] == name {
        go db.RemoveTable(name)
      }
    }
  }
  return true, nil
}

func (db *ProximaDatabase) GetCurrentDatabaseConfig() (map[string]interface{}, error) {
  var dbConfig map[string]interface{} = make(map[string]interface{})
  tables :=  make([]interface{}, 0)
  dbConfig["name"] = db.name
  dbConfig["id"] = db.id
  dbConfig["version"] = db.version
  intervalConfig := make(map[string]interface{})
  intervalConfig["sleep"] = db.sleep.String()
  intervalConfig["compression"] = db.compression.String()
  intervalConfig["batching"] = db.batching.String()
  dbConfig["config"] = intervalConfig
      for _, table := range db.tables {
          c, err := table.GetCurrentTableConfig();
          if err != nil {
            tables = append(tables, c)
          }

  }
  dbConfig["tables"] = tables
  return dbConfig, nil
}

func (db *ProximaDatabase) GetLocalDatabaseConfig() (map[string]interface{}, error) {
  resp, err := db.Get(db.id, "config", nil)
  config := make(map[string]interface{})
  if err != nil || resp == nil  {
    return config, err
  }
  json.Unmarshal(resp.GetValue(), &config)
  return config, nil
}

func (db *ProximaDatabase) SetLocalDatabaseConfig() (map[string]interface{}, error) {
    currentConfig, _ := db.GetCurrentDatabaseConfig()
    _, err := db.Set(db.id, "config", currentConfig, nil)
    if err != nil {
      return nil, err
    } else {
    return currentConfig, nil
  }
}

func (db *ProximaDatabase) LoadTable(loadType string, tableConfig map[string]interface{}) (*ProximaTable, error) {
  tableName := tableConfig["name"].(string)
  table, err := db.GetTable(tableName)
  if err != nil {
    return nil, nil
  }
  if table == nil {
    table, err = db.NewDefaultTable(tableName)
  }
  table.Load(loadType, tableConfig)
  db.tables[tableName] = table
  return table, nil
}

func (db *ProximaDatabase) GetTable(name string)  (*ProximaTable, error) {
  if db.tables[name] != nil {
    return db.tables[name], nil
  }
  return nil, nil
}

func (db *ProximaDatabase) AddTable(name string, table *ProximaTable) {
  db.tables[name] = table
}

func (db *ProximaDatabase) RemoveTable(name string) (bool, error) {
  delete(db.tables, name)
  return true, nil
}

func (db *ProximaDatabase) Delete() (bool, error) {
    _, err := db.Remove(db.id, "config", nil)
    if err != nil {
      return false, err
    }  else {
      for _, table := range db.tables {
          table.Delete();
      }
      return true, nil
    }
}

func (db *ProximaDatabase) Open() (bool, error) {
    resp, err := db.client.Open(context.TODO(), &client.OpenRequest{Name: db.id})
    if err != nil {
      return false, err
    } else {
      for _, table := range db.tables {
          table.Open();
      }
      return resp.GetConfirmation(), nil
    }
}

func (db *ProximaDatabase) Close() (bool, error) {
    resp, err := db.client.Close(context.TODO(), &client.CloseRequest{Name: db.id})
    if err != nil {
      return false, err
    } else {
      for _, table := range db.tables {
          table.Close()
      }
      return resp.GetConfirmation(), nil
    }
}
