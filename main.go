package proxima_db_client_go


import (
  client "github.com/proxima-one/proxima-db-client-go/client"
  database "github.com/proxima-one/proxima-db-client-go/database"
)

func DefaultProximaServiceClient(dbIP, dbPort string) (client.ProximaServiceClient, error)   {
  address := dbIP + ":" + dbPort
  maxMsgSize := 1024*1024*1024
  conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithDefaultCallOptions(
      grpc.MaxCallRecvMsgSize(maxMsgSize),
      grpc.MaxCallSendMsgSize(maxMsgSize)))
  if err != nil {
    return nil, err
  }
  return client.NewProximaServiceClient(conn), nil
}

func NewDefaultDatabase(name, id string) (*database.ProximaDatabase, error) {
  ip := getEnv("DB_ADDRESS" , "0.0.0.0")
  port :=  getEnv("DB_PORT", "50051")

  proxima_client, err := DefaultProximaServiceClient(ip, port)
  if err != nil {
    proxima_client = make(*client.ProximaServiceClient)
  }
  clients := make([]interface{})
  sleepInterval := time.ParseDuration()
  compressionInterval := time.ParseDuration()
  batchingInterval := time.ParseDuration()

  return database.NewProximaDatabase(name, id, "0.0.0.0", proxima_client, clients, sleepInterval,
    compressionInterval, batchingInterval), nil
}
