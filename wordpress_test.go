package main

import (
	"testing"
)

func TestCallWP(t *testing.T) {
	// test with tags route because it returns minimal data
	_, err := callWP("tags", "")
	if err != nil {
		t.Error(err)
	}
}
