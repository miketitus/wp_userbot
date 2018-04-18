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

func TestUserExists(t *testing.T) {
	exists, err := userExists("mike@mike-titus.com")
	if err != nil {
		t.Error(err)
	} else if !exists {
		t.Error("user should have been found")
	}
	exists, err = userExists("nobody@nowhere.xyz")
	if err != nil {
		t.Error(err)
	} else if exists {
		t.Error("user should not have been found")
	}
}
