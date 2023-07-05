package main

import "testing"

func TestAddInts(t *testing.T) {
	sum := doAdd(int64(1), int64(2))
	expected := int64(3)

	if sum != expected {
		t.Errorf("expected '%d' but got '%d'", expected, sum)
	}
}
