package main

import (
	"fmt"
	"gopkg.in/mailgun/mailgun-go.v1"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type result struct {
	status string
	body   string
}

type user struct {
	first string
	last  string
	email string
}

var mgAPIKey, mgDomain, mgPublicAPIKey, mgValidSender string
var mg mailgun.Mailgun

func main() {
	/* read env settings */
	mgAPIKey = os.Getenv("MG_API_KEY")
	mgDomain = os.Getenv("MG_DOMAIN")
	mgPublicAPIKey = os.Getenv("MG_PUBLIC_API_KEY")
	mgValidSender = os.Getenv("MG_VALID_SENDER")
	/* listen for email POSTs from Mailgun */
	http.HandleFunc("/userbot", parseEmail)
	err := http.ListenAndServe(":8443", nil)
	if err != nil {
		log.Fatal("http.ListenAndServe", err)
	}
}

func parseEmail(w http.ResponseWriter, req *http.Request) {
	/* acknowledge POST from Mailgun  */
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(250)                  // SMTP OK
	log.Printf("Got: %s\n", req.Header) // TODO

	/* decode body */
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll - %s\n", err)
		emailAlert("Parse Error", err.Error())
		return
	}
	rawBody := string(bodyBytes)
	rawBody, err = url.QueryUnescape(rawBody)
	if err != nil {
		log.Printf("url.QueryUnescape - %s\n", err)
		emailAlert("Parse Error", err.Error())
		return
	}

	/* trim bulk of unused payload which interferes with regex processing */
	var body string
	i := strings.Index(rawBody, "Content-Type")
	if i < 0 {
		body = rawBody
		emailAlert("Parse Error", body)
	} else {
		body = rawBody[0:i]
	}

	var sender string
	sender, err = getSender(body)
	if err != nil {
		emailAlert("Illegal Sender", err.Error())
		return
	}
	log.Printf("sender = %s\n", sender)

	recipients := getRecipients(body)
	log.Printf("recipients = %s\n", recipients)
}

func getSender(body string) (string, error) {
	senderRE := regexp.MustCompile("from=([^&]*)")
	raw := senderRE.FindString(body)
	sender := raw[5:]
	if sender != mgValidSender {
		return sender, fmt.Errorf("Illegal sender: '%s'", sender)
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

func emailAlert(subject string, body string) {
	log.Println(subject)
	if mg == nil {
		mg = mailgun.NewMailgun(mgDomain, mgAPIKey, mgPublicAPIKey)
	}
	message := mg.NewMessage(
		"admin@ncwawood.org",
		"userbot: "+subject,
		body,
		"mike@mike-titus.com")
	_, _, err := mg.Send(message)
	if err != nil {
		log.Println(err)
	}
}
