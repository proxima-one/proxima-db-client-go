package proxima_db_client_go

import (
  "testing"
  "math/rand"
  "fmt"
  "time"
)

func NewDatabaseClient() (*ProximaDB, error) {
  ip := "0.0.0.0"
  port := "50051"
  return DefaultProximaServiceClient(ip, port)
}



var proximaClient *ProximaDatabase = NewDatabaseClient();
var tables map[string][]string = TableSetup();
var appName string = "NewApplication";

var id string = "NewDBID";
var tables map[string]*ProximaTable;
var tableList []string = make([]string, 0);
var sleepInterval time.Duration = 5*time.Minute;
var cacheExpiration time.Duration = 5*time.Minute;
var compressionInterval time.Duration = 5*time.Minute;
var batchingInterval time.Duration = 5*time.Minute;


func NewTestApplicationDB() (*ProximaDastabase, error) {
  return NewProximaDatabase(appName, id, proximaClient, tables, tableList, sleepInterval, compressionInterval, batchingInterval)
}

//test Application Loading

//test database loading

//test table loading

//sync + update db, table

//remove tables

//delete

//batching

//compression

//sleeping



func TableSetup() (map[string][]string) {
  var num int= 20
  var removeLen int= (num-1)/2
  var partialLen int = remove + 1
  full := make([]string, num)
  remove := make([]string, removeLen)
  partial := make([]string, partialLen)
  var removeCount int = 0
  var partialCount int = 0
  for i:=0; i < num; i++ {
    var tableName string = randomString()
    if (i%2 == 0) {
      partial[partialCount] = tableName
      partialCount += 1
    } else {
      remove[removeCount] = tableName
      removeCount += 1
    }
  }
  return map[string]{"remove": remove, "full": full, "partial": partial}
}


  func TestLoadApplication(t *testing.T) {
    appDB, appErr := NewTestApplicationDB();
    if (appErr != nil) {
      t.Error("Issue with creation of application database", appErr);
    }

    for _, tableName := range tables["full"] {
      table, tableErr := appDB.CreateTable(tableName, cacheExpiration);
      if (tableErr != nil) {
        t.Error("Issue with creation of table", tableErr);
      }
    }

    var removed bool = true;

    for _, tableName := range tables["remove"] {
      removeOutcome, removeErr := appDB.TableRemove(tableName);
      if (removeErr != nil) {
        t.Error("Issue with removing a table", removeErr);
      }
      removedOutcome, removedErr = appDB.GetTable(tableName);
      if (removedOutcome != nil) {
        t.Error("Issue with a table existing after removal");
      }
    }
    _, closeErr := appDB.Close();
    if (closeErr != nil) {
      t.Error("Issue with closing the application", closeErr);
    }

    loadedAppDB, loadErr := LoadProximaDatabase(appName);
    if (loadErr != nil) {
      t.Error("Issue with loading the application", loadErr);
    }

    for _, tableName := range tables["partial"] {
      tableP, tableErrs := appDB.GetTable(tableName)
      if (tableErrs != nil) {
        t.Error("Issue with getting a table", tableErrs);
      }
    }
    appDB.Delete();
  }

func TestApplicationUpdates(t *testing.T) {
  appDB, appErr := NewTestApplicationDB();
  for _, tableName := range tables["full"] {
    table, tableErr := appDB.CreateTable(tableName, cacheExpiration);
    if (tableErr) {
      t.Error("Issue with creating a table", tableErr);
    }
  }
  _, saveErr := appDB.Save();
  if (saveErr) {
    t.Error("Issue with saving the application", saveErr);
  }
  appDB.Close();
  loadedAppDB, loadErr := LoadProximaDatabase(appName);
  if (loadErr) {
    t.Error("Issue with loading the application", loadErr);
  }
  for _, tableName := range tables["full"] {
    table, tableErr := appDB.GetTable(tableName)
    if (tableErr != nil) {
      t.Error("Issue with getting a table", tableErr);
    }
  }
  appDB.Delete();
}
func TestApplicationSync(t *testing.T) {
  appDB, appErr := NewTestApplicationDB();
  for _, tableName := range tables["full"] {
    table, tableErr := appDB.CreateTable(tableName, cacheExpiration);
    if (tableErr) {
      t.Error("Issue with creating a table", tableErr);
    }
  }
  updated, updateErr := appDB.Update();
  if (updateErr != nil) {
    t.Error("Issue with updating the application", updateErr);
  }
  appDB.Close();
  appDB, _ := LoadProximaDatabase(appDBName);
  syncErr := appDB.Sync();
  if (syncErr != nil) {
    t.Error("Issue with syncing the table", syncErr);
  }
  appDB.Delete();
}

func TestTableOperations(*testing.T) {
  appDB, appErr := NewTestApplicationDB();
  for _, tableName := range tables["full"] {
    table, tableErr := appDB.CreateTable(tableName, cacheExpiration);
    if (tableErr) {
      t.Error("Issue with creating a table", tableErr);
    }
  }
  var args map[string]interface{} = map[string]interface{}{"prove":false};
  var sizeValues int = 50;
  var numEntries int = 150;

  for  _, tableName := range tables["full"] {
    var entries map[string]string = GenerateKeyValuePairs(numEntries, 32, sizeValues)
    table, _ := appDB.GetTable(tableName) //
    for key, value :=  range entries {
      _, putErr := table.Put(key, value, args);
      if (putErr) {
        t.Error("Issue with putting value into table", putErr);
      }
      _, getErr := table.Get(key, value, args);
      if (getErr) {
        t.Error("Issue with getting value from table", getErr);
      }
    }
    if (table["isIdle"]) {
      t.Error("Issue with putting value into table");
    }
  }
  updated, updateErr := appDB.Update();
  if (updateErr != nil) {
    t.Error("Issue with updating application", updateErr);
  }
  appDB.Delete();
}
