package proxima_db_client_go

import (
  "testing"
  "math/rand"
  //"fmt"
)

func randomString(size int) (string) {
  bytes := make([]byte, size)
  rand.Read(bytes)
  return string(bytes)
}


func TearDown(name string, db *ProximaDB) {
  db.Close(name)
  db.TableRemove(name)
}

func NewDatabase(name string) (*ProximaDB) {
  ip := "0.0.0.0"
  port := "50051"
  proximaClient := NewProximaDB(ip, port)
  proximaClient.Open(name)
  return proximaClient
}

func Setup(name string, numEntries int, sizeValues int, prove bool) (*ProximaDB, map[string]string, map[string]interface{}) {
  proximaClient := NewDatabase(name)
  entries := generateKeyValuePairs(numEntries, 32, sizeValues)
  args := make(map[string]interface{})
  args["prove"] = prove
  return proximaClient, entries, args
}

// func generateEntries(num int, sizeValues int, prove bool) (map[string]interface{}){
//   mapping := make(map[string]string);
//   for i := 0; i < num; i++ {
//     key := randomString(32)
//     value := randomString(valSize)
//     mapping[key] = value
//   }
//   return mapping
// }

func generateKeyValuePairs(num int, keySize int, valSize int) (map[string]string){
  mapping := make(map[string]string)
  for i := 0; i < num; i++ {
    key := randomString(keySize)
    value := randomString(valSize)
    mapping[key] = value
  }
  return mapping
}

var tableName string = "NewTable";
var name string = "NewTable";
var proximaClient *ProximaDB = NewDatabase(name);
var args map[string]interface{} = map[string]interface{}{"prove":false};
var keyValues map[string]string = generateKeyValuePairs(100000, 32, 300)



func TestPut(t *testing.T) {
  for key, value := range keyValues {
    _, putErr := proximaClient.Set(name, key, value, args)
    if putErr != nil {
      t.Error("Cannot put value: ", putErr)
    }
    //fmt.Println(string(result.GetProof().GetRoot()))
  }
}


func TestGet(t *testing.T) {
  // name := "NewTable"
  // ip := "0.0.0.0"
  // port := "50051"
  // proximaClient := NewProximaDB(ip, port)
  // proximaClient.Open(name)
  // keyValues := generateKeyValuePairs(1, 32, 300)
  // args := make(map[string]interface{})
  // args["prove"] = true

  for key, _ := range keyValues {
    _, getErr := proximaClient.Get(name, key, args)
    if getErr != nil {
      t.Error("Cannot get value: ", getErr)
    }
    //fmt.Println(string(result.GetProof().GetRoot()))
  }

  //proximaClient.Close(name)
}

// func TestArbitraryKeyLength(t *testing.T) {
//   name := "NewTable"
//   ip := "0.0.0.0"
//   port := "50051"
//   proximaClient := NewProximaDB(ip, port)
//   proximaClient.Open(name)
//   keyValues := generateKeyValuePairs(1, 20, 300)
//   args := make(map[string]interface{})
//   //args["prove"] = false
//
//   for key, value := range keyValues {
//     _, putErr := proximaClient.Set(name, key, value, args)
//     if putErr != nil {
//       t.Error("Cannot put value: ", putErr)
//     }
//   }
//
//   for key, _ := range keyValues {
//     _, getErr := proximaClient.Get(name, key, args)
//     if getErr != nil {
//       t.Error("Cannot get value: ", getErr)
//     }
//     //fmt.Println(string(result))
//   }
//   proximaClient.Close(name)
// }

// func TestQuery(t *testing.T) {
//   name := "NewTable"
//   ip := "0.0.0.0"
//   port := "50051"
//   proximaClient := NewProximaDB(ip, port)
//   proximaClient.Open(name)
//   keyValues := generateKeyValuePairs(100, 32, 300)
//   args := make(map[string]interface{})
//   args["prove"] = false
//
//   for key, value := range keyValues {
//     _, putErr := proximaClient.Set(name, key, value, args)
//     if putErr != nil {
//       t.Error("Cannot put value: ", putErr)
//     }
//   }
//
//   for key, _ := range keyValues {
//     _, getErr := proximaClient.Get(name, key, args)
//     if getErr != nil {
//       t.Error("Cannot get value: ", getErr)
//     }
//   }
//   proximaClient.Close(name)
// }
