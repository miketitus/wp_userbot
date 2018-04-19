package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// User contains parsed user JSON returned by WordPress.
type User struct {
	Context     string `json:"context,omitempty"`
	ID          int32  `json:"id,omitempty"`
	Username    string `json:"username,omitempty"`
	DisplayName string `json:"name,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
}

var wpBaseURL, wpPassword, wpUser string

// initWordPress reads env vars for WordPress API.
func initWordPress() {
	wpBaseURL = fmt.Sprintf("https://%s/", os.Getenv("WP_BASE_URL"))
	wpPassword = os.Getenv("WP_PASSWORD")
	wpUser = os.Getenv("WP_USER")
}

// wpAPI sends a generic HTTP to the WordPress API.
func wpAPI(method, route, data string) (*http.Response, error) {
	var body io.Reader
	if wpBaseURL == "" {
		initWordPress()
	}
	url := wpBaseURL + "wp-json/wp/v2/" + route
	if method == "GET" && data != "" {
		url = url + "?" + data
	} else if method == "PUT" && data != "" {
		body = strings.NewReader(data)
	}
	log.Printf("url: %s\n", url)
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatal("http.NewRequest", err)
	}
	request.SetBasicAuth(wpUser, wpPassword)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	return client.Do(request)
}

// userExists determines whether a user account already exists (based on email address).
func userExists(email string) (bool, error) {
	response, err := wpAPI("GET", "users", "search="+email)
	if err != nil {
		log.Printf("client.Do: %s\n", response.Header)
		return false, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %s\n", err)
		return false, err
	}
	var users []User
	err = json.Unmarshal(body, &users)
	if err != nil {
		log.Printf("unmarshall: %s\n", err)
		return false, err
	}
	return len(users) > 0, nil
}

// TODO
func createUser(first, last, email string) (bool, error) {
	// ensure user doesn't already exist
	exists, err := userExists(email)
	if err != nil {
		return false, err
	} else if exists {
		return false, nil
	}
	// build user object
	user := User{FirstName: first, LastName: last, Email: email}
	user.Context = "edit" // TODO
	user.Username = strings.ToLower(first[:1] + last)
	user.DisplayName = first + ` ` + last
	user.Password = "aBcDeFgHiJ"
	j, err := json.Marshal(user)
	log.Printf("user: %s\n", j)
	// send user to WP
	response, err := wpAPI("PUT", "users", string(j))
	if err != nil {
		log.Printf("client.Do: %s\n", err)
		return false, err
	}
	log.Printf("header: %s\n", response.Header)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %s\n", err)
	}
	log.Printf("body: %s", body)
	return true, nil
}
