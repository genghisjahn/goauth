package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
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
}

type SignedMessage struct {
	MessageInfo BaseMessage
	Hash        string
}

func (sm *SignedMessage) SetHash(privkey string) {
	jsonbody, _ := json.Marshal(sm.MessageInfo)
	key := []byte(privkey)
	h := hmac.New(sha512.New, key)
	h.Write([]byte(jsonbody))
	sm.Hash = base64.StdEncoding.EncodeToString(h.Sum(nil))
}
func main() {

	_, nonce := GenerateKey(32)

	sMsg := SignedMessage{}
	sMsg.MessageInfo.PublicKey = pubkey
	sMsg.MessageInfo.TimeStamp = time.Now().Local()
	sMsg.MessageInfo.Verb = "get"
	sMsg.MessageInfo.Nonce = nonce
	sMsg.MessageInfo.Body = "Get item id 42"
	sMsg.SetHash(privkey)
	jsonmsg, _ := json.Marshal(sMsg)
	fmt.Println("Hashed: ", sMsg.Hash)

	fmt.Printf("SignedMessage: %v.\n", string(jsonmsg))
	TestEqual(string(jsonmsg))
}

func TestEqual(jsonbody string) {
	var sm SignedMessage
	err := json.Unmarshal([]byte(jsonbody), &sm)
	if err != nil {
		fmt.Println("Something went wrong.")
	}
	fmt.Printf("\nReceived Hash: %v\n", sm.Hash)
	sm.SetHash(privkey)
	fmt.Printf("Computed Hash: %v\n", sm.Hash)

}

func GenerateKey(keylength int) ([]byte, string) {
	key := make([]byte, keylength)
	rand.Read(key)
	base64str := base64.StdEncoding.EncodeToString(key)
	return key, base64str
}
