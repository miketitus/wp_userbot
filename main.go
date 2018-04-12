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

	log.Printf("Got: %s\n", req.Header)

	/* decode body */
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll - %s\n", err)
		return
	}
	rawBody := string(bodyBytes)
	rawBody, err = url.QueryUnescape(rawBody)
	if err != nil {
		log.Printf("url.QueryUnescape - %s\n", err)
		return
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

func illegalSenderAlert(e error) {
	log.Println(e)
	if mg == nil {
		mg = mailgun.NewMailgun(mgDomain, mgAPIKey, mgPublicAPIKey)
	}
	message := mg.NewMessage(
		"admin@ncwawood.org",
		"userbot: Illegal Sender",
		e.Error(),
		"mike@mike-titus.com")
	_, _, err := mg.Send(message)
	if err != nil {
		log.Println(err)
	}
}
