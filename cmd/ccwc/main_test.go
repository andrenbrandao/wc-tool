package main

import (
	"bytes"
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
		buffer := &bytes.Buffer{}
		commandLine := NewNullCommandLine(args, buffer)
		app := App{commandLine}
		app.run()
		return buffer.String(), nil
	}

	t.Run("without arguments", func(t *testing.T) {
		got, _ := run([]string{"testdata/data.txt"})
		want := "8 351 2432 testdata/data.txt"

		if got != want {
			t.Errorf("got %+v want %+v", got, want)
		}
	})

	t.Run("with -l", func(t *testing.T) {
		got, _ := run([]string{"-l", "testdata/data.txt"})
		want := "8 testdata/data.txt"

		if got != want {
			t.Errorf("got %+v want %+v", got, want)
		}
	})

	t.Run("with -w", func(t *testing.T) {
		got, _ := run([]string{"-w", "testdata/data.txt"})
		want := "351 testdata/data.txt"

		if got != want {
			t.Errorf("got %+v want %+v", got, want)
		}
	})

	t.Run("with -c", func(t *testing.T) {
		got, _ := run([]string{"-c", "testdata/data.txt"})
		want := "2432 testdata/data.txt"

		if got != want {
			t.Errorf("got %+v want %+v", got, want)
		}
	})

	t.Run("with -m", func(t *testing.T) {
		got, _ := run([]string{"-m", "testdata/data.txt"})
		want := "2430 testdata/data.txt"

		if got != want {
			t.Errorf("got %+v want %+v", got, want)
		}
	})
}
