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

type ICommandLine interface {
	args() *Args
	writeOutput(s string) error
}
type CommandLine struct {
	output io.Writer
	flag   *flag.FlagSet
}

type StubbedCommandLine struct {
	_args  []string
	output io.Writer
	flag   *flag.FlagSet
}

func (c *CommandLine) args() *Args {
	var printLineBreaks, printWords, printChars, printBytes bool
	flag.BoolVar(&printLineBreaks, "l", false, "print line breaks")
	flag.BoolVar(&printWords, "w", false, "print words")
	flag.BoolVar(&printChars, "m", false, "print chars")
	flag.BoolVar(&printBytes, "c", false, "print bytes")

	flag.Parse()

	filename := flag.CommandLine.Arg(0)

	return &Args{printLineBreaks, printWords, printChars, printBytes, filename}
}

func (c *CommandLine) writeOutput(s string) error {
	_, err := fmt.Fprint(c.output, s+"\n")
	return err
}

func NewCommandLine() *CommandLine {
	return &CommandLine{output: os.Stdout, flag: flag.NewFlagSet(os.Args[0], flag.ExitOnError)}
}

func (c *StubbedCommandLine) writeOutput(s string) error {
	_, err := fmt.Fprint(c.output, s)
	return err
}

func (c *StubbedCommandLine) args() *Args {
	var printLineBreaks, printWords, printChars, printBytes bool

	c.flag.BoolVar(&printLineBreaks, "l", false, "print line breaks")
	c.flag.BoolVar(&printWords, "w", false, "print words")
	c.flag.BoolVar(&printChars, "m", false, "print chars")
	c.flag.BoolVar(&printBytes, "c", false, "print bytes")

	c.flag.Parse(c._args)

	filename := c.flag.Arg(0)

	return &Args{printLineBreaks, printWords, printChars, printBytes, filename}
}

func NewNullCommandLine(args []string, output io.Writer) *StubbedCommandLine {
	return &StubbedCommandLine{_args: args, output: output, flag: flag.NewFlagSet(os.Args[0], flag.ExitOnError)}
}

type App struct {
	commandLine ICommandLine
}

type Args struct {
	printLineBreaks,
	printWords,
	printChars,
	printBytes bool
	filename string
}

func (a *App) run() {
	var printLineBreaks, printWords, printChars, printBytes bool
	var filename string

	args := a.commandLine.args()
	printLineBreaks = args.printLineBreaks
	printWords = args.printWords
	printChars = args.printChars
	printBytes = args.printBytes
	filename = args.filename

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
	a.commandLine.writeOutput(res)
}

func main() {
	commandLine := NewCommandLine()
	app := App{commandLine}
	app.run()
}
