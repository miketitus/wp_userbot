package main

import (
	"fmt"
	"io/ioutil"
	"testing"
)

var from = `Content-Disposition: form-data; name="From"

%s`

var to = `Content-Disposition: form-data; name="To"

%s`

func TestWriteBody(t *testing.T) {
	testBytes := []byte("Testing one, two, three.\n")
	err := writeBody(testBytes)
	if err != nil {
		t.Error(err)
	}
}

func TestGetSender(t *testing.T) {
	s := "No Body <nobody@nowhere.xyz>"
	e, err := getEmail(s)
	if err != nil {
		t.Error(err)
	}
	test := fmt.Sprintf(from, s)
	sender, err := getSender(test)
	if err != nil {
		t.Error(err)
	} else if e.Address != sender.Address {
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
	s, err := getEmail("No Body <nobody@nowhere.xyz>")
	if err != nil {
		t.Error(err)
	} else if senderIsAdmin(s) {
		t.Errorf("'%s' was declared valid", s)
	}
}

func TestGetRecipients(t *testing.T) {
	s := "John1 Doe <john1@john.doe>, John2 Doe <john2@john.doe>"
	test := fmt.Sprintf(to, s)
	recipients, err := getRecipients(test)
	if err != nil {
		t.Error(err)
	}
	e1, err := getEmail(recipients[0])
	if err != nil {
		t.Error(err)
	} else if e1.Address != "john1@john.doe" {
		t.Errorf("e1.Address should not be %s", e1.Address)
	}
	e2, err := getEmail(recipients[1])
	if err != nil {
		t.Error(err)
	} else if e2.Address != "john2@john.doe" {
		t.Errorf("e2.Address should not be %s", e2.Address)
	}
}

func TestGetEmail(t *testing.T) {
	// test valid format
	email, err := getEmail("John Doe <john@john.doe>")
	if err != nil {
		t.Error(err)
	} else if email.First != "John" {
		t.Errorf("First should not be %s", email.First)
	} else if email.Last != "Doe" {
		t.Errorf("Last should not be %s", email.Last)
	} else if email.Address != "john@john.doe" {
		t.Errorf("Address should not be %s", email.Address)
	}
	// test invalid format
	_, err = getEmail("John <john@john.doe>")
	if err == nil {
		t.Error("getEmail() failed to return an error")
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
	if mgAdmins == nil {
		initMain()
	}
	email := Email{"John", "Doe", mgAdmins[0].Address}
	emailUser(email, "jdoe", "secret")
}

func TestParsing(t *testing.T) {
	testFile := "./assets/152753133.txt"
	body, err := readTestFile(testFile)
	if err != nil {
		// don't fail test, just return
		t.Logf("can't open %s, skipping TestParsing()\n", testFile)
		return
	}
	sender, err := getSender(body)
	if err != nil {
		t.Error(err)
	}
	if !senderIsAdmin(sender) {
		t.Errorf("'%s' was declared to not be an admin", sender)
	}
	recipients, err := getRecipients(body)
	if err != nil {
		t.Error(err)
	}
	for _, r := range recipients {
		if !isValidEmail(r) {
			t.Errorf("'%s' was declared to not be a valid email address", r)
		}
	}
}

func readTestFile(testFile string) (string, error) {
	bytes, err := ioutil.ReadFile(testFile)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}
