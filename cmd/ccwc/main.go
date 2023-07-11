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

type OSProcess interface {
	args() []string
	writeOutput(s string) error
}

type Process struct {
	output io.Writer
}

func (p *Process) args() []string {
	return os.Args
}

func (p *Process) writeOutput(s string) error {
	_, err := fmt.Fprint(p.output, s+"\n")
	return err
}

type StubbedProcess struct {
	_args  []string
	output io.Writer
}

func (p *StubbedProcess) args() []string {
	a := []string{"nulled_process_go"}
	a = append(a, p._args...)
	return a
}

func (p *StubbedProcess) writeOutput(s string) error {
	_, err := fmt.Fprint(p.output, s)
	return err
}

type CommandLine struct {
	process OSProcess
}

func (c *CommandLine) args() []string {
	return c.process.args()
}

func NewCommandLine() *CommandLine {
	return &CommandLine{&Process{output: os.Stdout}}
}

func NewNullCommandLine(args []string, output io.Writer) *CommandLine {
	return &CommandLine{&StubbedProcess{_args: args, output: output}}
}

type App struct {
	commandLine *CommandLine
}

func (a App) run() {
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
				log.Fatal("invalid command argument")
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
	res := strings.Join(cols, " ")
	a.commandLine.process.writeOutput(res)
}

func main() {
	commandLine := NewCommandLine()
	app := App{commandLine}
	app.run()
}
