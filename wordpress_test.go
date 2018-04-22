package main

import (
	"fmt"
	"testing"
)

func TestWpAPI(t *testing.T) {
	// test with tags route because it returns minimal data
	_, err := wpAPI("GET", "tags", "")
	if err != nil {
		t.Error(err)
	}
}

func TestUserExists(t *testing.T) {
	if mgAdmins == nil {
		initMain()
	}
	for _, user := range mgAdmins {
		fields := getFields(user)
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

func TestCreateUser(t *testing.T) {
	success, err := createUser("Aaron", "Aardvark", "acct@mike-titus.com") // TODO
	if err != nil {
		t.Error(err)
	} else if !success {
		t.Error("user creation not successful")
	}
	// TODO delete test user
}

func TestGeneratePassword(t *testing.T) {
	testLengths := []uint8{1, 8, 16}
	for _, l := range testLengths {
		p := generatePassword(l)
		if len(p) != int(l) {
			msg := fmt.Sprintf("'%s' is not length %d\n", p, l)
			t.Error(msg)
		}
	}
	p := generatePassword(0)
	if p != "" {
		msg := fmt.Sprintf("length 0 should not have returned '%s'\n", p)
		t.Error(msg)
	}
}
