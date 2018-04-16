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
