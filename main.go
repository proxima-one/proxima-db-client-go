package proxima_db_client_go

import (
	//client "github.com/proxima-one/proxima-db-client-go/pkg/client"
	"fmt"
	"os"
	"time"

	database "github.com/proxima-one/proxima-db-client-go/pkg/database"
)

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}
	return defaultVal
}

func NewDefaultDatabase(name, id string) (*database.ProximaDatabase, error) {
	ip := getEnv("DB_ADDRESS", "0.0.0.0")
	port := getEnv("DB_PORT", "50051")

	proxima_client, err := database.DefaultProximaServiceClient(ip, port)
	if err != nil {
		fmt.Println("Error creating database", err)
		proxima_client = nil
	}
	clients := make([]interface{}, 0)
	sleepInterval, _ := time.ParseDuration("5m")
	compressionInterval, _ := time.ParseDuration("5m")
	batchingInterval, _ := time.ParseDuration("5m")
	//cacheExpiration, _ := time.ParseDuration("5m")

	return database.NewProximaDatabase(name, id, "0.0.0.0", proxima_client, clients, sleepInterval,
		compressionInterval, batchingInterval)
}
