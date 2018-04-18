package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var wpBaseURL, wpPassword, wpUser string

func initWordPress() {
	wpBaseURL = fmt.Sprintf("https://%s/", os.Getenv("WP_BASE_URL"))
	wpPassword = os.Getenv("WP_PASSWORD")
	wpUser = os.Getenv("WP_USER")
}

func callWP(api, opts string) (*http.Response, error) {
	if wpBaseURL == "" {
		initWordPress()
	}
	url := wpBaseURL + "wp-json/wp/v2/" + api
	if opts != "" {
		url = url + "?" + opts
	}
	// TODO log.Printf("URL: %s\n", url)
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

func userExists(email string) (bool, error) {
	response, err := callWP("users", "search="+email)
	if err != nil {
		return false, err
	}
	body, _ := getResponseBody(response)
	// TODO log.Println(body)
	return strings.Contains(body, "id"), nil
}
