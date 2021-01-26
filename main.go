package proxima_db_client_go


import (
  //client "github.com/proxima-one/proxima-db-client-go/pkg/client"
  database "github.com/proxima-one/proxima-db-client-go/pkg/database"
  "os"
  "time"
)


func getEnv(key, defaultVal string) (string) {
  val := os.Getenv(key)
  if val != "" {
    return val
  }
  return defaultVal
}

func NewDefaultDatabase(name, id string) (*database.ProximaDatabase, error) {
  ip := getEnv("DB_ADDRESS" , "0.0.0.0")
  port :=  getEnv("DB_PORT", "50051")

  proxima_client, err := database.DefaultProximaServiceClient(ip, port)
  if err != nil {
    proxima_client = nil
  }
  clients := make([]interface{}, 0)
  sleepInterval, _ := time.ParseDuration("5m")
  compressionInterval, _ := time.ParseDuration("5m")
  batchingInterval, _ := time.ParseDuration("5m")

  return database.NewProximaDatabase(name, id, "0.0.0.0", proxima_client, clients, sleepInterval,
    compressionInterval, batchingInterval)
}
