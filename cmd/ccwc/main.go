package main

import (
	"fmt"
	"log"
	"os"
)

func CountBytes(file *os.File) (int64, error) {
	info, err := file.Stat()
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}

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

	bytes, err := CountBytes(file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(bytes, filename)
}
