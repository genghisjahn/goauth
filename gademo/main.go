package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/genghisjahn/goauth"
)

var (
	pubkey   = "mbRgpR2eYAdJkhvrfwjlmMC+L/0Vbrj4KvVo5nvnScwsx25LK+tPE3AM/IMcHuDW5zzp4Kup9xKd5YXupRJHzw=="
	privkey  = "7F22ZeY+mlHtALq3sXcjrLdcID7whhVIQ5zD4bl4raKdBTYVgAjfdbvdfB5lmQa4wVP1o4frD5tmUcKON4ueVA=="
	httpAddr = flag.String("http", "www.order-demo.com:8090", "Server address")
)

func main() {
	flag.Parse()
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/send", sendHandler)
	log.Printf("Remote requests will be sent to %v\n", *httpAddr)
	log.Printf("Listening on localhost:8080")
	http.ListenAndServe("localhost:8080", nil)

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
	remoteUrl := fmt.Sprintf("http://%v/process", *httpAddr)
	numshares, _ := strconv.Atoi(r.FormValue("numshares"))
	maxprice, _ := strconv.Atoi(r.FormValue("maxprice"))
	order := BuildOrder(numshares, maxprice, remoteUrl, "POST")
	signedMsg := SignedMessage{Order: order}

	signedMsg.SetHash([]byte(privkey))
	sm, _ := json.Marshal(signedMsg)
	client := &http.Client{}

	req, _ := http.NewRequest("POST", remoteUrl, bytes.NewBufferString(string(sm)))
	resp, _ := client.Do(req)
	fmt.Println(resp.Status)
	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body)
	w.Header().Set("Content-Type", "application/json")
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		http.Error(w, string(contents), resp.StatusCode)
		return
	}
	w.Write(contents)
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
