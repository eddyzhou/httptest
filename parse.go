package main

import (
	"fmt"
	"io/ioutil"
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type RunInfo struct {
	loopDepth      int
	rootLoopLength int
	loopCode       []Line
	line           string
	lineNumber     int
	splitLine      []string
	forSetStack    setStack
	forEachVar     string
}

type setStack [][]interface{}

func (s *setStack) Push(v []interface{}) {
	(*s) = append(*s, v)
}

func (s *setStack) Pop() []interface{} {
	res := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return res
}

type Line struct {
	lineNumber int
	line       string
}

func EnterREPL() {
	runInfo := new(RunInfo)
	input := bufio.NewScanner(os.Stdin)
	CurrLineNum = 1
	fmt.Print(CurrLineNum, "> ")
	for input.Scan() {
		line := strings.TrimSpace(input.Text())
		if line == "quit" || line == "q" {
			return
		}

		runInfo.line = line
		runInfo.lineNumber = CurrLineNum
		parse(runInfo)
		CurrLineNum++
		fmt.Print(CurrLineNum, "> ")
	}
}

func ParseFile(file string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse script file failed: %v", err)
		os.Exit(1)
	}
	var codes []Line
	for i, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		codes = append(codes, Line{i+1, line})
	}

	defer func() {
		if r := recover(); r != nil {
			log.Println("http test failed. current url: ", CurrUrl)
			log.Println(r)
			os.Exit(1)
		}
	}()
	parseLoop(codes, 1)
}

func parseLoop(codes []Line, loop int) {
	for i := 0; i < loop; i++ {
		runInfo := new(RunInfo)
		for _, l := range codes {
			runInfo.line = l.line
			runInfo.lineNumber = l.lineNumber
			CurrLineNum = l.lineNumber
			parse(runInfo)
		}
	}
}

func parse(runInfo *RunInfo) {
	line := removeComments(runInfo.line)
	if line == "" {
		return
	}
	runInfo.line = line
	runInfo.splitLine = strings.Split(runInfo.line, " ")

	if len(runInfo.splitLine) < 1 || runInfo.line == "\n" {
		return
	}

	firstWord := runInfo.splitLine[0]
	switch firstWord {
	case "loop":
		loop(runInfo)
	case "endloop":
		endLoop(runInfo)
	case "for":
		foreach(runInfo)
	case "endfor":
		endFor(runInfo)
	default:
		if runInfo.loopDepth != 0 {
			runInfo.loopCode = append(runInfo.loopCode, Line{runInfo.lineNumber, runInfo.line})
		} else {
			eval(runInfo.splitLine)
		}
	}
}

func foreach(runInfo *RunInfo) {
	runInfo.loopDepth += 1
	if runInfo.loopDepth == 1 {
		if len(runInfo.splitLine) != 4 {
			panic(fmt.Sprintf("For command format Error on line %v", CurrLineNum))
		}
		setVar := GetTypeValue(runInfo.splitLine[3])
		arr, ok := setVar.([]interface{})
		if !ok {
			panic(fmt.Sprintf("For command Error: not loop a collection on line %v", CurrLineNum))
		}
		each := FindCustomVariable(runInfo.splitLine[1])
		if len(each) == 0 {
			panic(fmt.Sprintf("For command Error: format err on line %v", CurrLineNum))
		}
		runInfo.forSetStack.Push(arr)
		runInfo.forEachVar = each
	} else {
		runInfo.loopCode = append(runInfo.loopCode, Line{runInfo.lineNumber, runInfo.line})
	}
}

func endFor(runInfo *RunInfo) {
	if runInfo.loopDepth == 0 {
		panic(fmt.Sprintf("Loop error on line %v. endFor has no corresponding loop to begin from.", CurrLineNum))
	}
	runInfo.loopDepth -= 1
	if runInfo.loopDepth == 0 {
		parseFor(runInfo)
		runInfo.loopCode = make([]Line, 50)
	} else {
		runInfo.loopCode = append(runInfo.loopCode, Line{runInfo.lineNumber, runInfo.line})
	}
}

func parseFor(runInfo *RunInfo) {
	codes := runInfo.loopCode
	eachVar := runInfo.forEachVar
	setVar := runInfo.forSetStack.Pop()
	for _, v := range setVar {
		CustomVariables[eachVar] = v
		newRunInfo := new(RunInfo)
		newRunInfo.forSetStack = runInfo.forSetStack
		for _, l := range codes {
			newRunInfo.line = l.line
			newRunInfo.lineNumber = l.lineNumber
			CurrLineNum = l.lineNumber
			parse(newRunInfo)
		}
	}
}

func eval(splitLine []string) {
	keyword := splitLine[0]
	if fn, ok := FnKeywords[keyword]; ok {
		fn(splitLine)
	} else {
		panic(fmt.Sprintf("Incorrect keyword '%v' on line %v", keyword, CurrLineNum))
	}
}

func loop(runInfo *RunInfo) {
	runInfo.loopDepth += 1
	if runInfo.loopDepth == 1 {
		if len(runInfo.splitLine) != 2 {
			panic(fmt.Sprintf("Loop Error, loop requires exactly one argument on line %v", CurrLineNum))
		}
		if loopLen, err := strconv.Atoi(runInfo.splitLine[1]); err == nil {
			runInfo.rootLoopLength = loopLen
		} else {
			panic(fmt.Sprintf("Loop Error, loopLength not be number type on line %v", CurrLineNum))
		}
	} else {
		runInfo.loopCode = append(runInfo.loopCode, Line{runInfo.lineNumber, runInfo.line})
	}
}

func endLoop(runInfo *RunInfo) {
	if runInfo.loopDepth == 0 {
		panic(fmt.Sprintf("Loop error on line %v. endloop has no corresponding loop to begin from.", CurrLineNum))
	}
	runInfo.loopDepth -= 1
	if runInfo.loopDepth == 0 {
		parseLoop(runInfo.loopCode, runInfo.rootLoopLength)
		runInfo.loopCode = make([]Line, 50)
	} else {
		runInfo.loopCode = append(runInfo.loopCode, Line{runInfo.lineNumber, runInfo.line})
	}
}

func removeComments(line string) string {
	lines := strings.SplitN(line, "#", 2)
	return lines[0]
}
