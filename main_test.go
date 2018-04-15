package main

import (
	"testing"
)

func TestIsValidSender(t *testing.T) {
	for _, s := range mgValidSenders {
		if !IsValidSender(s) {
			t.Errorf("'%s' was declared invalid", s)
		}
	}
	s := "Invalid"
	if IsValidSender(s) {
		t.Errorf("'%s' was declared valid", s)
	}
}

func TestIsValidEmail(t *testing.T) {
	validEmails := []string{"abc@xyz.com", "abc+@xyz.us", "abc@wa.us"}
	inValidEmails := []string{"abcxyz.com", "abc@@xyz.us", "abc@waus", "abc@waus.", "abcxyz"}
	for _, e := range validEmails {
		if !IsValidEmail(e) {
			t.Errorf("'%s' was declared invalid", e)
		}
	}
	for _, e := range inValidEmails {
		if IsValidEmail(e) {
			t.Errorf("'%s' was declared valid", e)
		}
	}
}
