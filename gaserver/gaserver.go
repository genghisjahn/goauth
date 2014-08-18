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

type RequestPayload struct {
	Data
	PublicKey       []byte
	Nonce           []byte
	RequestDateTime time.Time
	Verb            string
	URL             string
}

type Data interface{}

type SignedRequest struct {
	Hash    string
	Payload RequestPayload
}

type ResponseMessage struct {
	Message      string
	ResponseData Data
	DateTime     time.Time
	Success      bool
}

const (
	PUBLIC_KEY_NOT_FOUND = "Public Key not found"
	DUPLICATE_NONCE      = "Duplicate Nonce"
	EXPIRED_TIMESTAMP    = "Expired Timestamp"
	INVALID_HASH         = "Invalid Hash"
	INVALID_JSON         = "Invalid JSON"
	INVALID_URL          = "Invalid URL %v"
	INVALID_ORDER        = "NumShares and MaxPrice must be integers greater than 0"
	ORDER_SUCCESS        = "Order processed successfully"
)

var (
	//These keys are purely for demonstration purposes.  They won't provide access to anything, anywhere at anytime.
	pubkey   = "mbRgpR2eYAdJkhvrfwjlmMC+L/0Vbrj4KvVo5nvnScwsx25LK+tPE3AM/IMcHuDW5zzp4Kup9xKd5YXupRJHzw=="
	privkey  = "7F22ZeY+mlHtALq3sXcjrLdcID7whhVIQ5zD4bl4raKdBTYVgAjfdbvdfB5lmQa4wVP1o4frD5tmUcKON4ueVA=="
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
	sr := SignedRequest{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err := json.Unmarshal([]byte(body), &sr); err != nil {
		http.Error(w, INVALID_JSON, http.StatusBadRequest)
		log.Println(INVALID_JSON)
		return
	}
	sr.Payload.Verb = strings.ToUpper(r.Method)

	rm := ProcessRequest(sr)
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

func ProcessRequest(sr SignedRequest) ResponseMessage {
	rm := ResponseMessage{}
	rm.DateTime = time.Now().Local()

	/* Check for the correct URL
	This demo should only process requests mean for: www.order-demo.com:8090
	*/
	if sr.Payload.URL != "http://www.order-demo.com:8090/process" {
		rm.Message = fmt.Sprintf(INVALID_URL, sr.Payload.URL)
		return rm
	}

	/*
	   Check for a valid public key.
	   This would be stored in some kind of data repository.
	*/
	if string(sr.Payload.PublicKey) != pubkey {
		rm.Message = PUBLIC_KEY_NOT_FOUND
		return rm
	}

	/*
		Make sure the nonce hasn't been used already.
		This prevents replay attacks.
	*/
	if !nonceLog[string(sr.Payload.Nonce)].IsZero() {
		rm.Message = DUPLICATE_NONCE
		return rm
	}
	nonceLog[string(sr.Payload.Nonce)] = time.Now().Local()

	/*
		Make sure that the request is within time contraints.
		This prevents delay attacks.
	*/
	duration := time.Since(sr.Payload.RequestDateTime)
	if duration > 3*time.Second {
		rm.Message = EXPIRED_TIMESTAMP
		return rm
	}

	//Calculate the hash using the server copy of the private key.
	sr2 := sr
	sr2.SetHash([]byte(privkey))

	/*Check if both NumShares and MaxPrice are both > 0.
	If either of the values are 0, then most likely non integer data was submitted.
	No matter the cause, the order should not be processed.
	*/
	/*
		TODO make this work with the empty interface.
		Check to see if the JSON changed.
		if sm.Order.MaxPrice <= 0 || sm.Order.NumShares <= 0 {
			rm.Message = INVALID_ORDER
			return rm
		}
	*/
	/*
		Compare the two hashes.
		If they don't match, then something has changed in the request in transit
		and it is not a valid request.
	*/
	if sr.Hash != sr2.Hash {
		rm.Message = INVALID_HASH
	} else {
		//The request is valid.
		rm.Success = true
		rm.ResponseData = sr.Payload.Data
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

func (sr *SignedRequest) SetHash(privkey []byte) {
	jsonbody, _ := json.Marshal(sr.Payload)
	h := hmac.New(sha512.New, privkey)
	h.Write([]byte(jsonbody))
	sr.Hash = base64.StdEncoding.EncodeToString(h.Sum(nil))
}
