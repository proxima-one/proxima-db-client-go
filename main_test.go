package proxima_db_client_go

// import (
//   "testing"
//   "math/rand"
//   "fmt"
//   "time"
// )
//
// func randomString(size int) (string) {
//   bytes := make([]byte, size)
//   rand.Read(bytes)
//   return string(bytes)
// }
//
//
// func MainTestTearDown(name string, db *ProximaDB) {
//   db.Close(name)
//   db.TableRemove(name)
// }
//
// // func NewDatabase(name string) (*ProximaDB) {
// //   ip := "0.0.0.0"
// //   port := "50051"
// //   proximaClient := (ip, port)
// //   proximaClient.Open(name)
// //   return proximaClient
// // }
//
// func MainTestSetup(name string, numEntries int, sizeValues int, prove bool) (*ProximaDB, map[string]string, map[string]interface{}) {
//   proximaClient := NewDatabase(name)
//   entries := generateKeyValuePairs(numEntries, 32, sizeValues)
//   args := make(map[string]interface{})
//   args["prove"] = prove
//   return proximaClient, entries, args
// }
//
// // func generateEntries(num int, sizeValues int, prove bool) (map[string]interface{}){
// //   mapping := make(map[string]string);
// //   for i := 0; i < num; i++ {
// //     key := randomString(32)
// //     value := randomString(valSize)
// //     mapping[key] = value
// //   }
// //   return mapping
// // }
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
//
// //update test
//
// var tableName string = "NewTable";
// var name string = "NewTable";
// var proximaClient *ProximaDB = NewDatabase(name);
// var args map[string]interface{} = map[string]interface{}{"prove":false};
// var keyValues map[string]string = generateKeyValuePairs(100, 32, 200)
// var sizeValues int = 50;
// var numEntries int = 15000;
// var entries map[string]string = generateKeyValuePairs(numEntries, 32, sizeValues)
// var batchEntries []interface{} = generateEntries(tableName, entries)
//
// //
//
// func TestCreateDatabase(t *testing.T) {
//   var proximaClient *ProximaDB = NewDatabase(name);
//     _, putErr := proximaClient.Set(name, key, value, args)
//     if putErr != nil {
//       t.Error("Cannot put value: ", putErr)
//     }
//
//     if tableErr != nil {
//       t.Error("Cannot make table: ", tableErr)
//     }
//     //fmt.Println(string(result.GetProof().GetRoot()))
// }
//
// func TestLoadDatabase(t *testing.T) {
//   if tableErr != nil {
//     t.Error("Cannot make table: ", tableErr)
//   }
// }
//
//
// func TestPut(t *testing.T) {
//   for key, value := range keyValues {
//     _, putErr := proximaClient.Set(name, key, value, args)
//     if putErr != nil {
//       t.Error("Cannot put value: ", putErr)
//     }
//     //fmt.Println(string(result.GetProof().GetRoot()))
//   }
// }
// //
// //
// func TestGet(t *testing.T) {
//   for key, _ := range keyValues {
//     _, getErr := proximaClient.Get(name, key, args)
//     if getErr != nil {
//       t.Error("Cannot get value: ", getErr)
//     }
//   }
// }
//
// func TestBatch(t *testing.T) {
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
