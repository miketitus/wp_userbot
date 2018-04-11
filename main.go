package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type user struct {
	first string
	last  string
	email string
}

var validSender = ""

func main() {
	/* read env settings */
	validSender = os.Getenv("USERBOT_SENDER")
	fmt.Printf("validSender: '%s'\n", validSender)
	/* lauch http server */
	http.HandleFunc("/userbot", parseEmail)
	err := http.ListenAndServe(":8443", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

func parseEmail(w http.ResponseWriter, req *http.Request) {
	/* acknowledge POST from Mailgun  */
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	/* decode body */
	bodyBytes, err1 := ioutil.ReadAll(req.Body)
	if err1 != nil {
		log.Fatal("ioutil.ReadAll", err1)
	}
	bodyString := string(bodyBytes)
	body, err2 := url.QueryUnescape(bodyString)
	if err2 != nil {
		log.Fatal("QueryUnescape", err2)
	}

	sender, err3 := getSender(body)
	if err3 != nil {
		illegalSenderAlert(err3)
	}
	fmt.Printf("sender = %s\n", sender)

	recipients := getRecipients(body)
	fmt.Printf("recipients = %s\n", recipients)
}

func getSender(body string) (string, error) {
	senderRE := regexp.MustCompile("from=(.*)&")
	sender := senderRE.FindString(body)[5:]
	if sender != validSender {
		return "", fmt.Errorf("Illegal sender: '%s'", sender)
	}
	return sender, nil
}

func getRecipients(body string) []string {
	recipientRE := regexp.MustCompile("To=(.*)&")
	raw := recipientRE.FindString(body)
	recipients := strings.Split(raw[3:], ", ")
	return recipients
}

func createUsers(recipients []string) {

}

func illegalSenderAlert(err error) {
	log.Fatal(err)
}
