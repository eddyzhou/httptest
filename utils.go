package main

import (
	"fmt"
	"regexp"
	"strings"
)

var CurrLineNum = 0
var CurrUrl string

var BaseVariables = map[string]interface{}{
	"HOST":     "",
	"RESP":     "",
	"POSTDATA": nil,
	"HEADER":   map[string]string{"User-Agent": "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:35.0) Gecko/20100101 Firefox/35.0"},
	"STATUS":   0,
}

var CustomVariables = make(map[string]interface{})
var varRe = regexp.MustCompile(".*[[:punct:]](.*)[[:punct:]].*")

var FnKeywords = map[string]func(splitLine []string){
	"req":      Req,
	"set":      Setvar,
	"echo":     Echo,
	"getvalue": Getvalue,
	"getarray": GetArray,
	"assert":   Assert,
}

func FindCustomVariable(s string) string {
	results := varRe.FindStringSubmatch(s)
	if len(results) > 1 {
		return results[1]
	}

	return ""
}

func GetTypeValue(s string) interface{} {
	if strings.HasPrefix(s, "$") {
		typ := s[1:]
		if value, ok := BaseVariables[typ]; ok {
			return value
		}
		if cv := FindCustomVariable(typ); cv != "" {
			if v, ok := CustomVariables[cv]; ok {
				return v
			}
		}
		panic(fmt.Sprintf("Error on line %v, variable %v is not a language variable.", CurrLineNum, typ))
	} else if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		data := s[1 : len(s)-1]
		arr := strings.Split(data, ",")
		r := make([]interface{}, len(arr))
		for i, x := range arr {
			r[i] = strings.TrimSpace(x)
		}
		return r
	}

	return s
}
