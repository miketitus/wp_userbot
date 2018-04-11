package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/userbot", parseEmail)
	err := http.ListenAndServe(":8443", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

func parseEmail(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal("ioutil.ReadAll", err)
	}
	body := string(bodyBytes)
	log.Println("********")
	log.Println(req.Header)
	log.Println("********")
	log.Println(body)
	log.Println("********")
	w.WriteHeader(http.StatusOK)
}
