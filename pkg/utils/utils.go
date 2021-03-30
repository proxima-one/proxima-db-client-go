package utils

import "fmt"

func ConvertMapTo(inputMap map[interface{}]interface{}) map[string]interface{} {
	var configMap map[string]interface{} = make(map[string]interface{})
	var key string
	for k, value := range inputMap {
		key = k.(string)
		valueType := fmt.Sprintf("%T", value)
		newValue := value
		if valueType == "map[interface  {}]interface {}" {
			//fmt.Println(value)
			var strMap map[string]interface{} = ConvertMapTo(value.(map[interface{}]interface{}))
			configMap[key] = strMap
			//fmt.Println(fmt.Sprintf("Value of map: %T", strMap))
		}
		if valueType == "[]interface {}" {
			newValue := make([]interface{}, len(value.([]interface{})))
			for i, v := range value.([]interface{}) {
				newV := v
				//fmt.Println(newV)
				if fmt.Sprintf("%T", v) == "map[interface {}]interface {}" {
					var strMap map[string]interface{} = ConvertMapTo(v.(map[interface{}]interface{}))
					newValue[i] = strMap

				} else {
					newValue[i] = newV
				}
			}
		}

		configMap[key] = newValue
	}
	return configMap
}
