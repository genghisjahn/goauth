package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
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
	http.HandleFunc("/process", processHandler)
	http.ListenAndServe(":8090", nil)
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	sm := SignedMessage{}
	body, _ := ioutil.ReadAll(r.Body)
	jsonBody, _ := url.QueryUnescape(string(body))
	jsonBody = strings.Replace(jsonBody, "signedMsg=", "", 1)
	json.Unmarshal([]byte(jsonBody), &sm)
	log.Println("PublicKey: ", string(sm.Order.PublicKey))
	log.Println("Nonce: ", string(sm.Order.Nonce))
	log.Println("OrderDateTime: ", sm.Order.OrderDateTime)
	log.Println("Verb: ", sm.Order.Verb)
	log.Println("Hash ", sm.Hash)

}
