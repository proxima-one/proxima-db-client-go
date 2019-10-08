package proxima_db_client_go


import (
  "context"
  "fmt"
  proxima_client "github.com/proxima-one/proxima-db-client-go/client"
  grpc "google.golang.org/grpc"
)


type ProximaDBResult struct {
  value []byte
  proof *ProximaDBProof
}

func (db *ProximaDBResult) GetValue() ([]byte) {
  return db.value
}

func (db *ProximaDBResult) GetProof() (*ProximaDBProof) {
  return db.proof
}

type ProximaDBProof struct {
  root []byte
  proof []byte
}

func (pf *ProximaDBProof) GetRoot() ([]byte) {
  return pf.root
}

func (pf *ProximaDBProof) GetProof() ([]byte) {
  return pf.proof
}

func NewProximaDBResult(value, proof, root []byte) (*ProximaDBResult) {
  return &ProximaDBResult{value: value, proof: &ProximaDBProof{root: root, proof: proof}}
}

type ProximaDB struct {
  client proxima_client.ProximaServiceClient
  tables []string
}

func NewProximaDB(dbIP, dbPort string) (*ProximaDB) {
  address := dbIP + ":" + dbPort
  conn, _ := grpc.Dial(address, grpc.WithInsecure())
  return &ProximaDB{client: proxima_client.NewProximaServiceClient(conn), tables: []string{}}
}

func (db *ProximaDB) Create(tableName string) (bool, error) {
  resp, err := db.client.Open(context.TODO(), &proxima_client.OpenRequest{Name: tableName})
  if err != nil {
    return false, err
  }
  db.tables = append(db.tables, tableName)
  return resp.GetConfirmation(), nil
}

func (db *ProximaDB) Open(tableName string) (bool, error) {
  resp, err:= db.client.Open(context.TODO(), &proxima_client.OpenRequest{Name: tableName})
  if err != nil {
    return false, nil
  }
  return resp.GetConfirmation(), nil
}

func (db *ProximaDB) Close(tableName string) (bool, error) {
  _, err:= db.client.Close(context.TODO(), &proxima_client.CloseRequest{Name: tableName})
  if err != nil {
    return false, err
  }
  return true, nil
}

func (db *ProximaDB) CreateAll(tableList []string) (bool, error) {
  for _, tableName := range tableList {
    _, err := db.Create(tableName)
    if err != nil {
      return false, err
    }
  }
  return true, nil
}

func (db *ProximaDB) OpenAll(tableList []string) (bool, error) {
    for _, tableName := range tableList {
      _, err := db.Open(tableName)
      if err != nil {
        return false, err
      }
    }
    return true, nil
  }

func (db *ProximaDB) CloseAll(tableList []string) (bool, error) {
    for _, tableName := range tableList {
      _, err := db.Close(tableName)
      if err != nil {
        return false, err
      }
    }
    return true, nil
  }

func (db *ProximaDB) Query(table string, data string, args map[string]interface{}) (*[]ProximaDBResult, error) {
  prove := args["prove"].(bool)
  responses , err := db.client.Query(context.TODO(), &proxima_client.QueryRequest{Name: table, Query: data, Prove: prove})
  if err != nil {
    return nil, err
  }
  proximaResults := make([]ProximaDBResult, 0)
  for _, response :=  range responses.GetResponses() {
    proximaResults = append(proximaResults, *NewProximaDBResult(response.GetValue(), response.GetProof(), response.GetRoot()))
  }
  return &proximaResults, nil
}

func (db *ProximaDB) Get(table string, k interface{}, args map[string]interface{}) (*ProximaDBResult, error){
  prove := args["prove"].(bool)
  key := []byte(fmt.Sprintf("%v", k.(interface{})))
  resp, err := db.client.Get(context.TODO(), &proxima_client.GetRequest{Name: table, Key: key, Prove: prove})
  if err != nil {
    return nil, err
  }
  return NewProximaDBResult(resp.GetValue(), resp.GetProof(), resp.GetRoot()), nil
}

func (db *ProximaDB) Set(table string, k interface{}, value interface{}, args map[string]interface{}) (*ProximaDBResult, error) {
  prove := args["prove"].(bool)
  key := []byte(fmt.Sprintf("%v", k.(interface{})))
  val := []byte(fmt.Sprintf("%v", value.(interface{}))) //have to convert the value to bytes
  resp, err := db.client.Put(context.TODO(), &proxima_client.PutRequest{Name: table, Key: key, Value: val, Prove: prove})
  if err != nil {
    return nil, err
  }
  return NewProximaDBResult([]byte{}, resp.GetProof(), resp.GetRoot()), nil
}

func (db *ProximaDB) Remove(table string, k interface{}, args map[string]interface{}) (*ProximaDBResult, error) {
  prove := args["prove"].(bool)
  key := []byte(fmt.Sprintf("%v", k.(interface{})))
  resp, err := db.client.Remove(context.TODO(), &proxima_client.RemoveRequest{Name: table, Key: key, Prove: prove})
  if err != nil {
    return nil, err
  }
  return NewProximaDBResult([]byte{}, resp.GetProof(), resp.GetRoot()), nil
}
