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
