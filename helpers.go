package proxima_db_client_go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// func RandomString(size int) string {
// 	bytes := make([]byte, size)
// 	rand.Read(bytes)
// 	return string(bytes)
// }

func makeQueryString(queryString string, variables map[string]interface{}) string {
	var encodedValue []byte
	for name, value := range variables {
		valueType := fmt.Sprintf("%T", value)
		//fmt.Println(name)
		//fmt.Println(value)
		if name == "input" {
			encodedValue, _ = json.Marshal(value)
			queryString = strings.ReplaceAll(queryString, "$"+name, string(encodedValue))
			var inputVars map[string]interface{} = value.(map[string]interface{})
			for na, _ := range inputVars {
				varName := fmt.Sprintf("\"%s\"", na)
				queryString = strings.ReplaceAll(queryString, varName, na)
			}
		} else if name == "queryText" {
			// QUESTION: Query
			varValue := fmt.Sprintf(`%q`, value)
			queryString = strings.ReplaceAll(queryString, "$"+name, varValue)
			//fmt.Println(queryString)
		} else if valueType == "string" || valueType == "String" {
			varValue := fmt.Sprintf("\"%s\"", fmt.Sprint(value))
			//type
			queryString = strings.ReplaceAll(queryString, "$"+name, varValue)
		} else {
			queryString = strings.ReplaceAll(queryString, "$"+name, fmt.Sprint(value))
		}
	}
	//every variable in vars, replace with the value of query with string value (name + $)
	return queryString
}

// func NewDatabase(name string) (*ProximaDB) {
//   ip := "0.0.0.0"
//   port := "50051"
//   proximaClient := NewProximaDB(ip, port)
//   proximaClient.Open(name)
//   return proximaClient
// }

// func MainSetup(name string, numEntries int, sizeValues int, prove bool) (*ProximaDB, map[string]string, map[string]interface{}) {
//   proximaClient := NewDatabase(name)
//   entries := generateKeyValuePairs(numEntries, 32, sizeValues)
//   args := make(map[string]interface{})
//   args["prove"] = prove
//   return proximaClient, entries, args
// }

func GenerateTableList(num int) []string {
	tableList := make([]string, 0)
	for i := 0; i < num; i++ {
		tableName := RandomString(32)
		tableList = append(tableList, tableName)
	}
	return tableList
}

func GenerateKeyValuePairs(keySize int, valSize int, num int) map[string]string {
	mapping := make(map[string]string)
	for i := 0; i < num; i++ {
		key := RandomString(keySize)
		value := RandomString(valSize)
		mapping[key] = value
	}
	return mapping
}

func GenerateDatabaseGetOperations(tableName string, pairs map[string]string) []interface{} {
	entries := make([]interface{}, 0)
	for key, value := range pairs {
		entry := map[string]interface{}{"key": key, "value": value, "table": tableName, "prove": false}
		entries = append(entries, entry)
	}
	return entries
}

func GenerateDatabasePutOperations(tableName string, pairs map[string]string) []interface{} {
	entries := make([]interface{}, 0)
	for key, value := range pairs {
		entry := map[string]interface{}{"key": key, "value": value, "table": tableName, "prove": false}
		entries = append(entries, entry)
	}
	return entries
}

func GenerateDatabaseRemoveOperations(tableName string, pairs map[string]string) []interface{} {
	entries := make([]interface{}, 0)
	for key, value := range pairs {
		entry := map[string]interface{}{"key": key, "value": value, "table": tableName, "prove": false}
		entries = append(entries, entry)
	}
	return entries
}

// func NewDatabase(name string) (*ProximaDB) {
//   ip := "0.0.0.0"
//   port := "50051"
//   proximaClient := (ip, port)
//   proximaClient.Open(name)
//   return proximaClient
// }

// func generateEntries(num int, sizeValues int, prove bool) (map[string]interface{}){
//   mapping := make(map[string]string);
//   for i := 0; i < num; i++ {
//     key := randomString(32)
//     value := randomString(valSize)
//     mapping[key] = value
//   }
//   return mapping
// }
//
// func generateKeyValuePairs(num int, keySize int, valSize int) (map[string]string){
//   mapping := make(map[string]string)
//   for i := 0; i < num; i++ {
//     key := randomString(keySize)
//     value := randomString(valSize)
//     mapping[key] = value
//   }
//   return mapping
// }
//
// func generateEntries(name string, pairs map[string]string) ([]interface{}){
//   entries := make([]interface{}, 0)
//   for key, value := range pairs {
//     entry:= map[string]interface{}{"key": key, "value": value, "table": name, "prove": false}
//     entries = append(entries, entry)
//   }
//   return entries
// }
//
// func makeBatches(batchSize, total int) ([]int) {
//   batches := make([]int, 0)
//   num := 0
//   for total > (num + 1) {
//     num += batchSize
//     batches = append(batches, num)
//   }
//   return append(batches, total)
// }

// func GenerateDatabaseQueryOperations(tableName string, queries []string) ([]interface{}){
//   entries := make([]interface{}, 0)
//   for key, value := range queries {
//     entry:= map[string]interface{}{"queryText": queryText, "table": tableName, "prove": false}
//     entries = append(entries, entry)
//   }
//   return entries
// }
// func (entityTest *EntityTestCase) generateSearchTest(queryString string, entities []interface{}) *GQLTest {
// 	schema := entityTest.schema
// 	operation := entityTest.operations["search"].(map[string]interface{})
// 	//queryStr := operation["type"]
// 	var operationName string = operation["type"].(string)
// 	//entity := entityMap["entityInput"].(map[string]interface{})
// 	vars := make(map[string]interface{})
// 	//entityTest.entity == map[string]interface{}
// 	//fmt.Println(entityTest.entity)
// 	queryText, _ := GenerateRandomSearchQueryText(entityTest.entity)
//
// 	vars["queryText"] = queryText
// 	vars["prove"] = false
// 	expectedResult := "[]"
// 	queryString = makeQueryString(queryString, vars)
// 	expectedErrors := gqlerror.List{}
// 	return NewGQLTest(schema, queryString, operationName, vars, expectedResult, expectedErrors)
// }

func GenerateRandomSearchQueryText(entityMap map[string]interface{}) (string, error) {
	//for  Float, int,
	//generate  filters
	//randomEntity := make(map[string]interface{})
	var filterExpressions = []string{">", ">=", "<", "<="}
	var name string
	var varType string
	//list
	filters := make([]interface{}, 0)
	var entityVar map[string]interface{}

	for _, eVar := range entityMap {
		entityVar = eVar.(map[string]interface{})
		varType = entityVar["type"].(string)

		if varType == "Int" || varType == "int" || varType == "Float" {

			name = entityVar["name"].(string)
			value, _ := GenerateRandomOfType(varType)
			//fmt.Println(randomVar)
			if value != nil && (rand.Intn(4) > 2) {
				filterMap := make(map[string]interface{})
				filterMap["name"] = name
				filterMap["value"] = value
				filterExpressionIndex := rand.Intn(len(filterExpressions))
				filterMap["expression"] = filterExpressions[filterExpressionIndex]
				filters = append(filters, filterMap)
				break
			}
		}
		//fmt.Println(name)
		//fmt.Println(varType)
		//fmt.Println(randomEntity[name])
		// 		c.MustPost(
		// 	`mutation($id: ID!, $text: String!) { updateTodo(id: $id, changes:{text:$text}) { text } }`,
		// 	&resp,
		// 	client.Var("id", 5),
		// 	client.Var("text", "Very important"),
		// )
	}
	filterString, err := JSONMarshal(filters)
	if err != nil {
		fmt.Println(filterString)
		fmt.Println(err)
		return "", err
	}
	//convert to JSON string...
	//fmt.Println(string(filterString))
	return string(filterString), nil
}

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func GenerateRandomStruct(entityStruct map[string]interface{}) (map[string]interface{}, error) {
	randomEntity := make(map[string]interface{})

	var name string
	var varType string
	var entityVar map[string]interface{}
	for _, eVar := range entityStruct {
		entityVar = eVar.(map[string]interface{})
		name = entityVar["name"].(string)
		varType = entityVar["type"].(string)
		randomType, _ := GenerateRandomOfType(varType)
		if randomType != nil {
			randomEntity[name] = randomType
		}
		//fmt.Println(name)
		//fmt.Println(varType)
		//fmt.Println(randomEntity[name])
		// 		c.MustPost(
		// 	`mutation($id: ID!, $text: String!) { updateTodo(id: $id, changes:{text:$text}) { text } }`,
		// 	&resp,
		// 	client.Var("id", 5),
		// 	client.Var("text", "Very important"),
		// )
	}
	return randomEntity, nil
}

func RandomString(size int) string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, size)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateRandomOfType(varType string) (interface{}, error) {
	rand.Seed(time.Now().UnixNano())
	switch varType {
	case "String":
		return RandomString(32), nil
	case "Float":
		//range
		val := float64(20.0)

		return float64(val), nil
	case "ID":
		return RandomString(32), nil
	case "Int":
		//range

		return rand.Intn(3), nil
	case "Boolean":
		return (rand.Intn(2) != 0), nil
	default:
		return nil, nil
	}
}

func formatJSON(data []byte) ([]byte, error) {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	formatted, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return formatted, nil
}
