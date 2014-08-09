package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/genghisjahn/goauth"
)

var pubkey = "mbRgpR2eYAdJkhvrfwjlmMC+L/0Vbrj4KvVo5nvnScwsx25LK+tPE3AM/IMcHuDW5zzp4Kup9xKd5YXupRJHzw=="
var privkey = "7F22ZeY+mlHtALq3sXcjrLdcID7whhVIQ5zD4bl4raKdBTYVgAjfdbvdfB5lmQa4wVP1o4frD5tfUcKON4ueVA=="

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/send", sendHandler)
	http.ListenAndServe(":8080", nil)
}

type Page struct {
	Title string
	Label string
}

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

func homeHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Place an Order!", Label: "Demo"}
	t, _ := template.ParseFiles("template1.html")
	t.Execute(w, p)
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	posturl := "http://localhost:8090/process"

	numshares, _ := strconv.Atoi(r.FormValue("numshares"))
	maxprice, _ := strconv.Atoi(r.FormValue("maxprice"))
	order := BuildOrder(numshares, maxprice, "GET", posturl)
	signedMsg := SignedMessage{Order: order}

	signedMsg.SetHash([]byte(privkey))
	fmt.Fprintf(w, "HELLO")

}

func BuildOrder(numshares int, maxprice int, url string, verb string) OrderMessage {
	result := OrderMessage{}

	result.NumShares = numshares
	result.MaxPrice = maxprice

	result.PublicKey = []byte(pubkey)
	result.Nonce, _ = goauth.GenerateKey(32)
	result.OrderDateTime = time.Now().Local()
	result.Verb = verb
	result.URL = url
	return result
}

func (sm *SignedMessage) SetHash(privkey []byte) {
	jsonbody, _ := json.Marshal(sm.Order)
	h := hmac.New(sha512.New, privkey)
	h.Write([]byte(jsonbody))
	sm.Hash = base64.StdEncoding.EncodeToString(h.Sum(nil))
}
