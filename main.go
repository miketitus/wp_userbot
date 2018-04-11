package main

import (
	"io/ioutil"
	"log"
	"net/QueryUnescape"
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
	bodyBytes, err1 := ioutil.ReadAll(req.Body)
	if err1 != nil {
		log.Fatal("ioutil.ReadAll", err1)
	}
	bodyString := string(bodyBytes)
	body, err2 := QueryUnescape(bodyString)
	if err2 != nil {
		log.Fatal("QueryUnescape", err2)
	}
	log.Println("********")
	log.Println(req.Header)
	log.Println("********")
	log.Println(bodyString)
	log.Println("********")
	w.WriteHeader(http.StatusOK)
}
