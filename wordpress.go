package main

import (
	"context"
	"log"
	"os"

	"github.com/robbiet480/go-wordpress"
)

var wpUser, wpPassword, wpBaseURL string
var wpClient *wordpress.Client
var wpContext context.Context

func initWordPress() {
	wpUser = os.Getenv("WP_USER")
	wpPassword = os.Getenv("WP_PASSWORD")
	wpBaseURL = os.Getenv("WP_BASE_URL")
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
