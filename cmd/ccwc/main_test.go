package main

import (
	"log"
	"os"
	"testing"
)

func TestCountBytes(t *testing.T) {
	f, err := os.Open("testdata/data.txt")
	if err != nil {
		log.Fatal(err)
	}
	got, _ := CountBytes(f)
	var want int64 = 2432

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func TestCountLineBreaks(t *testing.T) {
	f, err := os.Open("testdata/data.txt")
	if err != nil {
		log.Fatal(err)
	}
	got := CountLineBreaks(f)
	var want int64 = 8

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
