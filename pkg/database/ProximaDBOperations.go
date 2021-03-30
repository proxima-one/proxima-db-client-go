package database

import (
	"context"
	"encoding/json"

	client "github.com/proxima-one/proxima-db-client-go/pkg/client"
	iterables "github.com/proxima-one/proxima-db-client-go/pkg/iterables"

	//grpc "google.golang.org/grpc"
	//"fmt"

	"errors"
)

type ProximaDBResult struct {
	value []byte
	proof *ProximaDBProof
}

func (db *ProximaDBResult) GetValue() []byte {
	return db.value
}

func (db *ProximaDBResult) GetProof() *ProximaDBProof {
	return db.proof
}

type ProximaDBProof struct {
	root  []byte
	proof []byte
}

func (pf *ProximaDBProof) GetRoot() []byte {
	return pf.root
}

func (pf *ProximaDBProof) GetProof() []byte {
	return pf.proof
}

func NewProximaDBResult(value, proof, root []byte) *ProximaDBResult {
	return &ProximaDBResult{value: value, proof: &ProximaDBProof{root: root, proof: proof}}
}

func (db *ProximaDatabase) Query(table string, data string, args map[string]interface{}) ([]*ProximaDBResult, error) {
	prove := (args["prove"] != nil) && args["prove"].(bool)
	responses, err := db.client.Query(context.TODO(), &client.QueryRequest{Name: table, Query: data, Prove: prove})
	if err != nil {
		return nil, err
	}
	proximaResults := make([]*ProximaDBResult, 0)
	for _, response := range responses.GetResponses() {
		proximaResults = append(proximaResults, NewProximaDBResult(response.GetValue(), response.GetProof(), response.GetRoot()))
	}
	return proximaResults, nil
}

func (db *ProximaDatabase) Get(table string, k interface{}, args map[string]interface{}) (*ProximaDBResult, error) {
	prove := (args["prove"] != nil) && args["prove"].(bool)
	key := ProcessKey(k) //check cache first
	resp, err := db.client.Get(context.TODO(), &client.GetRequest{Name: table, Key: key, Prove: prove})
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
	requests := make([]*client.PutRequest, 0)
	for _, e := range entries {
		entry := map[string]interface{}(e.(map[string]interface{}))
		if entry["key"] != nil || entry["value"] != nil {
			key := ProcessKey(entry["key"])
			value := ProcessValue(entry["value"])
			requests = append(requests, &client.PutRequest{Name: string(entry["table"].(string)), Key: key, Value: value, Prove: false})
		}
	}
	responses, err := db.client.Batch(context.TODO(), &client.BatchRequest{Requests: requests, Prove: prove})
	if err != nil {
		return nil, err
	}
	proximaResults := make([]*ProximaDBResult, 0)
	for _, response := range responses.GetResponses() {
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
	resp, err := db.client.Put(context.TODO(), &client.PutRequest{Name: table, Key: key, Value: value, Prove: prove})
	if err != nil {
		return nil, err
	}
	return NewProximaDBResult([]byte{}, resp.GetProof(), resp.GetRoot()), nil
}

//Search on the db, only do it in tranches

//search filter for proximaDB with list of variables
func Search(objs []*ProximaDBResult, filterSet, orderBy, orderDirection interface{}, first, last int, args interface{}) ([]*ProximaDBResult, error) {
	mapFn := func(e interface{}) interface{} {
		val := e.(*ProximaDBResult).GetValue()
		var entity map[string]interface{}
		json.Unmarshal(val, &entity)
		return entity
	}
	result, err := iterables.Search(objs, filterSet, orderBy, orderDirection, first, last, mapFn)
	return result.([]*ProximaDBResult), err
}

func (db *ProximaDatabase) Scan(table string, first int, last int, limit int, args map[string]interface{}) ([]*ProximaDBResult, error) {
	finish := first
	direction := (first < last)
	if !direction {
		finish = last
	}
	return db.Range(table, 0, finish, direction, args)
}

func (db *ProximaDatabase) Range(table string, start int, finish int, direction bool, args map[string]interface{}) ([]*ProximaDBResult, error) {
	prove := (args["prove"] != nil) && args["prove"].(bool)
	limit := -1
	if direction {
		limit = 500
	}
	responses, err := db.client.Scan(context.TODO(), &client.ScanRequest{Name: table, First: int32(start), Last: int32(finish), Limit: int32(limit), Prove: prove})
	if err != nil {
		return nil, err
	}
	proximaResults := make([]*ProximaDBResult, 0)
	for _, response := range responses.GetResponses() {
		proximaResults = append(proximaResults, NewProximaDBResult(response.GetValue(), response.GetProof(), response.GetRoot()))
	}
	return proximaResults, nil
}

//
func (db *ProximaDatabase) Stat(table string) (*ProximaDBResult, error) {
	resp, err := db.client.Stat(context.TODO(), &client.StatRequest{Name: table})
	if err != nil {
		return nil, err
	}
	return NewProximaDBResult(resp.Stats, resp.Proof, resp.Root), nil
}

//
func (db *ProximaDatabase) Compact(table string) (bool, error) {
	resp, err := db.client.Compact(context.TODO(), &client.CompactRequest{Name: table})
	if err != nil {
		return false, err
	}
	return resp.Confirmation, nil
}

func (db *ProximaDatabase) Checkout(table string) (bool, error) {
	resp, err := db.client.Checkout(context.TODO(), &client.CheckoutRequest{Name: table})
	if err != nil {
		return false, err
	}
	return resp.Confirmation, nil
}

//
func (db *ProximaDatabase) Commit(table string) (bool, error) {
	resp, err := db.client.Commit(context.TODO(), &client.CommitRequest{Name: table})
	if err != nil {
		return false, err
	}
	return resp.Confirmation, nil
}

func (db *ProximaDatabase) Remove(table string, k interface{}, args map[string]interface{}) (*ProximaDBResult, error) {
	prove := (args["prove"] != nil) && args["prove"].(bool)
	key := ProcessKey(k)
	resp, err := db.client.Remove(context.TODO(), &client.RemoveRequest{Name: table, Key: key, Prove: prove})
	if err != nil {
		return nil, err
	}
	return NewProximaDBResult([]byte{}, resp.GetProof(), resp.GetRoot()), nil
}
