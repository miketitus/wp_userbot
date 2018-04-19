package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// User contains parsed user JSON returned by WordPress
type User struct {
	ID          int32  `json:"id"`
	DisplayName string `json:"name"`
	Username    string `json:"slug"`
}

var wpBaseURL, wpPassword, wpUser string

// initWordPress reads env vars for WordPress API
func initWordPress() {
	wpBaseURL = fmt.Sprintf("https://%s/", os.Getenv("WP_BASE_URL"))
	wpPassword = os.Getenv("WP_PASSWORD")
	wpUser = os.Getenv("WP_USER")
}

// callWP sends a generic HTTP GET to the WordPress API
func callWP(api, opts string) (*http.Response, error) {
	if wpBaseURL == "" {
		initWordPress()
	}
	url := wpBaseURL + "wp-json/wp/v2/" + api
	if opts != "" {
		url = url + "?" + opts
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("http.NewRequest", err)
	}
	request.SetBasicAuth(wpUser, wpPassword)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	return client.Do(request)
}

// userExists determines whether a user account already exists (based on email address)
func userExists(email string) (bool, error) {
	response, err := callWP("users", "search="+email)
	if err != nil {
		return false, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll - %s\n", err)
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
	exists, err := userExists(email)
	if err != nil {
		return false, err
	} else if exists {
		return false, nil
	}
	return true, nil
}
