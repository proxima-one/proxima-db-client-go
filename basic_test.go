package proxima_db_client_go

import (
  "testing"
  database "github.com/proxima-one/proxima-db-client-go/database"
  _ "math/rand"
  _ "fmt"
  _ "time"
)


var databaseName string = "DefaultDatabaseName";
var databaseID string = "DefaultDatabaseID";
var tableName string = "DefaultTableName";

var valueSize int = 50;
var numEntries int = 1500;
var numBatches int = 5;
var keySize int = 32;
var args map[string]interface{} = map[string]interface{}{"prove":false};



func BasicDatabaseTest(t *testing.T) {

  db, databaseErr := database.NewDefaultDatabase(databaseName, databaseID)
  if databaseErr != nil {
    t.Error("Cannot create database: ", databaseErr)
  }

  table, tableErr := db.NewDefaultTable(tableName)
  if tableErr != nil {
    t.Error("Cannot make table: ", tableErr)
  }

  var testEntries map[string]string = GenerateKeyValuePairs(keySize, valueSize, numEntries, args)
  t.Run("Put", t.TestPut(table, testEntries))
  t.Run("Get", t.TestGet(table, testEntries)
  t.Run("Remove", t.TestRemove(table, testEntries))
  db.Delete()
}

// func AdvancedDatabaseTest1(t *testing.T) {
//   //db operations
//   db, databaseErr := NewDefaultDatabase(databaseName, databaseID)
//   if databaseErr != nil {
//     t.Error("Cannot create database: ", databaseErr)
//   }
//
//   table, tableErr := db.NewDefaultTable(tableName)
//   if tableErr != nil {
//     t.Error("Cannot make table: ", tableErr)
//   }
//   db.Delete()
// }
//
// func AdvancedTableTest2(t *testing.T) {
//   var testEntries map[string]string = GenerateKeyValuePairs(keySize, valueSize, numEntries, args)
//   t.Run("Put", t.TestPut(table, testEntries))
//   t.Run("Get", t.TestGet(table, testEntries)
//   t.Run("Remove", t.TestRemove(table, testEntries))
// }

func (t *testing.T) TestGet(table *ProximaTable, entries map[string]interface{}) {
  for key, _ := range entries {
    _, getErr := table.Get(key, args)
    if getErr != nil {
      t.Error("Cannot get value: ", getErr)
    }
  }
}

func (t *testing.T) TestRemove(table *ProximaTable, entries map[string]interface{}) {
  for key, _ := range entries {
    _, removeErr := table.Remove(key, args)
    if removeErr != nil {
      t.Error("Cannot remove value from key: ", removeErr)
    }
  }
}

func (t *testing.T) TestPut(table *ProximaTable, entries map[string]interface{}) {
  for key, value := range entries {
    _, putErr := table.Put(key, value, args)
    if putErr != nil {
      t.Error("Cannot put key and value: ", getErr)
    }
  }
}

// func (t *testing.T) BasicTableTest() {
//   //create, get, Put, remove ...
//   t.Run("Put", t.TestPut(db, tableName, testEntries))
//   t.Run("Get", t.TestGet(db, tableName, testEntries)
//   t.Run("Remove", t.TestRemove(db, tableName, testEntries))
//   //TestUpdate()
//   //TestSync()
//   //TestTableConfig()
//   //TestLoad
//   //TestCheckMax
//   //TestGetClients
// }
//
//
// func (t *testing.T) TestCreation() {
//     //client, clientErr :=
//     if clientErr != nil {
//       t.Error("Error with creating the client: ", clientErr)
//     }
//     //db, databaseErr := NewDefaultDatabse()
//     if databaseErr != nil {
//       t.Error("Cannot create database: ", databaseErr)
//     }
//     //table, tableErr := db.NewDefaultTable()
//     if tableErr != nil {
//       t.Error("Cannot make table: ", tableErr)
//     }
//     //fmt.Println(string(result.GetProof().GetRoot()))
// }







// func (t *testing.T) TestBatch(db *ProximaDatabase, tableName string, entries []interface{}) {
//   // sizeValues := 300
//   // numEntries := 3000
//   // entries := generateKeyValuePairs(numEntries, 32, sizeValues)
//   // batchEntries := generateEntries(tableName, entries)
//   start := time.Now()
//   _, err := proximaClient.Batch(batchEntries, args)
//   if err != nil {
//     t.Error("Cannot batch values: ", err)
//   }
//   end := time.Now()
//   elapsed := end.Sub(start)
//   fmt.Println(elapsed)
// }
