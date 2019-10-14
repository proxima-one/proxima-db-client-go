package proxima_db_client_go



import (
  "testing"
  _ "math/rand"
  _ "fmt"
)
var db *ProximaDB = NewDatabase(tableName);

func BenchmarkGet(t *testing.B) {
  //tableName := "NewTable"
  numEntries := t.N
  sizeValues := 300
  entries := generateKeyValuePairs(numEntries, 32, sizeValues)

  for key, value := range entries {
    db.Set(tableName, key, value, args)
  }

  for key, _ := range entries {
    db.Get(tableName, key, args)
  }
  //TearDown(tableName, proximaClient)
}

func BenchmarkPut(t *testing.B) {
  numEntries := t.N
  sizeValues := 300
  entries := generateKeyValuePairs(numEntries, 32, sizeValues)

  for key, value := range entries {
    db.Set(tableName, key, value, args)
    // if putErr != nil {
    //   t.Error("Cannot put value: ", putErr)
    // }
  }
  //TearDown(tableName, proximaClient)
}
