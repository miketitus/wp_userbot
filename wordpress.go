package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/robbiet480/go-wordpress"
)

var wpBaseURL, wpPassword, wpUser string
var wpClient *wordpress.Client
var wpContext context.Context

func initWordPress() {
	wpBaseURL = fmt.Sprintf("https://%s/", os.Getenv("WP_BASE_URL"))
	wpPassword = os.Getenv("WP_PASSWORD")
	wpUser = os.Getenv("WP_USER")
	wpContext = context.Background()
	transport := wordpress.BasicAuthTransport{
		Username: wpUser,
		Password: wpPassword,
	}
	var err error
	wpClient, err = wordpress.NewClient(wpBaseURL, transport.Client())
	if err != nil {
		log.Fatal("wordpress.NewClient", err)
	}
}

/*func userExists(email string) bool {
	if wpBaseURL == "" {
		initWordPress()
	}
	users, resp, err := wpClient.Users.List(wpContext, "search="+email)
}*/

func callWP(api, opts string) (*http.Response, error) {
	wpBaseURL = fmt.Sprintf("https://%s/", os.Getenv("WP_BASE_URL"))
	wpPassword = os.Getenv("WP_PASSWORD")
	wpUser = os.Getenv("WP_USER")
	url := wpBaseURL + "wp-json/wp/v2/" + api
	if opts != "" {
		url = url + "?" + opts
	}
	log.Printf("URL: %s\n", url)
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
