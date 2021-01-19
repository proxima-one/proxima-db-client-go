
package proxima_db_client_go


import (
  "context"
  proxima_client "github.com/proxima-one/proxima-db-client-go/client"
  grpc "google.golang.org/grpc"
  "io/ioutil"
  //"fmt"
)

var DefaultDatabaseConfig = make(map[string]interface{}, 0)

func DefaultProximaServiceClient(dbIP, dbPort string) (proxima_client.ProximaServiceClient, error)   {
  address := dbIP + ":" + dbPort
  maxMsgSize := 1024*1024*1024
  conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithDefaultCallOptions(
      grpc.MaxCallRecvMsgSize(maxMsgSize),
      grpc.MaxCallSendMsgSize(maxMsgSize)))
  if err != nil {
    return nil, err
  }
  return proxima_client.NewProximaServiceClient(conn), nil
}

func  (db *ProximaDatabase) NewDefaultTable(name, id string) (*ProximaDatabase, error) {
  return NewProximaTable(db, name, id, db.sleep), nil
}

func NewDefaultDatabase(name, id string) (*ProximaDatabase, error) {
  client := make(*ProximaServiceClient)
  clients := make([]interface{})
  sleepInterval := time.ParseDuration()
  compressionInterval := time.ParseDuration()
  batchingInterval := time.ParseDuration()

  return NewProximaDatabase(name, id, "0.0.0.0", client, clients, sleepInterval,
    compressionInterval, batchingInterval), nil
}

func NewProximaDatabase(name, version, id string, client *ProximaServiceClient, clients []interface{}, sleepInterval time.Duration,
  compressionInterval time.Duration,
  batchingInterval time.Duration) (*ProximaDatabase, error) {

  db := &ProximaDatabase{name: name, id: id, version: version, client: client, clients: clients, tables: nil, sleep: sleepInterval, compression: compressionInterval, batching: batchingInterval}
  return db, nil
}

func CheckLatest(checkType string, config map[string]map[string]) (map[string]interface{}, error) {
  currentType := ""
  currentVersion := ""
  currentName := ""
  returnValue := make(map[string]interface{})
  for newType, newConfig := range config {
    newVersion := newConfig[checkType]
    if currentVersion < newVersion || currentName == "" {
      currentType = newType
      currentVersion = newVersion
      currentName = newConfig["name"]
    }
  }
  returnValue["type"] = currentType.(interface{})
  returnValue["config"] = config.(interface{})
  return returnValue
}

func (db *ProximaDatabase) UpdateClients(newClients []interface{}) {
  extend(db.clients, newClients)
}

func GetClients(config map[string]interface{}) ([]interface{}, error) {
  clients := config["clients"].([]interface{})
  return clients, nil
}

func GetClient(clients []interface{}) (proxima_client.ProximaServiceClient, error) {
  ip := getEnv("DB_ADDRESS" , "0.0.0.0")
  port :=  getEnv("DB_PORT", "50051")
  client, err := DefaultProximaServiceClient(ip, port)
  while err != nil && clients.length > 0 {
    i := math.randomInt(clients.length)
    clientConfig := clients.pop(i)
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
  clients, err := db.GetClients(config)
  client, clientErr := db.GetClient(clients)

  db := NewProximaDatabase(config["name"].(string), config["id"].(string), client, clients, time.ParseDuration(config["sleep"].(string)), time.ParseDuration(config["compression"].(string)), time.ParseDuration(config["batching"].(string)))

  for name, tableConfig := range tables {
        table, err := db.LoadTable(config)
        //db.tables[name] = table
  }
  db.Update()
  return db, nil
}

type ProximaDatabase struct {
  client proxima_client.ProximaServiceClient
  name string
  id string
  tables map[string]*ProximaTable
  version string
  clients []interface{}
  sleep time.Duration
  compression time.Duration
  batching time.Duration
}

func (db *ProximaDatabase) PushNetworkConfig(method string) (bool, error){
  return true, nil
}

func (db *ProximaDatabase) PullNetworkConfig(method string) (map[string]interface{}, error) {
  return make(map[string]interface{}), nil
}

func (db *ProximaDatabase) GetNetworkDatabaseConfig(method string) (map[string]interface{}, error) {
  return make(map[string]interface{}), nil
}

func (db *ProximaDatabase) GetAllDatabaseConfig(methodType string) (map[string]map[string]interface{}, error) {
  config := make(map[string]interface{})
  config["local"] = db.GetLocalConfig()
  config["current"] = db.GetCurrentConfig()
  config["node"] = db.GetNetworkConfig("node")
  if methodType == "global" {
    config["network"] = db.GetNetworkConfig("global")
  }
  return config, nil
}

func (db *ProximaDatabase) GetLatestDatabaseConfig(methodType string) (map[string]map[string]interface{}, error) {
  config, err := GetAllDatabaseConfig(methodType)
  return GetMaxConfig("type", "version", config)
}

func (db *ProximaDatabase) Sync() (bool, error) {
  db.Update()
  config, err := db.GetLatestDatabaseConfig("global")
  syncType := config["type"].(string)
  syncConfig := config[syncType].(map[string]interface{})
  db.SetCurrentTableConfig(syncConfig)

  for _, tableConfig := range newTables {
    table, _ := db.LoadTable(syncType, tableConfig.(map[string]interface{}))
    go table.Update()
    go table.Sync(syncType, tableConfig.(map[string]interface{}))
  }
  db.Update()
  return true, nil
}

func (db *ProximaDatabase) Update() (bool, error) {
  newConfig, err := db.GetLatestDatabaseConfig("node")
  syncType := config["type"].(string)
  syncConfig := config[syncType].(map[string]interface{})
  for name, table := range db.tables {
    if !newTables.Contains(name) {
      go db.RemoveTable(name)
    } else {
      go table.Update()
    }
  }
  db.SetCurrentTableConfig(syncConfig)
  for name, tableConfig := range newTables {
    go db.LoadTable(syncType, tableConfig)
  }
  db.SetLocalDatabaseConfig();
  db.PushNetworkDatabaseConfig("node"); //Local Node
  db.PushNetworkDatabase("node");
  return true, nil
}

func (db *ProximaDatabase) GetCurrentDatabaseConfig() (map[string]interface{}, error) {
  var dbConfig map[string]interface{}
  tables :=  make([]interface{})
  dbConfig["name"] = db.name
  dbConfig["id"] = db.id
  dbConfig["sleepInterval"] = db.sleepInterval.String()
  dbConfig["compressionInterval"] = db.compressionInterval.String()
  dbConfig["batchingInterval"] = db.batchingInterval.String()
      for name, table := range db.tables {
          c, err := table.GetCurrentConfig();
          if err != nil {
            return nil, err
          }
          append(tables, c)
      }
  dbConfig["tables"] = tables
  return dbConfig, nil
}

func (db *ProximaDatabase) GetLocalDatabaseConfig() (map[string]interface{}, error) {
  return db.Get(db.name, name)
}

func (db *ProximaDatabase) SetLocalDatabaseConfig() (map[string]interface{}, error) {
    result, err := db.Set(db.id, "config", db.GetCurrentConfig(), nil)
    if err != nil {
      return false, err
    } else {
    return true, nil
  }
}

func (db *ProximaDatabase) LoadTable(loadType, tableConfig map[string]interface{}) (*ProximaTable, error) {
  table, err := db.GetTable(tableConfig["name"].(string))
  if err != nil {
    return nil, nil
  }
  if table == nil {
    table = db.NewDefaultTable(tableConfig)
  }
  table.Load(loadType, tableConfig)
  db.tables[name] = table
  return table, nil
}

func (db *ProximaDatabase) GetTable(name string)  (*ProximaTable, error) {
  if (db.tables.Contains(name) {
    return db.tables[name], nil
  }
  return nil, nil
}

func (db *ProximaDatabase) AddTable(name string, table *ProximaTable) {
  db.tables[name] = table
}

func (db *ProximaDatabase) RemoveTable(name string) {
  delete(db.tables, name)
}

func (db *ProximaDatabase) Delete() (bool, error) {
    resp, err := db.Remove(db.id, "config", nil)
    if err != nil {
      return false, err
    }  else {
      for name, table := range db.tables {
          table.Delete();
      }
      return true, nil
    }
}

func (db *ProximaDatabase) Open() (bool, error) {
    resp, err := db.client.Open(context.TODO(), &proxima_client.OpenRequest{Name: db.id})
    if err != nil {
      return false, err
    } else {
      for name, table := range db.tables {
          table.open();
      }
      return resp.GetConfirmation(), nil
    }
}

func (db *ProximaDatabase) Close() (bool, error) {
    resp, err := db.client.Close(context.TODO(), &proxima_client.CloseRequest{Name: db.id})
    if err != nil {
      return false, err
    } else {
      for name, table := range db.tables {
          table.close()
      }
      return resp.GetConfirmation(), nil
    }
}
