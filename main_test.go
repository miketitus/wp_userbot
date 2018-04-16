package main

import (
	"testing"
)

func TestSenderIsAdmin(t *testing.T) {
	for _, s := range mgAdmins {
		if !senderIsAdmin(s) {
			t.Errorf("'%s' was declared invalid", s)
		}
	}
	s := "Invalid"
	if senderIsAdmin(s) {
		t.Errorf("'%s' was declared valid", s)
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
