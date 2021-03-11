package client

import (
  context "context"
  "testing"
  grpc "google.golang.org/grpc"
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




func Setup() (ProximaServiceClient) {
  conn, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
  if err != nil {
    panic(err)
  }
  return NewProximaServiceClient(conn)
}

func TestTable(t *testing.T) {
  name := "NewTable"
  proximaClient := Setup()
  openResponse, openErr := proximaClient.Open(context.Background(), &OpenRequest{Name: name})
  if openErr != nil || !openResponse.GetConfirmation() {
    fmt.Println(openErr)
    closeResponse, closeErr := proximaClient.Close(context.TODO(), &CloseRequest{Name: name})
    if closeErr != nil || !closeResponse.GetConfirmation() {
      t.Error("Cannot close table: ", closeErr)
    }
  }
}

func TestBasicCRUD(t *testing.T) {
  proximaClient := Setup()
  name := "NewTable"
  keySize := 32
  valSize := 300
  keyValues := generateKeyValuePairs(1000, keySize, valSize)
  prove := false
  proximaClient.Open(context.Background(), &OpenRequest{Name: name})
  
  for key, value := range keyValues {
    _, putErr :=  proximaClient.Put(context.TODO(), &PutRequest{Name:  name, Key: []byte(key),
      Value: []byte(value), Prove: prove})
    if putErr != nil {
      t.Error("Issue with putting value into table", putErr)
    }
  }

  for key, _ := range keyValues {
    _, getErr :=  proximaClient.Get(context.TODO(), &GetRequest{Name: name, Key: []byte(key), Prove: prove})
    if getErr != nil {
      t.Error("Issue with putting value into table", getErr)
    }
  }


  for i := 0; i < 10; i++ {
    resp, scanErr :=  proximaClient.Scan(context.TODO(), &ScanRequest{Name: name, First: int32(-1), Last: int32(10+i), Limit: int32(100), Prove: prove})
    //fmt.Println(resp.GetResponses())
    if scanErr != nil {
      fmt.Println(resp.GetResponses())
      t.Error("Issue with scanning table: ", scanErr)
    }
  }

  // for i := 0; i < 10; i++ {
  //   resp, scanErr :=  proximaClient.Scan(context.TODO(), &ScanRequest{Name: name, First: int32(-1), Last: int32(10+i), Limit: int32(100), Prove: prove})
  //   //fmt.Println(resp.GetResponses())
  //   if scanErr != nil {
  //     t.Error("Issue with scanning table: ", scanErr)
  //   }
  // }


}



//
// func TestAdvancedQuery(t *testing.T) {
//   proximaClient := Setup()
//   name := "NewTable"
//   keySize := 32
//   valSize := 300
//   keyValues := generateKeyValuePairs(100, keySize, valSize)
//   prove := false
//   proximaClient.Open(context.Background(), &OpenRequest{Name: name})
//
// }
//
// func TestMultipleTables(t *testing.T) {
//   panic("Not implemented")
// }
//
// func TestSimultaneousOperations(t *testing.T) {
//   panic("Not implemented")
// }
