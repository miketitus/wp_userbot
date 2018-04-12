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

var validSender string

func main() {
	/* read env settings */
	validSender = os.Getenv("USERBOT_SENDER")
	/* launch http server */
	http.HandleFunc("/userbot", parseEmail)
	err := http.ListenAndServe(":8443", nil)
	if err != nil {
		log.Fatal("http.ListenAndServe", err)
	}
}

func parseEmail(w http.ResponseWriter, req *http.Request) {
	/* acknowledge POST from Mailgun  */
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(250) // SMTP OK

	/* decode body */
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll - %s\n", err)
	}
	rawBody := string(bodyBytes)
	rawBody, err = url.QueryUnescape(rawBody)
	if err != nil {
		log.Printf("url.QueryUnescape - %s\n", err)
	}

	/* trim bulk of unused payload which interferes with regex processing */
	var body string
	i := strings.Index(rawBody, "Content-Type")
	if i < 0 {
		// TODO
		body = rawBody
	} else {
		body = rawBody[0:i]
	}

	var sender string
	sender, err = getSender(body)
	if err != nil {
		illegalSenderAlert(err)
	}
	log.Printf("sender = %s\n", sender)

	recipients := getRecipients(body)
	log.Printf("recipients = %s\n", recipients)
}

func getSender(body string) (string, error) {
	senderRE := regexp.MustCompile("from=([^&]*)")
	raw := senderRE.FindString(body)
	sender := raw[5:]
	if sender != validSender {
		return "", fmt.Errorf("Illegal sender: '%s'", sender)
	}
	return sender, nil
}

func getRecipients(body string) []string {
	recipientRE := regexp.MustCompile("To=([^&]*)")
	raw := recipientRE.FindString(body)
	recipients := strings.Split(raw[3:], ", ")
	return recipients
}

func createUsers(recipients []string) {

}

func illegalSenderAlert(err error) {
	log.Println(err)
}
