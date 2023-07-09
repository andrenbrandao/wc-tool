package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type fileStats struct {
	bytes          int64
	lineBreakCount int64
	wordCount      int64
}

func GetFileStats(file *os.File) fileStats {
	file.Seek(0, io.SeekStart)

	var wordCount int64
	var bytes int64
	var lineBreakCount int64

	reader := bufio.NewReader(file)

	inWord := false
	for {
		c, sz, err := reader.ReadRune()
		bytes += int64(sz)

		if err != nil {
			if err == io.EOF {
				if inWord {
					wordCount++
				}
				break
			} else {
				log.Fatal(err)
			}
		}

		if unicode.IsSpace(c) {
			if inWord {
				wordCount++
			}
			if c == '\n' {
				lineBreakCount++
			}
			inWord = false
		} else {
			inWord = true
		}
	}

	return fileStats{bytes: bytes, lineBreakCount: lineBreakCount, wordCount: wordCount}
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

	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}

	var file *os.File
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		file = os.Stdin
	} else {
		file, err = os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}

	if len(options) == 0 {
		options = getSupportedOptions()
	}

	fStats := GetFileStats(file)
	var cols []string
	for _, supportedOption := range getSupportedOptions() {
		for _, option := range options {
			if supportedOption != option {
				continue
			}

			switch option {
			case "-c":
				cols = append(cols, strconv.FormatInt(fStats.bytes, 10))

			case "-l":
				cols = append(cols, strconv.FormatInt(fStats.lineBreakCount, 10))

			case "-w":
				cols = append(cols, strconv.FormatInt(fStats.wordCount, 10))

			default:
				log.Fatal("invalid command argument")
			}
		}
	}

	cols = append(cols, filename)
	fmt.Println(strings.Join(cols, " "))
}
