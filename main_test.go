package proxima_db_client_go

import (
	"testing"
	proxima_database "github.com/proxima-one/proxima-db-client-go/pkg/database"
	_ "fmt"
	"math/rand"
	_ "time"
)

var databaseName string = "DefaultDatabaseName"
var databaseID string = "DefaultDatabaseID"
var tableName string = "DefaultTableName"

var valueSize int = 50
var numEntries int = 1500
var numBatches int = 5
var keySize int = 32
var args map[string]interface{} = map[string]interface{}{"prove": false}
var testTableConfig map[string]interface{} = map[string]interface{}{"name": "DPoolLists",
		  "id": "DefaultDB-DPoolLists",
		  "version": "0.0.0",
		  "blockNum": 0,
		  "header": "Root",
		  "compression": "36h",
		  "batching": "500ms",
		  "sleep": "10m",
		  "cacheExpiration": "5m"}

var testDatabaseConfig map[string]interface{} = map[string]interface{}{
  "name": "DefaultDB",
  "id": "DefaultID",
  "owner": "None",
  "version": "0.0.0",
  "config": map[string]interface{}{
    "sleep": "5m",
    "compression": "36h",
    "batching": "500ms",
  },
  "tables": []interface{}{testTableConfig},
}

func CopyMapNewVersion(newVersion string, originalMap map[string]interface{}) (map[string]interface{}) {
  CopiedMap:= make(map[string]interface{})

  /* Copy Content from Map1 to Map2*/
  for index, element  := range originalMap{
     CopiedMap[index] = element
  }
  CopiedMap["version"] = newVersion
  return CopiedMap
}

func TestCheckLatest(t *testing.T) {
  testConfigMap := make(map[string]interface{})

  networkTableConfig := CopyMapNewVersion("0.0.1", testTableConfig)
  nodeTableConfig := CopyMapNewVersion("0.0.5", testTableConfig)
  localTableConfig := CopyMapNewVersion("0.0.2", testTableConfig)
  currentTableConfig := CopyMapNewVersion("0.0.1", testTableConfig)

  tableTestConfigMap := map[string]interface{}{
    "network":networkTableConfig,
    "node":nodeTableConfig,
    "local":localTableConfig,
    "current":currentTableConfig,
  }
  testConfigMap["table"] = tableTestConfigMap

  networkDatabaseConfig := CopyMapNewVersion("0.0.1", testDatabaseConfig)
  nodeDatabaseConfig := CopyMapNewVersion("0.0.5", testDatabaseConfig)
  localDatabaseConfig := CopyMapNewVersion("0.0.2", testDatabaseConfig)
  currentDatabaseConfig := CopyMapNewVersion("0.0.1", testDatabaseConfig)

  databaseTestConfigMap := map[string]interface{}{
    "network":networkDatabaseConfig,
    "node":nodeDatabaseConfig,
    "local":localDatabaseConfig,
    "current":currentDatabaseConfig,
  }
  testConfigMap["database"] = databaseTestConfigMap

  var resp map[string]interface{}
  for configTestName, latestTestConfig := range testConfigMap {
    if configTestName == "table" {
        resp, _ = proxima_database.CheckLatest("version", latestTestConfig.(map[string]interface{}))
    } else {
        resp, _ = proxima_database.CheckLatest("version", latestTestConfig.(map[string]interface{}))
    }

    latestName := resp["type"].(string)
    if latestName == "" || latestName != "node" {
      t.Errorf("Issues with checking the latest. Expected %v, Actual %v", "latest", latestName)
    }
  }
}

func TestDatabaseCreation(t *testing.T) {
    db, databaseErr := NewDefaultDatabase(databaseName, databaseID)
    if databaseErr != nil {
      t.Error("Cannot create database: ", databaseErr)
    }

    dbConfig := testDatabaseConfig
    dbConfig["version"] = "0.0.1"

    db.SetCurrentDatabaseConfig(dbConfig, true)
    actualConfig, _ := db.GetCurrentDatabaseConfig()
    if actualConfig == nil ||  actualConfig["version"].(string) != dbConfig["version"].(string) {
      t.Errorf("Issues with setting current database config expected: %v but got: %v", dbConfig["version"].(string), actualConfig["version"].(string))
    }

    dbConfig["version"] = "0.0.2"

    db, dbErr := proxima_database.LoadProximaDatabase(dbConfig)
    if dbErr != nil {
      t.Error("Issues with loading database from config: ", dbErr)
    }
    actualConfig, _ = db.GetCurrentDatabaseConfig()
    if actualConfig["version"].(string) != dbConfig["version"].(string) {
      t.Errorf("Issues with loading database config. expected: %v but got: %v", dbConfig["version"].(string), actualConfig["version"].(string))
    }
}



func TestTableCreation(t *testing.T) {
    db, databaseErr := NewDefaultDatabase(databaseName, databaseID)
    if databaseErr != nil {
      t.Error("Cannot create database: ", databaseErr)
    }

    table, tableErr := db.NewDefaultTable(tableName)
    if tableErr != nil {
      t.Error("Cannot make table: ", tableErr)
    }
    tableConfig := testTableConfig
    tableConfig["version"] = "0.0.1"

    table.SetCurrentTableConfig(tableConfig)
    actualConfig, _ := table.GetCurrentTableConfig()
    if actualConfig["version"].(string) != tableConfig["version"].(string) {
      t.Errorf("Issues with setting current table config expected: %v but got: %v", tableConfig["version"].(string), actualConfig["version"].(string))
    }

    tableConfig["version"] = "0.0.2"

    table.Load("", tableConfig)
    actualConfig, _ = table.GetCurrentTableConfig()
		//fmt.Println(actualConfig)
    if actualConfig["version"].(string) != tableConfig["version"].(string) {
      t.Errorf("Issues with loading table config. expected: %v but got: %v", tableConfig["version"].(string), actualConfig["version"].(string))
    }
}



func TestBasicDatabase(t *testing.T) {
	db, databaseErr := NewDefaultDatabase(databaseName, databaseID)
	if databaseErr != nil {
		t.Error("Cannot create database: ", databaseErr)
	}
	var numTables int = 10
	var numRemovedTables int = 2

	tableList := GenerateTableList(numTables)

	for _, tableName := range tableList {
		_, tableErr := db.NewDefaultTable(tableName)
		if tableErr != nil {
			t.Error("Cannot make table: ", tableErr)
		}
	}

	for _, tableName := range tableList {
		_, tableErr := db.GetTable(tableName)
		if tableErr != nil {
			t.Error("Cannot make table: ", tableErr)
		}
		//additional test objectives, get, put, remove, etc
      //get
      //put
      //remove
	}

	for i := 0; i < numRemovedTables; i++ {
		i := rand.Intn(numTables - i)
		tableName := tableList[i]
		_, removeErr := db.RemoveTable(tableName)
		if removeErr != nil {
			t.Error("Issue with removing table: ", removeErr)
		}

		table, tableErr := db.GetTable(tableName)
		if tableErr != nil {
			t.Error("Error on removed table: ", tableErr)
		}
		if table != nil {
			t.Error("Issue with removing table, table is not nil")
		}
	}
}
//
//
//
// func TestDatabaseConfig(t *testing.T) {
// 	db, databaseErr := NewDefaultDatabase(databaseName, databaseID)
// 	if databaseErr != nil {
// 		t.Error("Cannot create database: ", databaseErr)
// 	}
//   configUpdates := make(map[string]interface{})
//   expectedVersion := "0.0.0"
// 	version1 := "0.0.1"
// 	version2 := "0.0.2"
//
//   config, _ := db.GetCurrentDatabaseConfig()
//   actualVersion := config["version"].(string)
//   if actualVersion != expectedVersion {
//     t.Errorf("Did not initialize database correctly, got version: %v, expectedVersion: %v", actualVersion, expectedVersion)
//   }
//   configUpdates["version"] = version2
//   db.SetCurrentDatabaseConfig(configUpdates, false)
//   expectedVersion = version2
//
//   config, _ = db.GetCurrentDatabaseConfig()
//   actualVersion = config["version"].(string)
//   if actualVersion != "" || actualVersion != expectedVersion {
//     t.Errorf("Did not update version correctly, got version: %v, expectedVersion: %v", actualVersion, expectedVersion)
//   }
// 	db.Update()
//
//   config, _ = db.GetCurrentDatabaseConfig()
//   actualVersion = config["version"].(string)
//   if actualVersion != expectedVersion {
//     t.Errorf("Did not update version correctly, got version: %v, expectedVersion: %v", actualVersion, expectedVersion)
//   }
//
//   configUpdates["version"] = version1
//   db.SetCurrentDatabaseConfig(configUpdates, false)
//
// 	db.Update()
//
//   config, _ = db.GetCurrentDatabaseConfig()
//   actualVersion = config["version"].(string)
//   if actualVersion != expectedVersion {
//     t.Errorf("Did not update version correctly, got version: %v, expectedVersion: %v", actualVersion, expectedVersion)
//   }
//
// 	var numTables int = 10
// 	var numRemovedTables int = 2
// 	tableList := GenerateTableList(numTables)
// 	removedTableList := make([]string, 0)
// 	for _, tableName := range tableList {
// 		_, tableErr := db.NewDefaultTable(tableName)
// 		if tableErr != nil {
// 			t.Error("Cannot make table: ", tableErr)
// 		}
// 	}
// 	db.Update()
// 	for i := 0; i < numRemovedTables; i++ {
// 		i := rand.Intn(numTables - i)
// 		tableName := tableList[i]
// 		_, removeErr := db.RemoveTable(tableName)
// 		if removeErr != nil {
// 			t.Error("Issue with removing table: ", removeErr)
// 		}
// 		removedTableList = append(removedTableList, tableName)
// 	}
// 	db.Update()
// 	for _, tableName := range removedTableList {
// 		table, tableErr := db.GetTable(tableName)
// 		if tableErr != nil {
// 			t.Error("Cannot make table: ", tableErr)
// 		}
// 		if table == nil {
// 			t.Error("Updating tables does not work correctly")
// 		}
// 	}
// }
//
//
//
func TestTableConfig(t *testing.T) {
	db, databaseErr := NewDefaultDatabase(databaseName, databaseID)
	if databaseErr != nil {
		t.Error("Cannot create database: ", databaseErr)
	}

	table, tableErr := db.NewDefaultTable(tableName)
	if tableErr != nil {
		t.Error("Cannot make table: ", tableErr)
	}
  configUpdatev2 := CopyMapNewVersion("0.0.2", testTableConfig)
	configUpdatev1 := CopyMapNewVersion("0.0.1", testTableConfig)
	blockNum1 := 1
//here

  table.SetCurrentTableConfig(configUpdatev2)
  expectedVersion := "0.0.2"
  config, _ := table.GetCurrentTableConfig()
  actualVersion := config["version"].(string)
  if actualVersion != expectedVersion {
    t.Errorf("Did not update version correctly, got version: %v, expectedVersion: %v", actualVersion, expectedVersion)
  }

	table.Update()

	//Check local
	//localConfig, _ := table.GetLocalTableConfig()
	//currentConfig, _ := table.GetCurrentTableConfig()


	var entries map[string]string = GenerateKeyValuePairs(keySize, valueSize, numEntries)
	for key, value := range entries {
		_, putErr := table.Put(key, value, false, args)
		if putErr != nil {
			t.Error("Cannot put key and value: ", putErr)
		}
	}

  configUpdatev2["blockNum"] = blockNum1
  table.SetCurrentTableConfig(configUpdatev2)
  expectedBlockNum := blockNum1

	table.Update()

  config, _ = table.GetCurrentTableConfig()
  actualBlockNum := config["blockNum"].(int)
  if actualBlockNum != expectedBlockNum {
    t.Errorf("Did not update blockNum correctly, got blockNum: %v, expected blockNum: %v", actualBlockNum, expectedBlockNum)
  }

	for key, value := range entries {
		resp, getErr := table.Get(key, false)
		if getErr != nil {
			t.Error("Cannot get value: ", getErr)
		}
		if resp == nil {
			t.Error("Key is nil should have value: ", value)
		}
	}

	config, _ = table.GetCurrentTableConfig()
	actualVersion = config["version"].(string)

	if actualVersion != expectedVersion {
		t.Errorf("Did not update version correctly, got version: %v, expectedVersion: %v", actualVersion, expectedVersion)
	}
	table.SetCurrentTableConfig(configUpdatev1)
	config, _ = table.GetCurrentTableConfig()
	actualVersion = config["version"].(string)
	if actualVersion != "0.0.1" {
		t.Errorf("Did not update the version correctly, got outdated version: %v, expectedVersion: %v", actualVersion, "0.0.1")
	}
	table.Update()

	newConfig, _ := table.GetCurrentTableConfig()
	actualVersion = newConfig["version"].(string)
	if actualVersion != expectedVersion {
		t.Errorf("Did not update the version correctly, got outdated version: %v, expectedVersion: %v", actualVersion, expectedVersion)
	}
}

func TestGet(t *testing.T) {
	db, databaseErr := NewDefaultDatabase(databaseName, databaseID)
	if databaseErr != nil {
		t.Error("Cannot create database: ", databaseErr)
	}

	table, tableErr := db.NewDefaultTable(tableName)
	if tableErr != nil {
		t.Error("Cannot make table: ", tableErr)
	}

	var entries map[string]string = GenerateKeyValuePairs(keySize, valueSize, numEntries)
	for key, _ := range entries {
		_, getErr := table.Get(key, false)
		if getErr != nil {
			t.Error("Cannot get value: ", getErr)
		}
	}
}



func TestPut(t *testing.T) {
	db, databaseErr := NewDefaultDatabase(databaseName, databaseID)
	if databaseErr != nil {
		t.Error("Cannot create database: ", databaseErr)
	}

	table, tableErr := db.NewDefaultTable(tableName)
	if tableErr != nil {
		t.Error("Cannot make table: ", tableErr)
	}

	var entries map[string]string = GenerateKeyValuePairs(keySize, valueSize, numEntries)
	for key, value := range entries {
		_, putErr := table.Put(key, value, false, args)
		if putErr != nil {
			t.Error("Cannot put key and value: ", putErr)
		}
	}

	for key, value := range entries {
		resp, getErr := table.Get(key, false)
		if getErr != nil {
			t.Error("Cannot get value: ", getErr)
		}
		if resp == nil {
			t.Error("Value is 0 should be: ", value)
		}
	}
}


// func TestRemove(t *testing.T) {
//   db, databaseErr := NewDefaultDatabase(databaseName, databaseID)
//   if databaseErr != nil {
//     t.Error("Cannot create database: ", databaseErr)
//   }
//
//   table, tableErr := db.NewDefaultTable(tableName)
//   if tableErr != nil {
//     t.Error("Cannot make table: ", tableErr)
//   }
//
//   var entries map[string]string = GenerateKeyValuePairs(keySize, valueSize, numEntries)
//   for key, _ := range entries {
//     _, removeErr := table.Remove(key, false)
//     if removeErr != nil {
//       t.Error("Cannot remove value from key: ", removeErr)
//     }
//   }
// }



// func (t *testing.T) TestBatch(db *ProximaDatabase, tableName string, entries []interface{}) {
//   // sizeValues := 300
//   // numEntries := 3000
//   // entries := generateKeyValuePairs(numEntries, 32, sizeValues)
//   // batchEntries := generateEntries(tableName, entries)
//   start := time.Now()
//   _, err := proximaClient.Batch(batchEntries, args)
//   if err != nil {
//     t.Error("Cannot batch values: ", err)
//   }
//   end := time.Now()
//   elapsed := end.Sub(start)
//   fmt.Println(elapsed)
// }
