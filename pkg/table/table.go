package proxima_db_client_go

//this should be for the proxima-db-client


import (
  json "github.com/json-iterator/go"
  //proxima "github.com/proxima-one/proxima-db-client-go"
  proxima_client "github.com/proxima-one/proxima-db-client-go/client"
  "time"
  "fmt"
)

type ProximaTable struct {
  name *string
  dbId *string
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

func (db *ProximaDatabase) GetTable(name string)  (*ProximaTable, error) {
  result, err := db.Get(db.name, name)
  table, err := TableFromConfig(db, result)
  if err != nil {
    return nil, err
  } else {
    return table, nil
  }
}

func TableFromConfig(db *ProximaDatabase, tableConfig map[string]interface{}) (*ProximaTable, error) {
  return NewProximaTable(db, tableConfig["name"].(string), tableConfig["dbId"].(string), tableConfig["cacheExpiration"].(time.Duration))
}

func (table *ProximaTable) Config() (map[string]interface{}) {
  var config map[string]interface{};
  config["name"] = table.name
  config["cacheExpiration"] = table.cacheExpiration
  config["dbId"] = table.dbId

  return config
}

func (table *ProximaTable) Delete() (bool, error) {
    table.Close();
    table.db.Delete(table.name);
    _ , err:= db.client.TableRemove(context.TODO(), &proxima_client.TableRemoveRequest{Name: table.dbId})
    if err != nil {
      return false, err
    }
    return true, nil
}
//header
func NewProximaTable(db *proxima.ProximaDatabase, name, dbId string, cacheExpiration time.Duration) (*ProximaTable) {
  table :=  &ProximaTable{db: db, name: name, dbId: dbId, cache: NewTableCache(cacheExpiration), isOpen: false, isIdle: false, sleepInterval: db.sleepInterval, compressionInterval: db.compressionInterval, batchingInterval: db.batchingInterval, header: "Root", blockNum: 0}
  return table
}

func (table *ProximaTable) Save(name string) (bool, error) {
  resp, err := table.db.Set(table.dbId, table.name, table.Config(), nil)
  if err != nil {
    return false, err
  }
  return resp.GetConfirmation(), nil
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

func (table *ProximaTable) Load() (error) {
  //load
  return nil
}

func (table *ProximaTable) Sync() (error) {
  //sync
  return nil
}

func (table *ProximaTable) Update() (bool, error) {
  //update the headers (table.getHeader)
  return true, nil
}

func (table *ProximaTable) Compare() (bool, error) {
  //compare
  return true, nil
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
  return table.db.Query(table.dbId, queryString, prove);
}

func (table *ProximaTable) Get(key string,  prove bool) (*ProximaDBResult, error) {
  var result *ProximaDBResult;
  table.isIdle = false
  if cached, found := table.cache.Get(key); found {
  result = cached
  } else {
  result, err := table.db.Get(table.dbId, key, true) //cache result
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
  result, err := table.db.Set(table.dbId, key, value);
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
