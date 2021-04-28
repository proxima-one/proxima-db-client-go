package proxima_db_client_go

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"

	proxima_database "github.com/proxima-one/proxima-db-client-go/pkg/database"

	//"github.com/pkg/errors"
	//yaml "gopkg.in/yaml.v2"
	json "encoding/json"

	yaml "github.com/ghodss/yaml"
)

var databaseName string = "DefaultDatabaseName"
var databaseID string = "DefaultDatabaseID"
var tableName string = "DefaultTableName"
var databaseConfigFile string = "./helpers/db-config.yaml"

var valueSize int = 300
var numEntries int = 200
var numBatches int = 3
var keySize int = 32
var args map[string]interface{} = map[string]interface{}{"prove": false}
var testTableConfig map[string]interface{} = map[string]interface{}{"name": "DPoolLists",
	"id":              "DefaultDB-DPoolLists",
	"version":         "0.0.0",
	"blockNum":        0,
	"header":          "Root",
	"compression":     "36h",
	"batching":        "500ms",
	"sleep":           "10m",
	"cacheExpiration": "5m"}

var testDatabaseConfig map[string]interface{} = map[string]interface{}{
	"name":    "DefaultDB",
	"id":      "DefaultDB",
	"owner":   "None",
	"version": "0.0.0",
	"config": map[string]interface{}{
		"sleep":           "5m",
		"compression":     "36h",
		"batching":        "500ms",
		"cacheExpiration": "10s",
	},
	"tables": []interface{}{testTableConfig},
}

func getDBConfig(configPath string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return make(map[string]interface{}), nil
	}
	jsonData, _ := yaml.YAMLToJSON([]byte(data))
	var configMap map[string]interface{}
	err = json.Unmarshal(jsonData, &configMap)
	if err != nil {
		fmt.Printf("err: %+v\n", err)
		return make(map[string]interface{}), nil
	}
	return configMap, nil
}

func CopyMapNewVersion(newVersion string, originalMap map[string]interface{}) map[string]interface{} {
	CopiedMap := make(map[string]interface{})

	/* Copy Content from Map1 to Map2*/
	for index, element := range originalMap {
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
		"network": networkTableConfig,
		"node":    nodeTableConfig,
		"local":   localTableConfig,
		"current": currentTableConfig,
	}
	testConfigMap["table"] = tableTestConfigMap

	///networkDatabaseConfig := CopyMapNewVersion("0.0.1", testDatabaseConfig)
	//nodeDatabaseConfig := CopyMapNewVersion("0.0.5", testDatabaseConfig)
	localDatabaseConfig := CopyMapNewVersion("0.0.2", testDatabaseConfig)
	currentDatabaseConfig := CopyMapNewVersion("0.0.1", testDatabaseConfig)

	databaseTestConfigMap := map[string]interface{}{
		//	"network": networkDatabaseConfig,
		//	"node":    nodeDatabaseConfig,
		"local":   localDatabaseConfig,
		"current": currentDatabaseConfig,
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
		if latestName == "" {
			t.Errorf("Issues with checking the latest. Expected %v, Actual %v", "latest", latestName)
		}
	}
}

//Test Scan

//Test Put (how fast with checkout/commit )

//Test Compact

//Test Search

//Test Stat

func TestDatabaseCreation(t *testing.T) {
	db, databaseErr := NewDefaultDatabase(databaseName, databaseID)
	if databaseErr != nil {
		t.Error("Cannot create database: ", databaseErr)
	}

	dbConfig := testDatabaseConfig
	dbConfig["version"] = "0.0.1"

	db.SetCurrentDatabaseConfig(dbConfig, true)
	actualConfig, _ := db.GetCurrentDatabaseConfig()
	if actualConfig == nil || actualConfig["version"].(string) != dbConfig["version"].(string) {
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

func TestDatabaseLoad(t *testing.T) {

	dbConfig, err := getDBConfig(databaseConfigFile)
	if err != nil || dbConfig == nil {
		t.Error("Issues loading the configuration files.", dbConfig)
	}

	db, dbErr := proxima_database.LoadProximaDatabase(dbConfig)
	if dbErr != nil {
		t.Error("Issues with loading database from config: ", dbErr)
	}
	actualConfig, err := db.GetCurrentDatabaseConfig()
	if err != nil {
		t.Error("Issue with the loading the current database config", err)
	}

	if actualConfig["version"].(string) != dbConfig["version"].(string) {
		t.Errorf("Issues with loading database config. expected: %v but got: %v", dbConfig["version"].(string), actualConfig["version"].(string))
	}

	//check the tables
	expectedTables := dbConfig["tables"].([]interface{})
	for _, tConfig := range expectedTables {
		//get table name
		var tableConfig map[string]interface{} = tConfig.(map[string]interface{})
		//check the name of the table
		var tableName string = tableConfig["name"].(string)
		//check the config of the table
		//fmt.Println("Tables that are expected")
		//	fmt.Println(tableName)
		actualTable, err := db.GetTable(tableName)
		if err != nil || actualTable == nil {
			t.Error("Issue with getting the table after loading: ", err)
		}
		//	fmt.Println(tableConfig)
	}

	//check the tables
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

//test scan, and query

func TestBasicDatabase(t *testing.T) {
	db, databaseErr := NewDefaultDatabase(databaseName, databaseID)
	if databaseErr != nil {
		t.Error("Cannot create database: ", databaseErr)
	}
	var numTables int = 10
	var numRemovedTables int = 2

	tableList := GenerateTableList(numTables)

	for _, tableName := range tableList {
		newTable, tableErr := db.NewDefaultTable(tableName)
		if tableErr != nil {
			t.Error("Cannot make table: ", tableErr)
		}

		if newTable == nil {
			t.Error("Issue creating the table.", tableName)
		}
		db.AddTable(tableName, newTable)
	}
	for _, tableName := range tableList {

		table, tableErr := db.GetTable(tableName)
		if tableErr != nil {
			t.Error("Cannot make table: ", tableErr)
		}

		if table == nil {
			t.Error("Error the table is nil", tableName)
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
	var updated bool
	table, tableErr := db.NewDefaultTable(tableName)
	if tableErr != nil {
		t.Error("Cannot make table: ", tableErr)
	}
	table.Open()
	configUpdatev2 := CopyMapNewVersion("0.0.2", testTableConfig)
	configUpdatev1 := CopyMapNewVersion("0.0.1", testTableConfig)
	blockNum1 := 1
	//here

	table.SetCurrentTableConfig(configUpdatev2)
	expectedVersion := "0.0.2"
	config, err := table.GetCurrentTableConfig()
	if err != nil {
		t.Errorf("Error with getting current table config: %v", err)
	}
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

	updated, err = table.Update()
	if err != nil || updated == false {
		t.Errorf("Error with updating table config: %v", err)
	}

	config, err = table.GetCurrentTableConfig()

	if err != nil {
		t.Errorf("Error with getting current table config: %v", err)
	}

	actualBlockNum := config["blockNum"].(int)
	//fmt.Println(config)
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

	updated, err = table.Update()
	if err != nil || updated == false {
		t.Errorf("Error with updating table config: %v", err)
	}

	newConfig, _ := table.GetCurrentTableConfig()
	actualVersion = newConfig["version"].(string)
	if actualVersion != expectedVersion {
		t.Errorf("Did not update the version correctly, got outdated version: %v, expectedVersion: %v", actualVersion, expectedVersion)
	}

	updated, err = db.Update()
	if err != nil || updated == false {

		t.Errorf("Error with updating db config: %+v", err)
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
	err := table.Open()
	if err != nil {
		t.Error("Cannot open table: ", err)
	}
	num := 200
	var entries map[string]string = GenerateKeyValuePairs(keySize, valueSize, num)
	start := time.Now()
	for key, value := range entries {
		resp, putErr := table.Put(key, value, false, args)

		if putErr != nil {
			fmt.Println(resp)
			t.Error("Cannot put key and value: ", putErr)
		}
	}
	endT := time.Now()
	elapsed := endT.Sub(start)
	fmt.Println(elapsed)
	fmt.Println(num)

	start = time.Now()
	//range queries
	for key, value := range entries {
		resp, getErr := table.Get(key, false)
		if getErr != nil {
			t.Error("Cannot get value: ", getErr)
		}
		if resp == nil {
			t.Error("Value is 0 should be: ", value)
		}
	}
	endT = time.Now()
	elapsed = endT.Sub(start)
	fmt.Println(elapsed)
	fmt.Println(num)
}

func TestScan(t *testing.T) {
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

	var num int = 10

	start := time.Now()
	for i := 0; i < num; i++ {
		resp, scanErr := table.Scan(int(0), int(10+i), int(100), false, make(map[string]interface{}))
		//fmt.Println(resp)
		if scanErr != nil {
			fmt.Println(resp)
			t.Error("Issue with scanning table: ", scanErr)
		}
	}
	endT := time.Now()
	elapsed := endT.Sub(start)
	fmt.Println(elapsed)
	fmt.Println(num)
}

func TestQuery(t *testing.T) {
	db, databaseErr := NewDefaultDatabase(databaseName, databaseID)
	if databaseErr != nil {
		t.Error("Cannot create database: ", databaseErr)
	}
	var tableName string = "DPools"
	table, tableErr := db.NewDefaultTable(tableName)
	if tableErr != nil {
		t.Error("Cannot make table: ", tableErr)
	}

	//generate random structures
	var entries map[string]string = GenerateKeyValuePairs(keySize, valueSize, numEntries)

	//entity in entities
	for key, value := range entries {
		_, putErr := table.Put(key, value, false, args)
		if putErr != nil {
			t.Error("Cannot put key and value: ", putErr)
		}
	}
	// var queries []string = GenerateRandomSearchQueryText()
	// //generate random queries
	// for i := 0; i < 10; i++ {
	// 	resp, queryErr := table.Query(queries[i], false)
	// 	fmt.Println(resp)
	// 	if queryErr != nil {
	// 		fmt.Println(resp)
	// 		t.Error("Issue with scanning table: ", queryErr)
	// 	}
	// }
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
