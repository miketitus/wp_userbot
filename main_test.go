package main

import (
	"fmt"
	"testing"
)

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
	s := "John Doe <john@john.doe>"
	f := getFields(s)
	if f[0] != "John" {
		t.Errorf("f[0] should not be %s", f[0])
	}
	if f[1] != "Doe" {
		t.Errorf("f[1] should not be %s", f[1])
	}
	if f[2] != "john@john.doe" {
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
