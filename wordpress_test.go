package main

import (
	"strings"
	"testing"
)

func TestWPGet(t *testing.T) {
	// test with tags route because it returns minimal data
	_, err := wpGet("tags", "")
	if err != nil {
		t.Error(err)
	}
}

func TestUserExists(t *testing.T) {
	if mgAdmins == nil {
		initMain()
	}
	for _, user := range mgAdmins {
		fields := strings.Fields(user)
		exists, err := userExists(fields[2])
		if err != nil {
			t.Error(err)
		} else if !exists {
			t.Error("user should have been found")
		}
	}
	exists, err := userExists("nobody@nowhere.xyz")
	if err != nil {
		t.Error(err)
	} else if exists {
		t.Error("user should not have been found")
	}
}
