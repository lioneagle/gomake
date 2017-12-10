package main

import (
	"fmt"
	//"os"
	//"runtime"

	"config"
)

func main() {
	config := config.NewRunConfig()
	config.Parse()

	//fmt.Printf("config = %+v\n", config)

	var err error = nil

	switch config.Command {
	case "coverage":
		err = coverage(config)
	case "install":
		err = install(config)
	case "benchmark":
		err = benchmark(config)
	case "pprof":
		err = pprof(config)
	}
	if err != nil {
		fmt.Println("err =", err)
	}

}
