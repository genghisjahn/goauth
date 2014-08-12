package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"flag"
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

const (
	PUBLIC_KEY_NOT_FOUND = "Public Key not found"
	DUPLICATE_NONCE      = "Duplicate Nonce"
	EXPIRED_TIMESTAMP    = "Expired Timestamp"
	INVALID_HASH         = "Invalid Hash"
	INVALID_JSON         = "Invalid JSON"
	ORDER_SUCCESS        = "Order processed successfully | No. Shares %v | Max Price $%v |"
)

var (
	pubkey   = "mbRgpR2eYAdJkhvrfwjlmMC+L/0Vbrj4KvVo5nvnScwsx25LK+tPE3AM/IMcHuDW5zzp4Kup9xKd5YXupRJHzw=="
	privkey  = "7F22ZeY+mlHtALq3sXcjrLdcID7whhVIQ5zD4bl4raKdBTYVgAjfdbvdfB5lmQa4wVP1o4frD5tfUcKON4ueVA=="
	nonceLog = map[string]time.Time{}
	httpAddr = flag.String("http", "192.168.1.7:8090", "Listen address")
)

func main() {
	flag.Parse()
	go ClearNonces()
	http.HandleFunc("/process", processHandler)
	log.Printf("Listening on : %v\n", *httpAddr)
	http.ListenAndServe(*httpAddr, nil)
}

func ClearNonces() {
	for {
		for key, value := range nonceLog {
			duration := time.Since(value)
			if duration > 5*time.Second {
				log.Printf("EXPIRED NONCE %v  at %v\n", value, time.Now().Local())
				delete(nonceLog, key)
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
}
func processHandler(w http.ResponseWriter, r *http.Request) {
	sm := SignedMessage{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err := json.Unmarshal([]byte(body), &sm); err != nil {
		http.Error(w, INVALID_JSON, http.StatusBadRequest)
		log.Println(INVALID_JSON)
		return
	}
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
	if duration > 3*time.Second {
		rm.Message = EXPIRED_TIMESTAMP
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
		rm.Message = fmt.Sprintf(ORDER_SUCCESS, sm.Order.NumShares, sm.Order.MaxPrice)
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
