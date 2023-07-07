package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	command := os.Args[1]
	filename := os.Args[2]

	if command != "-c" {
		log.Fatal("invalid command argument")
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	info, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(info.Size(), filename)
}
