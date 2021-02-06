package main

import (
	"bytes"
	"testing"
)

func TestGreet(t *testing.T) {
	buffer := bytes.Buffer{}

	Greet(&buffer, "Leonardo")

	got := buffer.String()
	want := "Hello, Leonardo"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}