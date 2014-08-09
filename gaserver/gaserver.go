package main

import (
	"fmt"
	"net/http"
	"time"
)

type OrderMessage struct {
	NumShares     int
	MaxPrice      int
	PublicKey     []byte
	Nonce         []byte
	OrderDateTime time.Time
	Verb          string
	URL           string
}

type SignedMessage struct {
	Hash  string
	Order OrderMessage
}

func main() {
	http.HandleFunc("/process/", processHandler)
	http.ListenAndServe(":8090", nil)
}

func processHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "HELLO")
}
