
package proxima_db_client_go


import (
  "context"
  proxima_client "github.com/proxima-one/proxima-db-client-go/client"
  grpc "google.golang.org/grpc"
  "io/ioutil"
  //"fmt"
)

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


type ProximaDatabase struct {
  client proxima_client.ProximaServiceClient
  name *string
  id *string
  tables map[string]*ProximaTable //map
  sleepInterval time.Duration //goroutine
  compressionInterval time.Duration //goroutine
  batchingInterval time.Duration //goroutine
}

func (proxima *proxima_client.ProximaServiceClient) LoadProximaDatabase(name string) (*ProximaDatabase, error) {
    var result *ProximaDatabase;
    resp, err := db.Get(db.id, "config", nil)
    if err != nil {
      return nil, err
    }
    db, dbErr := DBFromConfig(config)
    if dbErr != nil {
      return nil, err
    }
    for name, table := range db.tables {
          db.tables[name], _ = db.GetTable(name);
    }
      return db, nil
    }
}

func (db *ProximaDatabase) Save() (bool, error) {
    result, err := db.Set(db.id, "config", db.Config(), nil)
    if err != nil {
      return false, err
    } else {
      for name, table := range db.tables {
          table.Save();
      }
      return true, nil
    }
}

/*
TODO
Checks the yaml file, if the yaml file differs from the db file, then check the latest version, and load, saves the result
*/



func DBFromConfig(proxima *proxima_client.ProximaServiceClient, config map[string]interface{}) (*ProximaDatabase, error) {
  
  return NewProximaDatabase(config["name"].(string), config["id"].(string), proxima, nil, config["tables"].([]string), config["sleepInterval"].(time.Duration), config["compressionInterval"].(time.Duration), config["batchingInterval"].(time.Duration))
}

func (db *ProximaDatabase) LoadDatabaseConfig() (map[string]interface{}, error) {
  var dbConfig map[string]interface{}
  var tables []string;
  dbConfig["name"] = db.name
  dbConfig["id"] = db.id
  dbConfig["sleepInterval"] = db.sleepInterval
  dbConfig["compressionInterval"] = db.compressionInterval
  dbConfig["batchingInterval"] = db.batchingInterval
  for name, _ := range db.tables {
      tables = append(tables, name)
  }
  dbConfig["tables"] = tables
  return dbConfig, nil
}

//

func (db *ProximaDatabase) LoadExternalDatabaseConfig() ([]map[string]interface{}, error) {
  var dbConfig map[string]interface{}
  var tables []string;
  //
  return dbConfig, nil
}

func (db *ProximaDatabase) LoadLocalDatabaseConfig() (map[string]interface{}, error) {
  var dbConfig map[string]interface{}
  var tables []string;
  configFile, err := ioutil.ReadFile("db-config.yml")
    if err != nil {
        return nil, err
    }
    err = yaml.Unmarshal(configFile, dbConfig)
    if err != nil {
        return nil, err
    }
    return dbConfig, nil
}

//remove database
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

//wrong
func NewProximaDatabase(name, id string, client *ProximaServiceClient, tables map[string]*ProximaTable, tableList []string, sleepInterval time.Duration,
  compressionInterval time.Duration,
  batchingInterval time.Duration) (*ProximaDatabase, error) {
    //create table
  for i, name := range tableList {
        db.tables[name], _ = db.GetTable(name);
  }
  db := &ProximaDatabase{name: name, id: id, client: client, tables: tables, sleepInterval: sleepInterval, compressionInterval: compressionInterval, batchingInterval: batchingInterval}
  resp, err := db.Save();
  if err != nil {
    return nil, err
  } else {
    return db, nil
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

func (db *ProximaDatabase) CreateTable(name string, cacheExpiration time.Duration) (*ProximaTable, error) {
    table :=  NewProximaTable(name, db.id, cacheExpiration);
    db.addTable(table)
    return table, nil
}

func (db *ProximaDatabase) AddTable(name string, table *ProximaTable) {
  db.tables[name] = table
  table.Save()
}

func (db *ProximaDatabase) Delete(name string) {
  delete(db.tables, name)
}

func (db *ProximaDatabase) Update() (bool, error) {
  for _, table := range db.tables {
    table.Update()
  }
  return db.Save()
}

func (db *ProximaDatabase) Sync() (bool, error) {
  //db.Load(), new table
  //db.version
  //db.Compare() between two

  //compare app with
//compare the table version //time-based
//compare the table header
//select the correct table and update
//write the table to file
//write the table to application database
//if incorrect, get the correct data

  return db.Update()
}

//func load


// func (db *ProximaDatabase) Compare(dbA)  {
//
// }
