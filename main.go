package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"time"
)

/*
	Just hold a store of
	username,pubkey,privkey,email
	We can work out user accout creation later


*/

/* These are just for demo purposes.  They don't work for any external service. */

var pubkey = "mbRgpR2eYAdJkhvrfwjlmMC+L/0Vbrj4KvVo5nvnScwsx25LK+tPE3AM/IMcHuDW5zzp4Kup9xKd5YXupRJHzw=="
var privkey = "7F22ZeY+mlHtALq3sXcjrLdcID7whhVIQ5zD4bl4raKdBTYVgAjfdbvdfB5lmQa4wVP1o4frD5tfUcKON4ueVA=="

type BaseMessage struct {
	PublicKey string
	TimeStamp time.Time
	Nonce     string
	Verb      string
	Body      string
	URL       string
}

type SignedMessage struct {
	MessageInfo BaseMessage
	Hash        string
}

func (sm *SignedMessage) SetHash(privkey []byte) {
	jsonbody, _ := json.Marshal(sm.MessageInfo)
	h := hmac.New(sha512.New, privkey)
	h.Write([]byte(jsonbody))
	sm.Hash = base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func main() {

}

func IsValid(privatekey []byte, jsonbody string) (bool, error) {
	var sm SignedMessage
	err := json.Unmarshal([]byte(jsonbody), &sm)
	if err != nil {
		return false, err
	}
	hash1 := sm.Hash
	sm.SetHash([]byte(privkey))
	hash2 := sm.Hash
	return hash1 == hash2, nil

}

func GenerateKey(keylength int) ([]byte, string) {
	key := make([]byte, keylength)
	rand.Read(key)
	base64str := base64.StdEncoding.EncodeToString(key)
	return key, base64str
}
