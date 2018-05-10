package main

import (
	"fmt"
	"testing"
)

func TestWriteBody(t *testing.T) {
	err := writeBody("Testing one, two, three.\n")
	if err != nil {
		t.Error(err)
	}
}

func TestSenderIsAdmin(t *testing.T) {
	if mgAdmins == nil {
		initMain()
	}
	for _, s := range mgAdmins {
		test := fmt.Sprintf("&from=%s&", s)
		if !senderIsAdmin(test) {
			t.Errorf("'%s' was declared invalid", s)
		}
	}
	s := "&from=nobody@nowhere.xyz&"
	if senderIsAdmin(s) {
		t.Errorf("'%s' was declared valid", s)
	}
}

func TestGetFields(t *testing.T) {
	var f []string
	// one field
	f = getFields("<john@john.doe>")
	if f[0] != "john@john.doe" {
		t.Errorf("f[0] should not be %s", f[0])
	}
	// two fields
	f = getFields("John <john@john.doe>")
	if f[0] != "John" {
		t.Errorf("f[0] should not be %s", f[0])
	} else if f[1] != "john@john.doe" {
		t.Errorf("f[1] should not be %s", f[2])
	}
	// three fields
	f = getFields("John Doe <john@john.doe>")
	if f[0] != "John" {
		t.Errorf("f[0] should not be %s", f[0])
	} else if f[1] != "Doe" {
		t.Errorf("f[1] should not be %s", f[1])
	} else if f[2] != "john@john.doe" {
		t.Errorf("f[2] should not be %s", f[2])
	}
}

func TestIsUserBot(t *testing.T) {
	s := make([]string, 1, 1)
	s[0] = mgUserBot
	if !isUserBot(s) {
		t.Errorf("'%s' was declared to not be a userbot", s)
	}
	s[0] = "foo@abc.xyz"
	if isUserBot(s) {
		t.Errorf("'%s' was declared to be a userbot", s)
	}
}

func TestIsValidEmail(t *testing.T) {
	validEmails := []string{"abc@xyz.com", "abc+@xyz.us", "abc@wa.us"}
	inValidEmails := []string{"abcxyz.com", "abc@@xyz.us", "abc@waus", "abc@waus.", "abcxyz"}
	for _, e := range validEmails {
		if !isValidEmail(e) {
			t.Errorf("'%s' was declared invalid", e)
		}
	}
	for _, e := range inValidEmails {
		if isValidEmail(e) {
			t.Errorf("'%s' was declared valid", e)
		}
	}
}

func TestEmailUser(t *testing.T) {
	emailUser("jdoe", "John", "Doe", "acct@mike-titus.com", "secret")
}
