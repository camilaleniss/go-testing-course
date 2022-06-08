package main

import "testing"

func TestAddSuccess(t *testing.T) {
	result := Add(2, 20)

	expected := 22

	if result != expected {
		t.Errorf("got %q, expected %q", result, expected)
	}
}
