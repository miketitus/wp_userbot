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

	"gopkg.in/mailgun/mailgun-go.v1"
)

type user struct {
	first string
	last  string
	email string
}

var mgAPIKey, mgDomain, mgPublicAPIKey, mgUserBot string
var mgAdmins []string
var mg mailgun.Mailgun

func main() {
	// read env settings
	mgAPIKey = os.Getenv("MG_API_KEY")
	mgDomain = os.Getenv("MG_DOMAIN")
	mgPublicAPIKey = os.Getenv("MG_PUBLIC_API_KEY")
	mgUserBot = os.Getenv("MG_USERBOT")
	mgAdmins = strings.Split(os.Getenv("MG_ADMINS"), ", ")
	// listen for email POSTs from Mailgun
	http.HandleFunc("/userbot", parseEmail)
	err := http.ListenAndServe(":8443", nil)
	if err != nil {
		log.Fatal("http.ListenAndServe", err)
	}
}

// parseEmail is the main event loop, executing for each received email.
func parseEmail(w http.ResponseWriter, req *http.Request) {
	// acknowledge POST from Mailgun
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(250)                  // SMTP OK
	log.Printf("Got: %s\n", req.Header) // TODO

	body, err := getBody(req)
	if err != nil {
		return
	}

	if !senderIsAdmin(body) {
		// ignore: spam, or a recipient hit "reply to all"
		emailResults("Illegal Sender", body) // TODO
		return
	}
	parseRecipients(body)
}

// getBody extracts and unescapes the email body from the Mailgun POST.
func getBody(req *http.Request) (string, error) {
	// read body
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll - %s\n", err)
		emailResults("Parse Error", err.Error())
		return "", err
	}
	// decode body
	rawBody := string(bodyBytes)
	rawBody, err = url.QueryUnescape(rawBody)
	if err != nil {
		log.Printf("url.QueryUnescape - %s\n", err)
		emailResults("Parse Error", err.Error())
		return "", err
	}
	// TODO trim bulk of unused payload which interferes with regex processing
	var body string
	i := strings.Index(rawBody, "Content-Type")
	if i < 0 {
		body = rawBody
		emailResults("Parse Error", body)
	} else {
		body = rawBody[0:i]
	}
	return body, nil
}

// senderIsAdmin verifies that the decoded email message came from an approved email address.
func senderIsAdmin(body string) bool {
	var sender string
	senderRE := regexp.MustCompile("from=([^&]*)")
	raw := senderRE.FindString(body)
	if raw == "" {
		// this should never happen, but let's keep on truckin'
		sender = raw
	} else {
		sender = raw[5:]
	}
	for _, s := range mgAdmins {
		if s == sender {
			return true
		}
	}
	return false
}

// isUserBot is used to detect and ignore the userbot while processing recipient addresses.
func isUserBot(fields []string) bool {
	for _, f := range fields {
		if f == mgUserBot {
			return true
		}
	}
	return false
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
		if isUserBot(fields) {
			continue // ignore
		} else if len(fields) == 3 {
			// valid structure
			if senderIsAdmin(fields[2]) {
				continue // admin email, skip
			} else if isValidEmail(fields[2]) {
				resultBody = append(resultBody, fmt.Sprintf("Invalid email: %s", r))
			} else {
				userResult := createUser(fields[0], fields[1], fields[2])
				resultBody = append(resultBody, fmt.Sprintf("%s: %s", userResult, r))
			}
		} else {
			// error
			resultBody = append(resultBody, fmt.Sprintf("Invalid format: %s", r))
			resultSubject = "Error(s) found"
		}
	}
	emailResults(resultSubject, strings.Join(resultBody, "\n"))
}

// isValidEmail verifies that a recipient email address is in valid format.
// Uses a very simple regex designed to catch basic errors, but not nearly all edge cases.
func isValidEmail(email string) bool {
	// mimimal validation regex, could be a lot more complex
	emailRE := regexp.MustCompile("[^@]+@[^@]+\\..+")
	return emailRE.FindStringIndex(email) != nil
}

func createUser(first, last, email string) string {
	return "success"
}

// emailResults notifies admins of successes and failures while trying to create users.
func emailResults(subject string, body string) {
	log.Println(subject)
	if mg == nil {
		mg = mailgun.NewMailgun(mgDomain, mgAPIKey, mgPublicAPIKey)
	}
	from := "no-reply@" + mgDomain
	subject = "userbot: " + subject
	to := os.Getenv("MG_ADMINS")
	message := mg.NewMessage(from, subject, body, to)
	_, _, err := mg.Send(message)
	if err != nil {
		log.Printf("mg.NewMessage: %s\n", err)
		log.Printf("from: %s\n", from)
		log.Printf("subject: %s\n", subject)
		log.Printf("body: %s\n", body)
		log.Printf("to: %s\n\n", to)
	}
}
