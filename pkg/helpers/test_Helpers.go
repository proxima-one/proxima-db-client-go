package proxima_db_client_go

import (
  "testing"
  "math/rand"
  "fmt"
  "time"
)

//generate Entries

//generateEntry

//generateTables

//generateTable

//generateTableCache

//generateDBApplication





func RandomString(size int) (string) {
  bytes := make([]byte, size)
  rand.Read(bytes)
  return string(bytes)
}

func NewDatabase(name string) (*ProximaDB) {
  ip := "0.0.0.0"
  port := "50051"
  proximaClient := NewProximaDB(ip, port)
  proximaClient.Open(name)
  return proximaClient
}

func MainSetup(name string, numEntries int, sizeValues int, prove bool) (*ProximaDB, map[string]string, map[string]interface{}) {
  proximaClient := NewDatabase(name)
  entries := generateKeyValuePairs(numEntries, 32, sizeValues)
  args := make(map[string]interface{})
  args["prove"] = prove
  return proximaClient, entries, args
}

func GenerateKeyValuePairings(num int, keySize int, valSize int) (map[string]string){
  mapping := make(map[string]string)
  for i := 0; i < num; i++ {
    key := randomString(keySize)
    value := randomString(valSize)
    mapping[key] = value
  }
  return mapping
}

func GenerateDatabaseGetOperations(tableName string, pairs map[string]string) ([]interface{}){
  entries := make([]interface{}, 0)
  for key, value := range pairs {
    entry:= map[string]interface{}{"key": key, "value": value, "table": tableName, "prove": false}
    entries = append(entries, entry)
  }
  return entries
}

func GenerateDatabasePutOperations(tableName string, pairs map[string]string) ([]interface{}){
  entries := make([]interface{}, 0)
  for key, value := range pairs {
    entry:= map[string]interface{}{"key": key, "value": value, "table": tableName, "prove": false}
    entries = append(entries, entry)
  }
  return entries
}

func GenerateDatabaseRemoveOperations(tableName string, pairs map[string]string) ([]interface{}){
  entries := make([]interface{}, 0)
  for key, value := range pairs {
    entry:= map[string]interface{}{"key": key, "value": value, "table": tableName, "prove": false}
    entries = append(entries, entry)
  }
  return entries
}


// func NewDatabase(name string) (*ProximaDB) {
//   ip := "0.0.0.0"
//   port := "50051"
//   proximaClient := (ip, port)
//   proximaClient.Open(name)
//   return proximaClient
// }


// func generateEntries(num int, sizeValues int, prove bool) (map[string]interface{}){
//   mapping := make(map[string]string);
//   for i := 0; i < num; i++ {
//     key := randomString(32)
//     value := randomString(valSize)
//     mapping[key] = value
//   }
//   return mapping
// }
// 
// func generateKeyValuePairs(num int, keySize int, valSize int) (map[string]string){
//   mapping := make(map[string]string)
//   for i := 0; i < num; i++ {
//     key := randomString(keySize)
//     value := randomString(valSize)
//     mapping[key] = value
//   }
//   return mapping
// }
//
// func generateEntries(name string, pairs map[string]string) ([]interface{}){
//   entries := make([]interface{}, 0)
//   for key, value := range pairs {
//     entry:= map[string]interface{}{"key": key, "value": value, "table": name, "prove": false}
//     entries = append(entries, entry)
//   }
//   return entries
// }
//
// func makeBatches(batchSize, total int) ([]int) {
//   batches := make([]int, 0)
//   num := 0
//   for total > (num + 1) {
//     num += batchSize
//     batches = append(batches, num)
//   }
//   return append(batches, total)
// }

// func GenerateDatabaseQueryOperations(tableName string, queries []string) ([]interface{}){
//   entries := make([]interface{}, 0)
//   for key, value := range queries {
//     entry:= map[string]interface{}{"queryText": queryText, "table": tableName, "prove": false}
//     entries = append(entries, entry)
//   }
//   return entries
// }
