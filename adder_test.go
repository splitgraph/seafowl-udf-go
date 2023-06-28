package main

import "testing"

func TestAddInts(t *testing.T) {
	sum := AddInts(1, 2)
	expected := 3

	if sum != expected {
		t.Errorf("expected '%d' but got '%d'", expected, sum)
	}
}