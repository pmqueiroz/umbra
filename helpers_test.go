package main

import (
	"fmt"
	"os"
	"testing"
)

func TestValidFile(t *testing.T) {
	file, err := os.CreateTemp("", "temp-*.txt")

	if err != nil {
		t.Fatalf("Error creating tmp file")
	}

	defer os.Remove(file.Name())

	_, err = file.WriteString("Temporary file contents")

	if err != nil {
		fmt.Println(err)
		return
	}

	dat, err := readFile(file.Name())

	if err != nil {
		t.Fatal(err.Error())
	}

	if dat != "Temporary file contents" {
		t.Error("should return valid file content but didn't")
	}
}

func TestNonexistentFile(t *testing.T) {
	_, err := readFile("nonexistent.file")

	if err == nil {
		t.Error("should return an error but didn't")
	}
}
