package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// User contains parsed user JSON returned by WordPress.
type User struct {
	ID        int32  `json:"id,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
}

var wpBaseURL, wpPassword, wpUser string

// initWordPress reads env vars for WordPress API.
func initWordPress() {
	wpBaseURL = fmt.Sprintf("https://%s/", os.Getenv("WP_BASE_URL"))
	wpPassword = os.Getenv("WP_PASSWORD")
	wpUser = os.Getenv("WP_USER")
	rand.Seed(time.Now().Unix())
}

// wpAPI sends a generic HTTP to the WordPress API.
func wpAPI(method, route, data string) (*http.Response, error) {
	var body io.Reader
	if wpBaseURL == "" {
		initWordPress()
	}
	route = wpBaseURL + "wp-json/wp/v2/" + route
	if data != "" {
		// query strings used for both GET and POST
		route = route + "?" + data
	}
	log.Printf("route: %s\n", route)
	request, err := http.NewRequest(method, route, body)
	if err != nil {
		log.Printf("http.NewRequest: %s\n", err)
		return nil, err
	}
	request.SetBasicAuth(wpUser, wpPassword)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	return client.Do(request)
}

// userExists determines whether a user account already exists (based on email address).
func userExists(email string) (bool, error) {
	response, err := wpAPI("GET", "users", "context=edit&search="+email)
	if err != nil {
		log.Printf("client.Do: %s\n", response.Header)
		return false, err
	}
	users := usersFromResponse(response)
	return len(users) > 0, nil
}

// usersFromResponse creates a slice of Users from body text returned by the WordPress API.
func usersFromResponse(response *http.Response) []User {
	var users []User
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %s\n", err)
		return users
	}
	// fmt.Printf("body: %s\n", body)
	err = json.Unmarshal(body, &users)
	if err != nil {
		log.Printf("json.Unmarshall: %s\n", err)
	}
	// fmt.Printf("users: %v\n", users)
	return users
}

// createUser creates a new user account (primary key is email).
// It parses the response to make sure creation was successful and
// returns the user id if successful, or -1 if unsuccessful.
func createUser(first, last, email string) (int, error) {
	// ensure user doesn't already exist
	exists, err := userExists(email)
	if err != nil {
		return -1, err
	} else if exists {
		msg := fmt.Sprintf("a user already exists with email: %s", email)
		return -1, errors.New(msg)
	}
	// build options string
	username := strings.ToLower(first[:1] + last)
	password := generatePassword(12)
	opts := fmt.Sprintf("username=%s&first_name=%s&last_name=%s&email=%s&password=%s",
		username, first, last, email, password)
	// send user to WP
	response, err := wpAPI("POST", "users", opts)
	if err != nil {
		log.Printf("wpAPI: %s\n", err)
		return -1, err
	}
	/* WP returns valid JSON upon user creation, but json.Unmarshall fails to parse
	it for unspecified reasons. So, instead of checking an actual user result like:
	users := usersFromResponse(response)
	We have to do a search on raw text like:
	*/
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %s\n", err)
		return -1, err
	}
	// get id of new user
	idRE := regexp.MustCompile("\"id\":([0-9]*),")
	ids := idRE.FindSubmatch(body)
	if len(ids) < 2 {
		// this should never happen
		return -1, nil
	}
	id := string(ids[1])
	id32, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		// this should never happen
		log.Printf("strconv.ParseUint: %s\n", err)
		return -1, err
	}
	emailUser(username, first, last, email, password)
	return int(id32), err
}

// generatePassword generates random passwords for new users.
// It omits special characters and characters that are easily confused with
// each other, depending on typeface: 'O' '0' 'l' '1'
func generatePassword(length uint8) string {
	if wpBaseURL == "" {
		initWordPress()
	}
	var charset = []byte("ABCDEFGHIJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789")
	var password = make([]byte, length)
	var i uint8
	for i = 0; i < length; i++ {
		r := rand.Intn(len(charset))
		password[i] = charset[r]
	}
	return string(password)
}

func deleteUser(id int) {
	if id <= 1 {
		// don't delete the admin!
		return
	}
	route := fmt.Sprintf("users/%d", id)
	response, err := wpAPI("DELETE", route, "force=true&reassign=1")
	if err != nil {
		log.Printf("wpAPI: %s\n", err)
		return
	}
	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %s\n", err)
		return
	}
	return
}
