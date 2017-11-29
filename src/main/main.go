package main

import (
	"fmt"

	"config"
)

func main() {
	config := config.NewRunConfig()
	config.Parse()

	fmt.Printf("config = %+v\n", config)
}
