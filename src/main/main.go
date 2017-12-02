package main

import (
	//"fmt"
	//"os"
	//"runtime"

	"config"
)

func main() {
	config := config.NewRunConfig()
	config.Parse()

	//fmt.Printf("config = %+v\n", config)

	switch config.Command {
	case "coverage":
		coverage(config)
	case "install":
		install(config)
	case "benchmark":
		benchmark(config)
	case "pprof":
		pprof(config)
	}
}
