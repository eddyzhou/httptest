package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func Getvalue(splitLine []string) {
	if BaseVariables["RESP"].(string) == "" {
		panic(fmt.Sprintf("RESP varable must be set before getvalue can be called. Line %v", CurrLineNum))
	}
	resp := BaseVariables["RESP"].(string)
	m, err := toMap(resp)
	if err != nil {
		panic(fmt.Sprintf("Convert json`%v` failed on line %v", resp, CurrLineNum))
	}
	cv := FindCustomVariable(splitLine[1])
	if cv == "" {
		panic(fmt.Sprintf("getvalue format error on Line %v", CurrLineNum))
	}
	toFind := splitLine[2]
	if val, ok := m[toFind]; ok {
		CustomVariables[cv] = val
	} else {
		panic(fmt.Sprintf("Getvalue error, either no such value `%v` in RESP. Line %v", toFind, CurrLineNum))
	}
}

func toMap(jsonStr string) (map[string]interface{}, error) {
	var v interface{}
	if err := json.Unmarshal([]byte(jsonStr), &v); err != nil {
		return nil, err
	}
	m := v.(map[string]interface{})
	return m, nil
}

func GetArray(splitLine []string) {
	if BaseVariables["RESP"].(string) == "" {
		panic(fmt.Sprintf("RESP varable must be set before getvalue can be called. Line %v", CurrLineNum))
	}
	resp := BaseVariables["RESP"].(string)
	m, err := toMap(resp)
	if err != nil {
		panic(fmt.Sprintf("Convert json`%v` failed on line %v", resp, CurrLineNum))
	}
	cv := FindCustomVariable(splitLine[1])
	if cv == "" {
		panic(fmt.Sprintf("getarray format error on Line %v", CurrLineNum))
	}
	toFind := splitLine[2]

	var results []interface{}
	for _, v := range m {
		if vv, ok := v.([]interface{}); ok {
			for _, u := range vv {
				d := u.(map[string]interface{})
				if val, ok := d[toFind]; ok {
					switch val.(type) {
					case float32, float64:
						val = int(reflect.ValueOf(val).Float())
					case int, int8, int16, int32, int64:
						val = int(reflect.ValueOf(val).Int())
					case uint, uint8, uint16, uint32, uint64:
						val = int(reflect.ValueOf(val).Uint())
					}
					results = append(results, val)
				}
			}
		}
	}

	if len(results) > 0 {
		CustomVariables[cv] = results
	} else {
		panic(fmt.Sprintf("Getvalue error, either no such values `%v` in RESP. Line %v", toFind, CurrLineNum))
	}
}
