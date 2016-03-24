package main

import "fmt"

func Echo(splitLine []string) {
	fmt.Println(GetTypeValue(splitLine[1]))
}
