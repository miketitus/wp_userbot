package main

import (
	"testing"
)

func TestInitWordPress(t *testing.T) {
	initWordPress()
	_, _, err := wpClient.Users.Me(wpContext, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestCallWP(t *testing.T) {
	// test with tags route because it returns minimal data
	_, err := callWP("tags", "")
	if err != nil {
		t.Error(err)
	}
}
