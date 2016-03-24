package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Assert(splitLine []string) {
	keyword := splitLine[1]
	switch keyword {
	case "STATUS":
		assertStatus(splitLine)
	case "JSON":
		assertJson(splitLine)
	case "RESP":
		assertResp(splitLine)
	default:
		if cv := FindCustomVariable(keyword); cv != "" {
			if v, ok := CustomVariables[cv]; ok {
				assertValue(splitLine, v)
			} else {
				panic(fmt.Sprintf("Assert error, variable'%v' not assign value on line %v", cv, CurrLineNum))
			}
		} else {
			panic(fmt.Sprintf("Assert error, no such key '%v' on line %v", keyword, CurrLineNum))
		}
	}
}

func assertValue(splitLine []string, value interface{}) {
	except := strings.Join(splitLine[2:], " ")
	switch value.(type) {
	case string:
		if except != value {
			panic(fmt.Sprintf("Assert error, except %v, but %v on line %v", except, value, CurrLineNum))
		}
	case int:
		if except != strconv.Itoa(value.(int)) {
			panic(fmt.Sprintf("Assert error, except %v, but %v on line %v", except, value, CurrLineNum))
		}
	case float32:
		ex, _ := strconv.ParseFloat(except, 32)
		if ex != value {
			panic(fmt.Sprintf("Assert error, except %v, but %v on line %v", except, value, CurrLineNum))
		}
	case float64:
		ex, _ := strconv.ParseFloat(except, 64)
		if ex != value {
			panic(fmt.Sprintf("Assert error, except %v, but %v on line %v", except, value, CurrLineNum))
		}
	default:
		panic(fmt.Sprintf("Assert error, except %v, but %v on line %v", except, value, CurrLineNum))
	}
}

func assertStatus(splitLine []string) {
	except := GetTypeValue(strings.Join(splitLine[2:], " "))
	exceptStatus, err := strconv.Atoi(except.(string))
	if err != nil {
		panic(fmt.Sprintf("Assert error, status %v not int type on line %v", except, CurrLineNum))
	}
	status := BaseVariables["STATUS"].(int)
	if status != exceptStatus {
		panic(fmt.Sprintf("Assert error, except %v, but %v on line %v", except, status, CurrLineNum))
	}
}

func assertJson(splitLine []string) {
	except := GetTypeValue(strings.Join(splitLine[2:], " "))
	resp := BaseVariables["RESP"].(string)
	em, err := toMap(except.(string))
	if err != nil {
		panic(fmt.Sprintf("Assert error, convert json`%v` failed on line %v", except, CurrLineNum))
	}
	rm, err := toMap(resp)
	if err != nil {
		panic(fmt.Sprintf("Assert error, resp`%v` convert json failed on line %v", resp, CurrLineNum))
	}
	if !reflect.DeepEqual(em, rm) {
		panic(fmt.Sprintf("Assert error, except %v, but '%v' on line '%v'", except, resp, CurrLineNum))
	}

}

func assertResp(splitLine []string) {
	except := GetTypeValue(strings.Join(splitLine[2:], " "))
	except = strings.TrimSpace(except.(string))
	resp := BaseVariables["RESP"].(string)
	resp = strings.TrimSpace(resp)
	if resp != except {
		panic(fmt.Sprintf("Assert error, except %v, but '%v' on line '%v'", except, resp, CurrLineNum))
	}
}
