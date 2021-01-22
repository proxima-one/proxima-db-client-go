package database
//package table

import (
  json "github.com/json-iterator/go"
  //proxima "github.com/proxima-one/proxima-db-client-go
  //client "github.com/proxima-one/proxima-db-client-go/client"
  "time"
  "fmt"
)

func NewProximaTable(db *ProximaDatabase, name, id string, cacheExpiration time.Duration) (*ProximaTable) {
  table :=  &ProximaTable{db: db, name: name, id: id, cache: NewTableCache(cacheExpiration), isOpen: false, isIdle: false, sleepInterval: db.sleepInterval, compressionInterval: db.compressionInterval, batchingInterval: db.batchingInterval, header: "Root", blockNum: 0}
  return table
}

type ProximaTable struct {
  name *string
  id *string
  version *string
  blockNum int
  header *string
  isOpen bool
  isIdle bool
  sleepInterval time.Duration
  compressionInterval time.Duration
  batchingInterval time.Duration

  db *ProximaDatabase
  cache *ProximaTableCache
}

func (table *ProximaTable) GetLatestTableConfig(methodType string) (map[string]map[string]interface{}, error) {
  config := make(map[string]interface{})
  config["node"] = table.GetNetworkTableConfig("node")
  config["local"] = table.GetLocalTableConfig()
  config["current"] = table.GetCurrentTableConfig()
  return GetLatestConfig("type", "blockNum", config)
}

func (table *ProximaTable) GetAllTableConfig(methodType string) (map[string]map[string]interface{}, error) {
  config := make(map[string]interface{})
  config["local"] = table.GetLocalTableConfig()
  config["current"] = table.GetCurrentTableConfig()
  config["node"] = table.GetNetworkTableConfig("node")

  if methodType == "global" {
    config["network"] = table.GetNetworkTableConfig("global")
  }
  return config, nil
}

func (table *ProximaTable) GetNetworkTableConfig(methodType string) (map[string]interface{}, error) {
  return make(map[string]interface{}), nil
}

func (table *ProximaTable) GetLocalTableConfig() (map[string]interface{}, error) {
  result, err := db.Get(db.name, name)
  if err != nil {
    return nil, err
  } else {
    return result, nil
  }
}

func (table *ProximaTable) SetLocalTableConfig() (bool, error) {
  resp, err := table.db.Set(table.id, table.name, config, nil)
  if err != nil {
    return false, err
  }
  return resp.GetConfirmation(), nil
}

func (table *ProximaTable) GetCurrentTableConfig() (map[string]interface{}, error) {
  config := make(map[string]interface{});
  config["name"] = table.name
  config["id"] = table.id
  config["version"] = table.version
  config["blockNum"] = table.blockNum
  config["header"] = table.header
  config["sleepInterval"] = table.sleepInterval.String()
  config["compressionInterval"] = table.compressionInterval.String()
  config["batchingInterval"] = table.batchingInterval.String()
  config["cacheExpiration"] = table.cache.cacheExpiration.String()
  return config, nil
}

func (table *ProximaTable) SetCurrentTableConfig(config map[string]interface{}) (bool, error) {
  table.name = config["name"].(string)
  table.id = config["id"].(string)
  table.version = config["version"].(string)
  table.blockNum = config["blockNum"].(int)
  table.header  = config["header"].(string)
  table.sleepInterval = time.ParseDuration(config["sleepInterval"].(string))
  table.compressionInterval = time.ParseDuration(config["compressionInterval"].(string))
  table.batchingInterval = time.ParseDuration(config["batchingInterval"].(string))
  table.cache = NewTableCache(time.ParseDuration(config["cacheExpiration"].(string)))
  return true, nil
}

func (table *ProximaTable) Sync(syncType string, syncConfig map[string]interface{}) (map[string]interface{}, error) {
  //tf
  newConfig, err := db.GetMaxExternalDatabaseConfig(syncConfig)
  if newConfig["type"].(string) == "network" {
      table.Load("global", syncConfig)
  }

  if newConfig["type"].(string) == "node" {
      table.Load("node", syncConfig)
  }
  table.Load("local", syncConfig)
  db.PushNetworkTableConfig("node");
  db.PushNetworkTable("node");
  return newConfig, nil
}

func (table *ProximaDB) Load(configType string, config map[string]interface{}) {
  table.Update()
  table.SetCurrentTableConfig(config)
  if configType == "global" {
      table.PullNetworkTable("global")
      table.PullNetworkTableConfig("global")
  }
  if configType == "node" {
      table.PullNetworkTable("node")
      table.PullNetworkTableConfig("node")
  }
}

func (table *ProximaTable) Update() (bool, error) {
  newConfig, err := db.GetMaxInternalDatabaseConfig()
  syncType := config["type"].(string)
  syncConfig := config[syncType].(map[string]interface{})

  table.SetCurrentTableConfig(syncConfig);
  db.SetLocalDatabaseConfig();
  return true, nil
}

func (table *ProximaTable) Delete() (bool, error) {
    table.Close();
    table.db.Delete(table.name);
    _ , err:= db.client.TableRemove(context.TODO(), &client.TableRemoveRequest{Name: table.id})
    if err != nil {
      return false, err
    }
    return true, nil
}

func (table *ProximaTable) Open() (error) {
  if table.isOpen {
    return nil
  }
  err := table.db.OpenTable(table.name, table);
  if err != nil {
    return err
  } else {
    table.isIdle = false;
    table.isOpen = true;
    go Compression(table, table.compressionInterval);
    go Batching(table, table.batchingInterval);
    go SleepSchedule(table, table.sleepInterval);
  }
  return nil
}

func Compression(table *ProximaTable, interval time.Duration) {
  ticker := time.NewTicker(interval)
  for ; true; <-ticker.C {
    select {
      case !table.isOpen:
        //ticker stop
                return
      case t := <-ticker.C:
          //compress the database ... table.Compress()
      }
  }
}

func Batching(table *ProximaTable, interval time.Duration) {
  ticker := time.NewTicker(interval)
  for ; true; <-ticker.C {
    select {
      case !table.isOpen: //is not open
      //compress
      //ticker stop
                return
      case t := <-ticker.C:
            //compress the transaction CheckoutTransaction with table
            //make a new transaction CheckIn
      }
  }
}

func SleepSchedule(table *ProximaTable, interval time.Duration) {
  ticker := time.NewTicker(interval)
  for ; true; <-ticker.C {
    select {
      case table.isIdle:
          //ticker stop
          ticker.Stop()
          table.Close()
          return
      case t := <-ticker.C:
          table.isIdle = true;
      }
  }
}

func (table *ProximaTable) Close() (error) {
  table.isIdle = false; //turns off of the sleep
  table.isOpen = false; //turns off compression and batching
  err := table.db.CloseTable(table.name);
  table.cache.cache.flush();
  if err != nil {
    return err
  }
  return nil
}

func (table *ProximaTable) Query(queryString string, prove bool) (*ProximaDBResult, error) {
  table.isIdle = false
  return table.db.Query(table.id, queryString, prove);
}

func (table *ProximaTable) Get(key string,  prove bool) (*ProximaDBResult, error) {
  var result *ProximaDBResult;
  table.isIdle = false
  if cached, found := table.cache.Get(key); found {
  result = cached
  } else {
  result, err := table.db.Get(table.id, key, true) //cache result
  if err != nil {
    return nil, err
  }
  table.cache.Set(key, result)
  }
  return result, nil
}
//fix
func (table *ProximaTable) Put(key string, value map[string]interface{}) (*ProximaDBResult, error) {
  var result *ProximaDBResult;
  table.isIdle = false
  result, err := table.db.Set(table.id, key, value);
  if err != nil {
    return nil, err
  }
  table.cache.Set(key, result);
  //update blockNum
  if value["blockNum"] != nil {
    table.blockNum = value["blockNum"].(int)
  }
  return result, err;
}

func (table *ProximaTable) Remove(key string) (*ProximaDBResult, error) {
  table.isIdle = false
  var result *ProximaDBResult;
  table.cache.Remove(key);
  return table.db.Remove(table.name, key);
}
