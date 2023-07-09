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

func main() {
	args := os.Args[1:]
	var filename string

	printLineBreaks := false
	printWords := false
	printBytes := false
	for _, arg := range args {
		if arg[0] == '-' {

			switch arg {
			case "-l":
				printLineBreaks = true

			case "-w":
				printWords = true

			case "-c":
				printBytes = true

			default:
				log.Fatal("invalid command argument")
			}
		} else {
			filename = arg
		}
	}

	if !printLineBreaks && !printWords && !printBytes {
		printLineBreaks = true
		printWords = true
		printBytes = true
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

	fStats := GetFileStats(file)
	var cols []string

	if printLineBreaks {
		cols = append(cols, strconv.FormatInt(fStats.lineBreakCount, 10))
	}

	if printWords {
		cols = append(cols, strconv.FormatInt(fStats.wordCount, 10))
	}

	if printBytes {
		cols = append(cols, strconv.FormatInt(fStats.bytes, 10))
	}

	cols = append(cols, filename)
	fmt.Println(strings.Join(cols, " "))
}
