package main

import (
	"fmt"
	"strings"
)

func Setvar(splitLine []string) {
	keyword := splitLine[1]
	switch keyword {
	case "HOST":
		basicSet(splitLine)
	case "POSTDATA":
		postData(splitLine)
	case "HEADER":
		headerSet(splitLine)
	default:
		if cv := FindCustomVariable(keyword); cv != "" {
			CustomVariables[cv] = GetTypeValue(strings.Join(splitLine[2:], " "))
		} else {
			panic(fmt.Sprintf("Set error, no such key '%v' on line %v", keyword, CurrLineNum))
		}

	}
}

func basicSet(splitLine []string) {
	keyword := splitLine[1]
	BaseVariables[keyword] = GetTypeValue(strings.Join(splitLine[2:], " "))
}

func headerSet(splitLine []string) {
	line := strings.Join(splitLine[2:], " ")
	data := strings.Split(line, ",")
	headers := BaseVariables["HEADER"].(map[string]string)
	for _, x := range data {
		parts := strings.Split(x, ":")
		key := parts[0]
		val := GetTypeValue(strings.Join(parts[1:], " ")).(string)
		if val == "_" {
			delete(headers, key)
		} else {
			headers[key] = strings.TrimSpace(val)
		}
	}
}

func postData(splitLine []string) {
	line := strings.Join(splitLine[2:], " ")
	data := strings.Split(line, ",")
	m := make(map[string]string)
	for _, x := range data {
		parts := strings.Split(x, "=")
		m[parts[0]] = GetTypeValue(strings.Join(parts[1:], " ")).(string)
	}
	BaseVariables["POSTDATA"] = m
}
