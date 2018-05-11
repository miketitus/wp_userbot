package main

import (
	"fmt"
	"testing"
)

func TestWriteBody(t *testing.T) {
	testBytes := []byte("Testing one, two, three.\n")
	err := writeBody(testBytes)
	if err != nil {
		t.Error(err)
	}
}

func TestGetSender(t *testing.T) {
	goodFrom := "[\"From\", \"%s\"], [\"More\"], \"Blah\""
	s := "No Body <nobody@nowhere.xyz>"
	test := fmt.Sprintf(goodFrom, s)
	sender, err := getSender(test)
	if err != nil {
		t.Error(err)
	} else if s != sender {
		t.Errorf("'%s' does not match '%s'", s, sender)
	}
}

func TestSenderIsAdmin(t *testing.T) {
	if mgAdmins == nil {
		initMain()
	}
	for _, s := range mgAdmins {
		if !senderIsAdmin(s) {
			t.Errorf("'%s' was declared invalid", s)
		}
	}
	s := "No Body <nobody@nowhere.xyz>"
	if senderIsAdmin(s) {
		t.Errorf("'%s' was declared valid", s)
	}
}

func TestGetRecipients(t *testing.T) {
	test := "[\"To\", \"John1 Doe <john1@john.doe>, John2 Doe <john2@john.doe>\"], [\"Foo\", \"Bar\""
	recipients := getRecipients(test)
	f1 := getFields(recipients[0])
	if f1[2] != "john1@john.doe" {
		t.Errorf("f1[2] should not be %s", f1[2])
	}
	f2 := getFields(recipients[1])
	if f2[2] != "john2@john.doe" {
		t.Errorf("f2[2] should not be %s", f2[2])
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
	emailUser("jdoe", "John", "Doe", "acct@mike-titus.com", "secret") // TODO
}
