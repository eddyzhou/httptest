package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func Req(splitLine []string) {
	switch splitLine[1] {
	case "GET":
		get(splitLine)
	case "POST":
		post(splitLine)
	default:
		panic(fmt.Sprintf("Do error, no such request '%v' on line %v", splitLine[1], CurrLineNum))
	}
}

func get(splitLine []string) {
	if BaseVariables["HOST"].(string) == "" {
		panic(fmt.Sprintf("Error on line %v, HOST not set.", CurrLineNum))
	}

	CurrUrl = BaseVariables["HOST"].(string) + getQueryPath(splitLine[2])
	client := &http.Client{}
	req, err := http.NewRequest("GET", CurrUrl, nil)
	if err != nil {
		panic(fmt.Sprintf("Error on line %v get url %v failed: %v", CurrLineNum, CurrUrl, err))
	}
	setHeaders(req)

	res, err := client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("Error on line %v url %v get data failed: %v", CurrLineNum, CurrUrl, err))
	}
	body, err := ioutil.ReadAll(res.Body)
	resp := string(body)
	res.Body.Close()
	if err != nil {
		panic(fmt.Sprintf("Error on line %v url %v get data failed: %v", CurrLineNum, CurrUrl, err))
	}

	BaseVariables["RESP"] = resp
	BaseVariables["STATUS"] = res.StatusCode
}

func getQueryPath(s string) string {
	if cv := FindCustomVariable(s); cv != "" {
		if v, ok := CustomVariables[cv]; ok {
			toReplace := "${" + cv + "}"
			switch v.(type) {
			case int:
				return strings.Replace(s, toReplace, strconv.Itoa(v.(int)), -1)
			case string:
				return strings.Replace(s, toReplace, v.(string), -1)
			}
		}

	}
	return s
}

func setHeaders(req *http.Request) {
	if h, ok := BaseVariables["HEADER"]; ok && h != nil {
		headers := h.(map[string]string)
		for k, v := range headers {
			if v != "" {
				req.Header.Set(k, v)
			}
		}
	}
}

func post(splitLine []string) {
	if BaseVariables["HOST"].(string) == "" {
		panic(fmt.Sprintf("Error on line %v, HOST not set.", CurrLineNum))
	}
	if _, ok := BaseVariables["POSTDATA"]; !ok {
		panic(fmt.Sprintf("POSTDATA must be set to do POST request, line %v", CurrLineNum))
	}

	currUrl := BaseVariables["HOST"].(string) + getQueryPath(splitLine[2])
	client := &http.Client{}
	values := url.Values{}
	for k, v := range BaseVariables["POSTDATA"].(map[string]string) {
		values.Set(k, v)
	}
	CurrUrl = currUrl + "?" + values.Encode()
	req, err := http.NewRequest("POST", CurrUrl, nil)
	if err != nil {
		panic(fmt.Sprintf("Error on line %v url %v does not exist.", CurrLineNum, CurrUrl))
	}
	setHeaders(req)

	res, err := client.Do(req)
	if err != nil {
		panic(fmt.Sprintf("Error on line %v url %v get data failed: %v", CurrLineNum, CurrUrl, err))
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	resp := string(body)
	if err != nil {
		panic(fmt.Sprintf("Error on line %v url %v get data failed.", CurrLineNum, CurrUrl))
	}

	BaseVariables["RESP"] = resp
	BaseVariables["STATUS"] = res.StatusCode
}
