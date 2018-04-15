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
		emailResults("Parse Error", err.Error())
		return
	}
	rawBody := string(bodyBytes)
	rawBody, err = url.QueryUnescape(rawBody)
	if err != nil {
		log.Printf("url.QueryUnescape - %s\n", err)
		emailResults("Parse Error", err.Error())
		return
	}

	/* trim bulk of unused payload which interferes with regex processing */
	var body string
	i := strings.Index(rawBody, "Content-Type")
	if i < 0 {
		body = rawBody
		emailResults("Parse Error", body)
	} else {
		body = rawBody[0:i]
	}

	// validate sender
	var sender string
	sender, err = getSender(body)
	if err != nil {
		emailResults("Illegal Sender", err.Error())
		return
	}
	log.Printf("sender = %s\n", sender)

	// turn recipients into users
	parseRecipients(body)
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

func parseRecipients(body string) {
	var resultBody []string
	resultSubject := "Success"
	recipientRE := regexp.MustCompile("To=([^&]*)")
	raw := recipientRE.FindString(body)
	recipients := strings.Split(raw[3:], ", ")
	log.Printf("recipients = %s\n", recipients)
	for _, r := range recipients {
		fields := strings.Fields(r)
		if len(fields) == 3 {
			// valid structure
			if fields[2] == "<mike@mike-titus.com>" {
				continue // skip
			} else if isValidEmail(fields[2]) {
				resultBody = append(resultBody, fmt.Sprintf("Invalid email: %s", r))
			} else {
				userResult := createUser(fields[0], fields[1], fields[2])
				resultBody = append(resultBody, fmt.Sprintf("%s: %s", userResult, r))
			}
		} else if fields[0] == "userbot" || fields[0] == "<userbot@ncwawood.org>" {
			// that's me! -- ignore
		} else {
			// error TODO
			resultBody = append(resultBody, fmt.Sprintf("Invalid format: %s", r))
			resultSubject = "Error(s) found"
		}
	}
	emailResults(resultSubject, strings.Join(resultBody, "\n"))
}

func createUser(first, last, email string) string {
	return "success"
}

func isValidEmail(email string) bool {
	// mimimal validation regex, could be a lot more complex
	emailRE := regexp.MustCompile(".*@.*\\..*")
	return emailRE.FindStringIndex(email) != nil
}

func emailResults(subject string, body string) {
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
