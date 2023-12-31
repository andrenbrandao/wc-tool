package main

import (
	"bufio"
	"flag"
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
	charsCount     int64
}

func GetFileStats(file *os.File) fileStats {
	var wordCount int64
	var bytes int64
	var lineBreakCount int64
	var charsCount int64

	reader := bufio.NewReader(file)

	inWord := false
	for {
		c, sz, err := reader.ReadRune()

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

		bytes += int64(sz)
		charsCount++

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

	return fileStats{bytes: bytes, lineBreakCount: lineBreakCount, wordCount: wordCount, charsCount: charsCount}
}

func main() {
	var printLineBreaks, printWords, printChars, printBytes bool

	flag.BoolVar(&printLineBreaks, "l", false, "print line breaks")
	flag.BoolVar(&printWords, "w", false, "print words")
	flag.BoolVar(&printChars, "m", false, "print chars")
	flag.BoolVar(&printBytes, "c", false, "print bytes")

	flag.Parse()

	filename := flag.CommandLine.Arg(0)

	if !printLineBreaks && !printWords && !printChars && !printBytes {
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

	if printChars {
		cols = append(cols, strconv.FormatInt(fStats.charsCount, 10))
	}

	if printBytes {
		cols = append(cols, strconv.FormatInt(fStats.bytes, 10))
	}

	cols = append(cols, filename)
	fmt.Println(strings.Join(cols, " "))
}
