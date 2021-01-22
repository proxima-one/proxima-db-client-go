package proxima_db_client_go

import (
  "fmt"
  //"encoding/json"
  json "github.com/json-iterator/go"
  )

func padOrTrimBytes(bb []byte, size int) ([]byte) {
    l := len(bb)
    if l == size {
        return bb
    }
    if l > size {
        return bb[l-size:]
    }
    tmp := make([]byte, size)
    copy(tmp[size-l:], bb)
    return tmp
}

func ProcessKey(key interface{}) ([]byte) {
  byteKey := []byte(fmt.Sprintf("%v", key.(interface{})))
  return padOrTrimBytes(byteKey, 32)
}


func ProcessValue(value interface{}) ([]byte) {
  byteValue, _ := json.Marshal(value)
  return []byte(byteValue)
}

func (db *ProximaDB) Open(tableList []string) (bool, error) {
    for _, tableName := range tableList {
      _, err := db.Open(tableName)
      if err != nil {
        return false, err
      }
    }
    return true, nil
  }

func (db *ProximaDB) Close(tableList []string) (bool, error) {
    for _, tableName := range tableList {
      _, err := db.Close(tableName)
      if err != nil {
        return false, err
      }
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

  func NewDatabaseClient(ip, port string, ) (*proxima_client.ProximaServiceClient, error) {
    address := dbIP + ":" + dbPort
    maxMsgSize := 1024*1024*1024
    conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithDefaultCallOptions(
        grpc.MaxCallRecvMsgSize(maxMsgSize),
        grpc.MaxCallSendMsgSize(maxMsgSize)))
    if err!= nil {
      return nil, err
    } else {
      return proxima_client.NewProximaServiceClient(conn), nil
    }
  }
