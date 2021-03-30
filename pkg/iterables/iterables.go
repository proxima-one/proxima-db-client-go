package iterables

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
)

func Search(objs interface{}, filterSet interface{}, orderBy interface{}, orderDirection interface{}, first int, last int, args interface{}) (interface{}, error) {
	var newObjs []interface{}

	if IsSlice(objs) {
		newObjs = objs.([]interface{})
	} else {
		strObj := fmt.Sprintf("%v", objs)
		err := json.Unmarshal([]byte(strObj), &newObjs)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
	}
	//maps
	limit := len(newObjs)
	var noLimit bool
	var isFirst bool

	if orderBy != nil {
		orderByFn, _ := createOrderByFn(orderBy.(string), orderDirection, args)
		resp, _ := Order(newObjs, orderByFn)
		newObjs = resp.([]interface{})

		orderLimit := last
		if first > last {
			isFirst = true
			orderLimit = first
		}
		if limit > orderLimit {
			limit = orderLimit
		}
		noLimit = (len(newObjs) == limit) || (limit <= 0)
	}

	if filterSet != nil {
		filterFn, _ := createFilterFn(filterSet, args)
		resp, _ := Filter(newObjs, filterFn)
		newObjs = resp.([]interface{})
	}

	if !noLimit && isFirst {
		return newObjs[0:limit], nil
	}

	if !noLimit && !isFirst {
		return newObjs[len(newObjs)-limit : len(newObjs)], nil
	}

	return newObjs, nil
}

func IsSlice(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Slice || reflect.TypeOf(v).Kind() == reflect.Array
}

func IsFunc(fn interface{}) bool {
	return reflect.TypeOf(fn).Kind() == reflect.Func
}

func Map(objs interface{}, mapFn func(interface{}) interface{}) (interface{}, error) {
	var newObjs []interface{}
	if !IsSlice(objs) {
		return newObjs, nil
	}
	if !IsFunc(mapFn) {
		return newObjs, nil
	}
	for _, obj := range objs.([]interface{}) {
		newObjs = append(newObjs, mapFn(obj))
	}
	return newObjs, nil
}

func Order(objs interface{}, orderByFn func(interface{}, interface{}) bool) (interface{}, error) {
	if !IsSlice(objs) {
		return objs, nil
	}
	var newObjs []interface{} = objs.([]interface{})
	if len(newObjs) > 0 {
		sort.Slice(newObjs, func(i, j int) bool {
			return orderByFn(newObjs[i], newObjs[j])
		})
	}
	return newObjs, nil
}

func createOrderByFn(sortFieldName string, orderDirection interface{}, args interface{}) (func(interface{}, interface{}) bool, error) {
	mapFn := func(e interface{}) interface{} { return e } //default
	if IsFunc(args) {
		mapFn = args.(func(interface{}) interface{})
	}

	orderByFn := func(entityI, entityJ interface{}) bool {
		i := reflect.ValueOf(&entityI).FieldByName(sortFieldName)
		j := reflect.ValueOf(&entityJ).FieldByName(sortFieldName)
		var expression string = ">"
		if orderDirection != nil {
			expression = "<"
		}
		return FilterExpression(i, expression, j)
	}
	return func(entityI, entityJ interface{}) bool {
		return orderByFn(mapFn(entityI), mapFn(entityJ))
	}, nil
}

func Filter(objs interface{}, filter func(interface{}) bool) (interface{}, error) {
	if !IsSlice(objs) {
		return objs, nil
	}
	var newObjs []interface{}
	for _, obj := range objs.([]interface{}) {
		if filter(obj) {
			newObjs = append(newObjs, obj)
		}
	}
	return newObjs, nil
}

func createFilterFn(filterList interface{}, args interface{}) (func(interface{}) bool, error) {
	mapFn := func(e interface{}) interface{} { return e } //default
	if IsFunc(args) {
		mapFn = args.(func(interface{}) interface{})
	}
	filterSet, err := LoadSearchFilterSet(filterList)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return func(obj interface{}) bool {
		newObj := mapFn(obj)
		var matchesFilter bool
		isSelecting := filterSet.isSelecting
		for _, filter := range filterSet.filters {
			result := reflect.ValueOf(&newObj).Elem().FieldByName(filter.varName)
			matchesFilter = FilterExpression(result, filter.expression, filter.value) //compare
			if (!matchesFilter && isSelecting) || (matchesFilter && !isSelecting) {
				return false
			}
		}
		return false
	}, nil
}

func LoadSearchFilterSet(filterList interface{}) (SearchFilterSet, error) {
	var filterSet SearchFilterSet = SearchFilterSet{isSelecting: true, filters: make([]SearchFilter, 0)}

	//if is string, then get the value of it, such that it is []interface{}, reflect, and set to new val

	var newFilters []interface{}

	if IsSlice(filterList) {
		newFilters = filterList.([]interface{})
	} else {
		strObj := fmt.Sprintf("%v", filterList)
		err := json.Unmarshal([]byte(strObj), &newFilters)
		if err != nil {
			fmt.Println(err.Error())
			return filterSet, err
		}
	}

	//processFilterList, check if is instance of []interface{},

	for _, filter := range newFilters {
		searchFilter, err := LoadSearchFilter(filter)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println(filter)
		} else {
			filterSet.filters = append(filterSet.filters, searchFilter)
		}
	}
	return filterSet, nil
}

func LoadSearchFilter(filter interface{}) (SearchFilter, error) {
	f, ok := filter.(SearchFilter)
	if ok {
		return f, nil
	}
	strFilter := fmt.Sprintf("%v", filter)
	var filterMap map[string]interface{} //map
	err := json.Unmarshal([]byte(strFilter), &filterMap)
	if err != nil {
		fmt.Println(err.Error())
		return SearchFilter{}, err
	} else {
		var varName string = filterMap["varName"].(string)
		var expression string = filterMap["expression"].(string)
		//get the fields of the searchFilter
		//find name (if only two values, then)
		//expression
		//value

		return SearchFilter{varName: varName, expression: expression, value: filterMap["value"]}, nil
	}
}

type SearchFilterSet struct {
	filters     []SearchFilter
	isSelecting bool
}

type SearchFilter struct {
	varName    string
	expression string
	value      interface{}
}

//marshall
//encoding

//sort by

func FilterExpression(result interface{}, exp string, expected interface{}) bool {
	//cast correctly
	//cast result
	//cast expected
	var isOrdered bool
	//diff/compare
	//var r interface{}
	//var e interface{}
	var diff bool

	switch result.(type) {
	case int:
		//set
		r := reflect.ValueOf(result).Int()
		e := reflect.ValueOf(expected).Int()
		diff = (r > e)
		isOrdered = true
		break
	case float64:
		r := reflect.ValueOf(result).Float()
		e := reflect.ValueOf(expected).Float()
		diff = (r > e)
		isOrdered = true
		break
	case string:
		r := reflect.ValueOf(result).String()
		e := reflect.ValueOf(expected).String()
		diff = (r > e)
		isOrdered = true
		break
	default:
		diff = false
		break
	}

	if !isOrdered {
		return false
	}

	//case string
	switch exp {
	case "=":
		return result == expected
	case "!=":
		return result != expected
	case ">":
		//get types
		//reflect.ValueOf(result)
		return diff
	case ">=":
		return diff || (result == expected)
	case "<":
		return !(diff || (result == expected))
	case "<=":
		return !diff
	default:
		return false
	}
	return false
}
