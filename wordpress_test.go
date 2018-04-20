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

func TestCreateUser(t *testing.T) {
	success, err := createUser("Aaron", "Aardvark", "acct@mike-titus.com")
	if err != nil {
		t.Error(err)
	} else if !success {
		t.Error("user creation not successful")
	}
}

func TestGeneratePassword(t *testing.T) {
	testLengths := []int{1, 8, 16}
	for _, l := range testLengths {
		p := generatePassword(l)
		if len(p) != l {
			msg := fmt.Sprintf("'%s' is not length %d\n", p, l)
			t.Error(msg)
		}
	}
	badLengths := []int{-1, 0}
	for _, l := range badLengths {
		p := generatePassword(l)
		if p != "" {
			msg := fmt.Sprintf("length %d should not have returned '%s'\n", l, p)
			t.Error(msg)
		}
	}
}
