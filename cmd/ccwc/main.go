package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func CountBytes(file *os.File) (int64, error) {
	info, err := file.Stat()
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}

func CountLineBreaks(file *os.File) int64 {
	file.Seek(0, io.SeekStart)
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
	file.Seek(0, io.SeekStart)
	var wordCount int64

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanWords)

	for fileScanner.Scan() {
		wordCount++
	}

	return wordCount
}

type fileStats struct {
	bytes          int64
	lineBreakCount int64
	wordCount      int64
}

func getSupportedOptions() []string {
	return []string{"-l", "-w", "-c"}
}

func main() {
	args := os.Args[1:]
	var options []string
	var filename string

	for _, arg := range args {
		if arg[0] == '-' {
			options = append(options, arg)
		} else {
			filename = arg
		}
	}

	if len(options) == 0 {
		options = getSupportedOptions()
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fStats := fileStats{}
	var cols []string
	for _, supportedOption := range getSupportedOptions() {
		for _, option := range options {
			if supportedOption != option {
				continue
			}

			switch option {
			case "-c":
				bytes, err := CountBytes(file)
				if err != nil {
					log.Fatal(err)
				}
				fStats.bytes = bytes
				cols = append(cols, strconv.FormatInt(bytes, 10))

			case "-l":
				lineBreakCount := CountLineBreaks(file)
				fStats.lineBreakCount = lineBreakCount
				cols = append(cols, strconv.FormatInt(lineBreakCount, 10))

			case "-w":
				wordCount := CountWords(file)
				fStats.wordCount = wordCount
				cols = append(cols, strconv.FormatInt(wordCount, 10))

			default:
				log.Fatal("invalid command argument")
			}
		}
	}

	cols = append(cols, filename)
	fmt.Println(strings.Join(cols, " "))
}
