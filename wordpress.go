package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dghubble/oauth1"
	"github.com/robbiet480/go-wordpress"
)

var wpBaseURL, wpKey, wpSecret string
var wpClient *wordpress.Client
var wpConfig oauth1.Config
var wpContext context.Context

func initWordPress() {
	wpBaseURL = fmt.Sprintf("https://%s/", os.Getenv("WP_BASE_URL"))
	log.Println(wpBaseURL)
	wpKey = os.Getenv("WP_KEY")
	wpSecret = os.Getenv("WP_SECRET")
	wpContext = context.Background()
	initOAuth()
}

// http://73.254.169.106:8080/

func initOAuth() {
	wpConfig = oauth1.Config{
		ConsumerKey:    wpKey,
		ConsumerSecret: wpSecret,
		CallbackURL:    "http://73.254.169.106:8080",
		Endpoint: oauth1.Endpoint{
			RequestTokenURL: wpBaseURL + "oauth1/request",
			AuthorizeURL:    wpBaseURL + "oauth1/authorize",
			AccessTokenURL:  wpBaseURL + "oauth1/access",
		},
	}

	requestToken, requestSecret, err := login()
	if err != nil {
		log.Fatalf("Request Token Phase: %s", err.Error())
	}
	accessToken, err := receiveVerifier(requestToken, requestSecret)
	if err != nil {
		log.Fatalf("Access Token Phase: %s", err.Error())
	}

	log.Println("Consumer was granted an access token to act on behalf of a user.")
	log.Printf("token: %s\nsecret: %s\n", accessToken.Token, accessToken.TokenSecret)

}

func login() (requestToken, requestSecret string, err error) {
	requestToken, requestSecret, err = wpConfig.RequestToken()
	if err != nil {
		return "", "", err
	}
	authorizationURL, err := wpConfig.AuthorizationURL(requestToken)
	if err != nil {
		return "", "", err
	}
	fmt.Printf("Open this URL in your browser:\n%s\n", authorizationURL.String())
	return requestToken, requestSecret, err
}

func receiveVerifier(requestToken, requestSecret string) (*oauth1.Token, error) {
	fmt.Printf("Choose whether to grant the application access.\nPaste " +
		"the oauth_verifier parameter from the address bar: ")
	var verifier string
	_, err := fmt.Scanf("%s", &verifier)
	accessToken, accessSecret, err := wpConfig.AccessToken(requestToken, requestSecret, verifier)
	if err != nil {
		return nil, err
	}
	return oauth1.NewToken(accessToken, accessSecret), err
}
