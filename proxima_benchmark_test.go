package proxima_db_client_go



import (
  "testing"
  _ "math/rand"
  _ "fmt"
)
var db *ProximaDB = NewDatabase(tableName);

// func BenchmarkGet(t *testing.B) {
//   //tableName := "NewTable"
//   numEntries := t.N
//   sizeValues := 300
//   entries := generateKeyValuePairs(numEntries, 32, sizeValues)
//
//   for key, value := range entries {
//     db.Set(tableName, key, value, args)
//   }
//   t.ResetTimer()
//   for key, _ := range entries {
//     _, err := db.Get(tableName, key, args)
//     if err!= nil {
//         t.Error("Cannot get value: ", err)
//     }
//   }
// }
//
// func BenchmarkPut(t *testing.B) {
//   numEntries := t.N
//   sizeValues := 300
//   entries := generateKeyValuePairs(numEntries, 32, sizeValues)
//   t.ResetTimer()
//   for key, value := range entries {
//     _, err := db.Set(tableName, key, value, args)
//     if err != nil {
//        t.Error("Cannot put value: ", err)
//     }
//   }
// }


func BenchmarkBatch(t *testing.B) {
  sizeValues := 300
  numEntries := t.N
  batchSize := 1000
  entries := generateKeyValuePairs(numEntries, 32, sizeValues)
  batchEntries := generateEntries(tableName, entries)
  batches := makeBatches(batchSize, numEntries)

  t.ResetTimer()
  for batch := range batches {
    if batchEntries[:batch] != nil {
      _, err := db.Batch(batchEntries[:batch], args)
      numEntries -= batch
      if err != nil {
        t.Error("Cannot batch values: ", err)
      }
      if batch < numEntries {
        batchEntries = batchEntries[batch:]
      }
    }
  }

}
