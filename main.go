package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	config, err := loadConfig(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(config)
}
