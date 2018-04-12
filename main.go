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
		log.Fatal("http.ListenAndServe", err)
	}
}

func parseEmail(w http.ResponseWriter, req *http.Request) {
	/* acknowledge POST from Mailgun  */
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	/* decode body */
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal("ioutil.ReadAll - ", err)
	}
	rawBody := string(bodyBytes)
	rawBody, err = url.QueryUnescape(rawBody)
	if err != nil {
		log.Fatal("url.QueryUnescape - ", err)
	}

	/* trim bulk of unused payload which interferes with regex processing */
	var body string
	i := strings.Index(rawBody, "Mailgun")
	fmt.Printf("i = %d", i)
	if i < 0 {
		body = rawBody
	} else {
		body = rawBody[0:i]
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
