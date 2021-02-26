package database
//package table

import (
  //json "github.com/json-iterator/go"
  //proxima "github.com/proxima-one/proxima-db-client-go
  client "github.com/proxima-one/proxima-db-client-go/pkg/client"
  "context"
  "time"
  "strconv"
  "fmt"
  "encoding/json"
)

func NewProximaTable(db *ProximaDatabase, name, id string, cacheExpiration time.Duration) (*ProximaTable) {
  table :=  &ProximaTable{db: db, name: name, id: id, cache: NewTableCache(cacheExpiration), isOpen: false, isIdle: false, sleep: db.sleep, compression: db.compression, batching: db.batching, header: "Root", blockNum: 0}
  return table
}

type ProximaTable struct {
  name string
  id string
  version string
  blockNum int
  header string
  isOpen bool
  isIdle bool
  sleep time.Duration
  compression time.Duration
  batching time.Duration

  db *ProximaDatabase
  cache *ProximaTableCache
}

func (table *ProximaTable) GetLatestTableConfig(methodType string) (map[string]interface{}, error) {
  config := make(map[string]interface{})
  config["node"], _ = table.GetNetworkTableConfig("node")
  config["local"], _ = table.GetLocalTableConfig()
  config["current"], _ = table.GetCurrentTableConfig()
  return CheckLatest("version", config)
}

func (table *ProximaTable) GetAllTableConfig(methodType string) (map[string]interface{}, error) {
  config := make(map[string]interface{})
  config["local"], _ = table.GetLocalTableConfig()
  config["current"], _ = table.GetCurrentTableConfig()
  config["node"], _ = table.GetNetworkTableConfig("node")

  if methodType == "global" {
    config["network"], _ = table.GetNetworkTableConfig("global")
  }
  return config, nil
}

func (table *ProximaTable) GetNetworkTableConfig(methodType string) (map[string]interface{}, error) {
  return make(map[string]interface{}), nil
}

func (table *ProximaTable) GetLocalTableConfig() (map[string]interface{}, error) {
  setConfig := make(map[string]interface{});
  setConfig["prove"] = false
  resp, err := table.db.Get(table.id, table.name, setConfig)
  if err != nil || resp == nil {
    return make(map[string]interface{}), err
  } else {
    config := make(map[string]interface{})
    json.Unmarshal(resp.GetValue(), &config)
    return config, nil
  }
}

func (table *ProximaTable) SetLocalTableConfig() (bool, error) {
  config, configErr := table.GetCurrentTableConfig()
  setConfig := make(map[string]interface{});
  setConfig["prove"] = false
  if configErr != nil {
    return false, configErr
  }
  _, err := table.db.Set(table.id, table.name, config, setConfig)
  if err != nil {
    return false, err
  }
  return true, nil
}

func (table *ProximaTable) GetCurrentTableConfig() (map[string]interface{}, error) {
  config := make(map[string]interface{});
  config["name"] = table.name
  config["id"] = table.id
  config["version"] = table.version
  config["blockNum"] = table.blockNum
  config["header"] = table.header
  config["sleep"] = table.sleep.String()
  config["compression"] = table.compression.String()
  config["batching"] = table.batching.String()
  config["cacheExpiration"] = table.cache.cacheExpiration.String()
  return config, nil
}

func (table *ProximaTable) SetCurrentTableConfig(config map[string]interface{}) (bool, error) {
  table.name = config["name"].(string)
  table.id = config["id"].(string)
  table.version = config["version"].(string)
  blockNum := fmt.Sprintf("%v", config["blockNum"])
  table.blockNum, _ = strconv.Atoi(blockNum)
  table.header  = config["header"].(string)
  cacheExpiration, err := time.ParseDuration(config["cacheExpiration"].(string))
  table.sleep, err = time.ParseDuration(config["sleep"].(string))
  table.compression, err = time.ParseDuration(config["compression"].(string))
  table.batching, err = time.ParseDuration(config["batching"].(string))
  if err != nil {
    return false, err
  }
  table.cache = NewTableCache(cacheExpiration)
  return true, nil
}

func (table *ProximaTable) Sync(syncType string, syncConfig map[string]interface{}) (map[string]interface{}, error) {
  //, syncConfig
  newConfig, _ := table.db.GetLatestDatabaseConfig(syncType)
  if newConfig["type"].(string) == "network" {
      table.Load("global", syncConfig)
  }

  if newConfig["type"].(string) == "node" {
      table.Load("node", syncConfig)
  }
  table.Load("local", syncConfig)
  table.PushNetworkTableConfig("node");
  table.PushNetworkTable("node");
  return newConfig, nil
}

func (table *ProximaTable) Load(configType string, config map[string]interface{}) {
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
  newConfig, _ := table.GetLatestTableConfig("local")
  syncType := newConfig["type"]
  config := newConfig["config"].(map[string]interface{})
  if syncType != nil && config[syncType.(string)] != nil {

    syncConfig := config[syncType.(string)].(map[string]interface{})
    table.SetCurrentTableConfig(syncConfig);
    table.SetLocalTableConfig();
    return true, nil
  }
  table.SetLocalTableConfig();
  return false, nil
}

func (table *ProximaTable) Delete() (bool, error) {
    table.Close();
    table.db.RemoveTable(table.name);
    _ , err:= table.db.client.TableRemove(context.TODO(), &client.TableRemoveRequest{Name: table.id})
    if err != nil {
      return false, err
    }
    return true, nil
}

func (table *ProximaTable) Open() (error) {
  if table.isOpen {
    return nil
  }
  //err := table.db.OpenTable(table.name, table);
  // if err != nil {
  //   return err
  // } else {
    table.isIdle = false;
    table.isOpen = true;
    go Compression(table, table.compression);
    go Batching(table, table.batching);
    go SleepSchedule(table, table.sleep);
  // }
  return nil
}

func Compression(table *ProximaTable, interval time.Duration) {
  // ticker := time.NewTicker(interval)
  // for ; true; <-ticker.C {
  //   select {
  //     case !table.isOpen:
  //       //ticker stop
  //       return
  //     case t := <-ticker.C:
  //         //compress the database ... table.Compress()
  //     }
  // }
  return
}

func Batching(table *ProximaTable, interval time.Duration) {
  // ticker := time.NewTicker(interval)
  // for ; true; <-ticker.C {
  //   select {
  //     case !table.isOpen: //is not open
  //     //compress
  //     //ticker stop
  //               return
  //     case t := <-ticker.C:
  //           //compress the transaction CheckoutTransaction with table
  //           //make a new transaction CheckIn
  //     }
  // }
  return
}

func SleepSchedule(table *ProximaTable, interval time.Duration) {
  // ticker := time.NewTicker(interval)
  // for ; true; <-ticker.C {
  //   if table.isIdle {
  //         //ticker stop
  //         ticker.Stop()
  //         table.Close()
  //         return
  //   }
  //   if t := <-ticker.C {
  //         table.isIdle = true;
  //   }
  // }
  return
}

func (table *ProximaTable) Close() (error) {
  table.isIdle = false; //turns off of the sleep
  table.isOpen = false; //turns off compression and batching
  // err := table.db.CloseTable(table.name);
  // table.cache.cache.flush();
  // if err != nil {
  //   return err
  // }
  return nil
}

func (table *ProximaTable) Query(queryString string, prove bool) ([]*ProximaDBResult, error) {
  table.isIdle = false

  return table.db.Query(table.id, queryString, map[string]interface{}{"prove": prove});
}

func (table *ProximaTable) Get(key string,  prove bool) (*ProximaDBResult, error) {
  var result *ProximaDBResult;
  var err error;
  table.isIdle = false
  if cached, found := table.cache.Get(key); found {
  result = cached
  } else {
  result, err = table.db.Get(table.id, key, map[string]interface{}{"prove": true}) //cache result
  if err != nil {
    return nil, err
  }
  if result != nil {
    table.cache.Set(key, map[string]interface{}{"prove": false})
  }
  }
  return result, nil
}
//fix
func (table *ProximaTable) Put(key interface{}, value interface{}, prove bool, args map[string]interface{}) (*ProximaDBResult, error) {
  var result *ProximaDBResult;
  table.isIdle = false
  result, err := table.db.Set(table.id, key, value, map[string]interface{}{"prove": prove});
  if err != nil {
    return nil, err
  }
  table.cache.Set(key, result);
  //update blockNum
  if args["blockNum"] != nil {
    table.blockNum = args["blockNum"].(int)
  }
  return result, err;
}

func (table *ProximaTable) Remove(key string, prove bool) (*ProximaDBResult, error) {
  table.isIdle = false
  var result *ProximaDBResult;
  var err error;
  //var result *ProximaDBResult;
  table.cache.Remove(key);
  result, err  = table.db.Remove(table.id, key, map[string]interface{}{"prove": prove});
  return result, err
}

func (table *ProximaTable) PushNetworkTableConfig(method string) (bool, error){
  return true, nil
}

func (table *ProximaTable) PushNetworkTable(method string) (bool, error){
  return true, nil
}

func (table *ProximaTable) PullNetworkTable(method string) (bool, error){
  return true, nil
}

func (table *ProximaTable) PullNetworkTableConfig(method string) (map[string]interface{}, error) {
  return make(map[string]interface{}), nil
}
