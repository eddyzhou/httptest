package main

import (
	"io"
	"flag"
	"log"
	"os"
)

func init() {
	log.SetOutput(io.Writer(os.Stderr))
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
}

func main() {
	var host string
	flag.StringVar(&host, "host", "http://m2.qiushibaike.com", "set test host")
	flag.Parse()

	BaseVariables["HOST"] = host
	args := flag.Args()
	if len(args) < 1 {
		EnterREPL()
		return
	}

	scriptFile := args[0]
	ParseFile(scriptFile)
}
