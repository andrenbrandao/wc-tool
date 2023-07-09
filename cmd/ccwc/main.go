package main

import (
	"bufio"
	"errors"
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
	file.Seek(0, io.SeekStart)

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

type OSProcess interface {
	args() []string
}

type Process struct{}

func (p *Process) args() []string {
	return os.Args
}

type StubbedProcess struct {
	_args []string
}

func (p *StubbedProcess) args() []string {
	a := []string{"nulled_process_go"}
	a = append(a, p._args...)
	return a
}

type CommandLine struct {
	process OSProcess
}

func (c *CommandLine) args() []string {
	return c.process.args()
}

func NewCommandLine() *CommandLine {
	return &CommandLine{&Process{}}
}

func NewNullCommandLine(args []string) *CommandLine {
	return &CommandLine{&StubbedProcess{args}}
}

type App struct {
	commandLine *CommandLine
}

func (a App) run() (string, error) {
	args := a.commandLine.args()[1:]
	var filename string

	printLineBreaks := false
	printWords := false
	printChars := false
	printBytes := false
	for _, arg := range args {
		if arg[0] == '-' {

			switch arg {
			case "-l":
				printLineBreaks = true

			case "-w":
				printWords = true

			case "-m":
				printChars = true

			case "-c":
				printBytes = true

			default:
				return "", errors.New("invalid command argument")
			}
		} else {
			filename = arg
		}
	}

	if !printLineBreaks && !printWords && !printChars && !printBytes {
		printLineBreaks = true
		printWords = true
		printBytes = true
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	var file *os.File
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		file = os.Stdin
	} else {
		file, err = os.Open(filename)
		if err != nil {
			return "", err
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
	return strings.Join(cols, " "), nil
}

func main() {
	commandLine := NewCommandLine()
	app := App{commandLine}

	res, err := app.run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}
