package proxima_db_client_go

import (
  "testing"
  "math/rand"
  "fmt"
)

func randomString(size int) (string) {
  bytes := make([]byte, size)
  rand.Read(bytes)
  return string(bytes)
}

func generateKeyValuePairs(num int, keySize int, valSize int) (map[string]string){
  mapping := make(map[string]string)
  for i := 0; i < num; i++ {
    key := randomString(keySize)
    value := randomString(valSize)
    mapping[key] = value
  }
  return mapping
}

func TestGetPut(t *testing.T) {
  name := "NewTable"
  ip := "0.0.0.0"
  port := "50051"
  proximaClient := NewProximaDB(ip, port)
  proximaClient.Open(name)
  keyValues := generateKeyValuePairs(100, 32, 300)
  args := make(map[string]interface{})
  args["prove"] = false

  for key, value := range keyValues {
    fmt.Println(key)
    _, putErr := proximaClient.Set(name, key, value, args)
    if putErr != nil {
      t.Error("Cannot put value: ", putErr)
    }
  }

  for key, _ := range keyValues {
    result, getErr := proximaClient.Get(name, key, args)
    if getErr != nil {
      t.Error("Cannot get value: ", getErr)
    }
    fmt.Println(result.GetValue())
  }
  proximaClient.Close(name)
}

func TestArbitraryKeyLength(t *testing.T) {
  name := "NewTable"
  ip := "0.0.0.0"
  port := "50051"
  proximaClient := NewProximaDB(ip, port)
  proximaClient.Open(name)
  keyValues := generateKeyValuePairs(100, 20, 300)
  args := make(map[string]interface{})
  //args["prove"] = false

  for key, value := range keyValues {
    _, putErr := proximaClient.Set(name, key, value, args)
    if putErr != nil {
      t.Error("Cannot put value: ", putErr)
    }
  }

  for key, _ := range keyValues {
    result, getErr := proximaClient.Get(name, key, args)
    if getErr != nil {
      t.Error("Cannot get value: ", getErr)
    }
    fmt.Println(result)
  }
  proximaClient.Close(name)
}

func TestQuery(t *testing.T) {
  name := "NewTable"
  ip := "0.0.0.0"
  port := "50051"
  proximaClient := NewProximaDB(ip, port)
  proximaClient.Open(name)
  keyValues := generateKeyValuePairs(100, 32, 300)
  args := make(map[string]interface{})
  args["prove"] = false

  for key, value := range keyValues {
    _, putErr := proximaClient.Set(name, key, value, args)
    if putErr != nil {
      t.Error("Cannot put value: ", putErr)
    }
  }

  for key, _ := range keyValues {
    _, getErr := proximaClient.Get(name, key, args)
    if getErr != nil {
      t.Error("Cannot get value: ", getErr)
    }
  }
  proximaClient.Close(name)
}
