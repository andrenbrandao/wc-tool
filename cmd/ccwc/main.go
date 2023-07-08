package main

import (
	"bufio"
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

func CountLineBreaks(file *os.File) int64 {
	var lineBreakCount int64

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanRunes)

	for fileScanner.Scan() {
		if fileScanner.Text() == "\n" {
			lineBreakCount++
		}
	}

	return lineBreakCount
}

func CountWords(file *os.File) int64 {
	var wordCount int64

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanWords)

	for fileScanner.Scan() {
		wordCount++
	}

	return wordCount
}

func main() {
	command := os.Args[1]
	filename := os.Args[2]

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	switch command {
	case "-c":
		bytes, err := CountBytes(file)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(bytes, filename)

	case "-l":
		lineBreaks := CountLineBreaks(file)
		fmt.Println(lineBreaks, filename)

	case "-w":
		wordCount := CountWords(file)
		fmt.Println(wordCount, filename)

	default:
		log.Fatal("invalid command argument")
	}

}
