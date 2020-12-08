
package proxima_db_client_go

import (
  "context"
  proxima_client "github.com/proxima-one/proxima-db-client-go/client"
  grpc "google.golang.org/grpc"
  //"fmt"
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


func (db *ProximaDatabase) Query(table string, data string, args map[string]interface{}) ([]*ProximaDBResult, error) {
  prove := (args["prove"] != nil) && args["prove"].(bool)
  responses , err := db.client.Query(context.TODO(), &proxima_client.QueryRequest{Name: table, Query: data, Prove: prove})
  if err != nil {
    return nil, err
  }
  proximaResults := make([]*ProximaDBResult, 0)
  for _, response :=  range responses.GetResponses() {
    proximaResults = append(proximaResults, NewProximaDBResult(response.GetValue(), response.GetProof(), response.GetRoot()))
  }
  return proximaResults, nil
}

func (db *ProximaDatabase) Get(table string, k interface{}, args map[string]interface{}) (*ProximaDBResult, error){
  prove := (args["prove"] != nil) && args["prove"].(bool)
  key := ProcessKey(k) //check cache first
  resp, err := db.client.Get(context.TODO(), &proxima_client.GetRequest{Name: table, Key: key, Prove: prove})
  if err != nil {
    return nil, err
  }
  if resp == nil || resp.GetValue() == nil || len(resp.GetValue()) <= 0 {
    return nil, nil
  }
  return NewProximaDBResult(resp.GetValue(), resp.GetProof(), resp.GetRoot()), nil
}

 func (db *ProximaDatabase) Batch(entries []interface{}, args map[string]interface{}) ([]*ProximaDBResult, error) {
   prove := (args["prove"] != nil) && args["prove"].(bool)
   requests := make([]*proxima_client.PutRequest, 0)
   for _, e:= range entries {
    entry:= map[string]interface{}(e.(map[string]interface{}))
    if entry["key"] != nil || entry["value"] != nil {
      key:= ProcessKey(entry["key"])
      value:= ProcessValue(entry["value"])
      requests = append(requests, &proxima_client.PutRequest{Name: string(entry["table"].(string)), Key: key, Value: value, Prove: false})
    }
   }
   responses , err := db.client.Batch(context.TODO(), &proxima_client.BatchRequest{Requests: requests, Prove: prove})
   if err != nil {
     return nil, err
   }
   proximaResults := make([]*ProximaDBResult, 0)
   for _, response :=  range responses.GetResponses() {
     proximaResults = append(proximaResults, NewProximaDBResult([]byte{}, response.GetProof(), response.GetRoot()))
   }
   return proximaResults, nil
 }

func (db *ProximaDatabase) Set(table string, k interface{}, v interface{}, args map[string]interface{}) (*ProximaDBResult, error) {
  prove := (args["prove"] != nil) && args["prove"].(bool)
  if k == nil || v == nil {
    return nil, errors.New("Error with key and value insertion")
  }
  key := ProcessKey(k)
  value := ProcessValue(v) //check cache first
  resp, err := db.client.Put(context.TODO(), &proxima_client.PutRequest{Name: table, Key: key, Value: value, Prove: prove})
  if err != nil {
    return nil, err
  }
  return NewProximaDBResult([]byte{}, resp.GetProof(), resp.GetRoot()), nil
}

func (db *ProximaDatabase) Remove(table string, k interface{}, args map[string]interface{}) (*ProximaDBResult, error) {
  prove := (args["prove"] != nil) && args["prove"].(bool)
  key := ProcessKey(k)
  resp, err := db.client.Remove(context.TODO(), &proxima_client.RemoveRequest{Name: table, Key: key, Prove: prove})
  if err != nil {
    return nil, err
  }
  return NewProximaDBResult([]byte{}, resp.GetProof(), resp.GetRoot()), nil
}
