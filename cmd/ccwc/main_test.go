package main

import (
	"log"
	"os"
	"testing"
)

func TestGetFileStats(t *testing.T) {
	f, err := os.Open("testdata/data.txt")
	if err != nil {
		log.Fatal(err)
	}
	got := GetFileStats(f)
	want := fileStats{bytes: 2432, lineBreakCount: 8, wordCount: 351, charsCount: 2430}

	if got != want {
		t.Errorf("got %+v want %+v", got, want)
	}
}

func TestApp(t *testing.T) {
	run := func(args []string) (string, error) {
		commandLine := NewNullCommandLine(args)
		app := App{commandLine}
		return app.run()
	}

	got, _ := run([]string{"testdata/data.txt"})
	want := "8 351 2432 testdata/data.txt"

	if got != want {
		t.Errorf("got %+v want %+v", got, want)
	}
}
