package proxima_client

import (
  context "context"
  "testing"
  grpc "google.golang.org/grpc"
  "math/rand"
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
    //t.Error("Cannot open table: ", openErr)
    closeResponse, closeErr := proximaClient.Close(context.Background(), &CloseRequest{Name: name})
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
  //_, openErr := proximaClient.Open(context.Background(), &OpenRequest{Name: name})

  for key, value := range keyValues {
    _, putErr :=  proximaClient.Put(context.Background(), &PutRequest{Name:  name, Key: []byte(key),
      Value: []byte(value), Prove: prove})
    if putErr != nil {
      t.Error("Issue with putting value into table", putErr)
    }
  }

  for key, _ := range keyValues {
    _, getErr :=  proximaClient.Get(context.Background(), &GetRequest{Name: name, Key: []byte(key), Prove: prove})
    if getErr != nil {
      t.Error("Issue with putting value into table", getErr)
    }
  }
}
//
// func TestAdvancedQuery(t *testing.T) {
//   panic("Not implemented")
// }
//
// func TestMultipleTables(t *testing.T) {
//   panic("Not implemented")
// }
//
// func TestSimultaneousOperations(t *testing.T) {
//   panic("Not implemented")
// }
