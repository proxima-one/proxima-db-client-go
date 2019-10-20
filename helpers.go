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
