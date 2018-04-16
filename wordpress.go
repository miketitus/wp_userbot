package main

import (
	"context"
	"fmt"
	"log"
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
