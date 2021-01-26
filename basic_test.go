package proxima_db_client_go

import (
  "testing"
  //proxima "github.com/proxima-one/proxima-db-client-go"
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


//docker run proxima-db with ip and port transformation



func BasicDatabaseTest(t *testing.T) {

  db, databaseErr := NewDefaultDatabase(databaseName, databaseID)
  if databaseErr != nil {
    t.Error("Cannot create database: ", databaseErr)
  }

  _, tableErr := db.NewDefaultTable(tableName)
  if tableErr != nil {
    t.Error("Cannot make table: ", tableErr)
  }
  //var pairs map[string]string = GenerateKeyValuePairs(keySize, valueSize, numEntries)
  //var getOperations []interface{} = GenerateDatabaseGetOperations(tableName, pairs)
  //var putOperations []interface{} = GenerateDatabasePutOperations(tableName, pairs)
  //var removeOperations []interface{} = GenerateDatabaseRemoveOperations(tableName, pairs)



  //t.Run("Put", TestPut(t, table, testEntries))
  //t.Run("Get", TestGet(t, table, testEntries))
  //t.Run("Remove", TestRemove(t, table, testEntries))
  // db.Delete()
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
//
func TestGet(t *testing.T) {
  db, databaseErr := NewDefaultDatabase(databaseName, databaseID)
  if databaseErr != nil {
    t.Error("Cannot create database: ", databaseErr)
  }

  table, tableErr := db.NewDefaultTable(tableName)
  if tableErr != nil {
    t.Error("Cannot make table: ", tableErr)
  }

  var entries map[string]string = GenerateKeyValuePairs(keySize, valueSize, numEntries)
  //var putOperations []interface{} = GenerateDatabasePutOperations(tableName, pairs)
//var args map
  for key, _ := range entries {
    _, getErr := table.Get(key, false)
    if getErr != nil {
      t.Error("Cannot get value: ", getErr)
    }
  }
}
// //
// func TestRemove(t *testing.T) {
//   db, databaseErr := NewDefaultDatabase(databaseName, databaseID)
//   if databaseErr != nil {
//     t.Error("Cannot create database: ", databaseErr)
//   }
//
//   table, tableErr := db.NewDefaultTable(tableName)
//   if tableErr != nil {
//     t.Error("Cannot make table: ", tableErr)
//   }
//
//   var entries map[string]string = GenerateKeyValuePairs(keySize, valueSize, numEntries)
//   for key, _ := range entries {
//     _, removeErr := table.Remove(key, false)
//     if removeErr != nil {
//       t.Error("Cannot remove value from key: ", removeErr)
//     }
//   }
// }
//
func TestPut(t *testing.T) {
    db, databaseErr := NewDefaultDatabase(databaseName, databaseID)
    if databaseErr != nil {
      t.Error("Cannot create database: ", databaseErr)
    }

    table, tableErr := db.NewDefaultTable(tableName)
    if tableErr != nil {
      t.Error("Cannot make table: ", tableErr)
    }

    var entries map[string]string = GenerateKeyValuePairs(keySize, valueSize, numEntries)
    //var putOperations []interface{} = GenerateDatabasePutOperations(tableName, pairs)
  //var args map
  for key, value := range entries {
    _, putErr := table.Put(key, value, false, args)
    if putErr != nil {
      t.Error("Cannot put key and value: ", putErr)
    }
  }

  //   for key, _ := range entries {
  //     _, getErr := table.Get(key, args)
  //     if getErr != nil {
  //       t.Error("Cannot get value: ", getErr)
  //     }
  //   }
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
