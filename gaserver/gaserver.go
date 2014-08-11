package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

type ReturnMessage struct {
	Message  string
	DateTime time.Time
	Success  bool
}

const PUBLIC_KEY_NOT_FOUND = "Public Key not found"
const DUPLICATE_NONCE = "Duplicate Nonce"
const INVALID_DATE = "Invalid Date"
const INVALID_HASH = "Invalid Hash"
const ORDER_SUCCESS = "Order processed successfully"

var pubkey = "mbRgpR2eYAdJkhvrfwjlmMC+L/0Vbrj4KvVo5nvnScwsx25LK+tPE3AM/IMcHuDW5zzp4Kup9xKd5YXupRJHzw=="
var privkey = "7F22ZeY+mlHtALq3sXcjrLdcID7whhVIQ5zD4bl4raKdBTYVgAjfdbvdfB5lmQa4wVP1o4frD5tfUcKON4ueVA=="

var nonceLog = map[string]time.Time{}

func main() {
	http.HandleFunc("/process", processHandler)
	http.ListenAndServe("192.168.1.7:8090", nil)
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	sm := SignedMessage{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	json.Unmarshal([]byte(body), &sm)
	sm.Order.Verb = strings.ToUpper(r.Method)
	rm := ProcessMessage(sm)
	if !rm.Success {
		if rm.Message == PUBLIC_KEY_NOT_FOUND {
			http.Error(w, rm.Message, http.StatusUnauthorized)
		} else {
			http.Error(w, rm.Message, http.StatusBadRequest)
		}
		log.Println(rm)
		return
	}
	log.Println(rm)
	rmJson, _ := json.Marshal(rm)
	w.Header().Set("Content-Type", "application/json")
	w.Write(rmJson)

}

func ProcessMessage(sm SignedMessage) ReturnMessage {
	rm := ReturnMessage{}
	rm.DateTime = time.Now().Local()

	/*
	   Check for a valid public key.
	   This would be stored in some kind of data repository.
	*/
	if string(sm.Order.PublicKey) != pubkey {
		rm.Message = PUBLIC_KEY_NOT_FOUND
		return rm
	}

	/*
		Make sure the nonce hasn't been used already.
		This prevents replay attacks.
	*/
	if !nonceLog[string(sm.Order.Nonce)].IsZero() {
		rm.Message = DUPLICATE_NONCE
		return rm
	}
	nonceLog[string(sm.Order.Nonce)] = time.Now().Local()

	/*
		Make sure that the request is within time contraints.
		This prevents delay attacks.
	*/
	duration := time.Since(sm.Order.OrderDateTime)
	if duration > 5*time.Second {
		rm.Message = INVALID_DATE
		return rm
	}

	//Calculate the hash using the server copy of the private key.
	sm2 := sm
	sm2.SetHash([]byte(privkey))

	/*
		Compare the two hashes.
		If they don't match, then something has changed in the request in transit
		and it is not a valid request.
	*/
	if sm.Hash != sm2.Hash {
		rm.Message = INVALID_HASH
	} else {
		//The request is valid.
		//ProcessOrder(sm)
		rm.Success = true
		rm.Message = ORDER_SUCCESS
	}

	return rm
}

func PrintNonceLog() {
	log.Println("Nonce Log:")
	for key, value := range nonceLog {
		fmt.Println("Key:", key, "Value:", value)
	}
}

func (sm *SignedMessage) SetHash(privkey []byte) {
	jsonbody, _ := json.Marshal(sm.Order)
	h := hmac.New(sha512.New, privkey)
	h.Write([]byte(jsonbody))
	sm.Hash = base64.StdEncoding.EncodeToString(h.Sum(nil))
}
