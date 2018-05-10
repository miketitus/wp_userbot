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
	"time"

	"gopkg.in/mailgun/mailgun-go.v1"
)

var mgAdmins []string
var mgAPIKey, mgDomain, mgListenPort, mgPublicAPIKey, mgUserBcc, mgUserBot string
var mg mailgun.Mailgun

// TODO HTTPS
func main() {
	initMain()
	// listen for email POSTs from Mailgun
	http.HandleFunc("/userbot", parseEmail)
	err := http.ListenAndServe(`:`+mgListenPort, nil)
	if err != nil {
		log.Fatal("http.ListenAndServe", err)
	}
}

// initMain reads env vars for Mailgun API.
func initMain() {
	mgAdmins = strings.Split(os.Getenv("MG_ADMINS"), ", ")
	mgAPIKey = os.Getenv("MG_API_KEY")
	mgDomain = os.Getenv("MG_DOMAIN")
	mgListenPort = os.Getenv("MG_LISTEN_PORT")
	mgPublicAPIKey = os.Getenv("MG_PUBLIC_API_KEY")
	mgUserBcc = os.Getenv("MG_USER_BCC")
	mgUserBot = os.Getenv("MG_USERBOT")
	mg = mailgun.NewMailgun(mgDomain, mgAPIKey, mgPublicAPIKey)
}

// parseEmail is the main event loop, executing for each received email.
func parseEmail(w http.ResponseWriter, request *http.Request) {
	log.Println("* * * start * * *")
	log.Printf("Header: %s\n", request.Header)
	// acknowledge POST from Mailgun
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(250) // SMTP OK
	// decode request
	body, err := getRequestBody(request)
	if err != nil {
		return
	}
	// validate request
	if !senderIsAdmin(body) {
		// ignore: spam, or a recipient hit "reply to all"
		emailResults("Illegal Sender", body)
		return
	}
	// process request
	parseRecipients(body)
	log.Println("* * * end * * *")
}

// getBody extracts and unescapes the email body from the Mailgun POST.
func getRequestBody(req *http.Request) (string, error) {
	// read body
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("getRequestBody: %s\n", err)
		emailResults("Parse Error", err.Error())
		return "", err
	}
	err = writeBody(bodyBytes)
	if err != nil {
		log.Printf("getRequestBody: %s\n", err)
		emailResults("writeBody Error", err.Error())
	}
	// decode body
	rawBody, err := url.QueryUnescape(string(bodyBytes))
	if err != nil {
		log.Printf("getRequestBody: %s\n", err)
		emailResults("Parse Error", err.Error())
		return "", err
	}
	return rawBody, nil
}

// writeBody writes raw body text to a temporary file for debugging.
func writeBody(body []byte) error {
	t := time.Now()
	fname := fmt.Sprintf("/tmp/%d.txt", t.Unix())
	return ioutil.WriteFile(fname, body, 0644)
}

// senderIsAdmin verifies that the decoded email message came from an approved email address.
func senderIsAdmin(body string) bool {
	var sender string
	senderRE := regexp.MustCompile("&from=([^&]*)")
	raw := senderRE.FindString(body)
	if raw == "" {
		// this should never happen, but let's keep on truckin'
		sender = raw
	} else {
		sender = raw[6:]
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

// parseRecipients decodes email recipients, and loops through them creating new
// WordPress users for each recipient, except for userbot and admin recipients.
func parseRecipients(body string) {
	var hadError bool
	var resultBody []string
	recipientRE := regexp.MustCompile("To=([^&]*)")
	raw := recipientRE.FindString(body)
	recipients := strings.Split(raw[3:], ", ")
	resultBody = append(resultBody, fmt.Sprintf("Recipient list: %s", recipients))
	for _, r := range recipients {
		fields := getFields(r)
		if isUserBot(fields) {
			continue // ignore
		} else if len(fields) == 3 {
			// valid structure
			if senderIsAdmin(fields[2]) {
				continue // admin email, skip
			} else if !isValidEmail(fields[2]) {
				hadError = true
				resultBody = append(resultBody, fmt.Sprintf("%s: Invalid email", r))
			} else {
				var result string
				id, err := createUser(fields[0], fields[1], fields[2])
				if err != nil {
					hadError = true
					result = err.Error()
				} else if id < 0 {
					hadError = true
					result = "Unknown Error"
				} else {
					result = "Success"
				}
				resultBody = append(resultBody, fmt.Sprintf("%s: %s", r, result))
			}
		} else {
			// error, invalid structure
			hadError = true
			resultBody = append(resultBody, fmt.Sprintf("%s: Invalid format", r))
		}
	}
	var resultSubject string
	if hadError {
		resultSubject = "Error(s) found"
	} else {
		resultSubject = "Success"
	}
	emailResults(resultSubject, strings.Join(resultBody, "\n"))
}

// getFields splits an email address into component strings, and cleans up the email address.
// e.g. "John Doe <john@john.doe>" --> [john doe john@john.doe]
func getFields(s string) []string {
	// first and/or last name is optional, but last field should always be the email field
	fields := strings.Fields(s)
	// cleanup email address
	i := len(fields) - 1
	if fields[i][0:1] == `<` {
		end := len(fields[i]) - 1
		fields[i] = fields[i][1:end]
	}
	return fields
}

// isValidEmail verifies that a recipient email address is in valid format.
// Uses a very simple regex designed to catch basic errors, but not nearly all edge cases.
func isValidEmail(email string) bool {
	emailRE := regexp.MustCompile("[^@]+@[^@]+\\..+")
	return emailRE.FindStringIndex(email) != nil
}

// emailResults notifies admins of processing successes and failures.
func emailResults(subject string, body string) {
	if mg == nil {
		initMain()
	}
	from := "no-reply@" + mgDomain
	subject = "userbot: " + subject
	to := os.Getenv("MG_ADMINS")
	message := mg.NewMessage(from, subject, body, to)
	_, _, err := mg.Send(message)
	if err != nil {
		log.Printf("emailResults: %s\n", err)
		log.Printf("from: %s\n", from)
		log.Printf("to: %s\n\n", to)
		log.Printf("subject: %s\n", subject)
		log.Printf("body: %s\n", body)
	}
}

// emailUser notifies the new user of their login username and password.
func emailUser(username, first, last, email, password string) {
	if mg == nil {
		initMain()
	}
	from := "no-reply@" + mgDomain
	subject := "NCWA forum login info" // TODO
	to := email
	plainBody := fmt.Sprintf(getPlainText(), first, last, username, password)
	htmlBody := fmt.Sprintf(getHTMLText(), first, last, username, password)
	message := mg.NewMessage(from, subject, plainBody, to)
	message.SetHtml(htmlBody)
	if mgUserBcc != "" {
		message.AddBCC(mgUserBcc)
	}
	_, _, err := mg.Send(message)
	if err != nil {
		log.Printf("emailUser: %s\n", err)
		log.Printf("from: %s\n", from)
		log.Printf("to: %s\n\n", to)
		log.Printf("subject: %s\n", subject)
		log.Printf("body: %s\n", plainBody)
	}
}

// getPlainText returns a plain email template for use with fmt.Sprintf()
func getPlainText() string {
	return `Hi %s %s,

	Here is the login info that will allow you to access the NCWA discussion forums.
	
	URL: https://ncwawood.org/wp-login.php
	Username: %s
	Password: %s
	
	Michael Titus
	NCWA Webmaster`
}

// getHTMLText returns an HTML email template for use with fmt.Sprintf()
func getHTMLText() string {
	return `<p>Hi %s %s,</p>

	<p>Here is the login info that will allow you to access the NCWA discussion forums.</p>
	
	<p>URL: <a href="https://ncwawood.org/wp-login.php">https://ncwawood.org/wp-login.php</a><br/>
	Username: %s<br/>
	Password: %s</p>
	
	<p>Michael Titus<br/>
	NCWA Webmaster</p>`
}
