package main

import (
	// "fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/userbot", parseEmail)
	err := http.ListenAndServeTLS(":8443", "userbot.crt", "userbot.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func parseEmail(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	log.Println(req.Header)
	log.Print('\n')
	log.Println(req.Body)
	log.Print('\n')
	w.WriteHeader(http.StatusOK)
}
