package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestValidAuth(t *testing.T) {
	sMsg := getTestMessage("http://www.jonwear.com", "get", "Get item 42.")
	sMsg.SetHash([]byte(privkey))
	jsonmsg, _ := json.Marshal(sMsg)
	checkTest(t, true, []byte(privkey), string(jsonmsg), fmt.Sprintf("Supplied hash %v does not equal computed hash.\n", sMsg.Hash))
}

func TestInvalidURL(t *testing.T) {
	sMsg := getTestMessage("http://www.jonwear.com", "get", "Get item 42.")
	sMsg.SetHash([]byte(privkey))
	sMsg.MessageInfo.URL = "http://www.slashdot.org"
	jsonmsg, _ := json.Marshal(sMsg)
	checkTest(t, false, []byte(privkey), string(jsonmsg), fmt.Sprintf("Changed URL %v should have generated an invalid hash.\n", sMsg.MessageInfo.URL))
}

func TestInvalidVerb(t *testing.T) {
	sMsg := getTestMessage("http://www.jonwear.com", "get", "Get item 42.")
	sMsg.SetHash([]byte(privkey))
	sMsg.MessageInfo.Verb = "delete"
	jsonmsg, _ := json.Marshal(sMsg)
	checkTest(t, false, []byte(privkey), string(jsonmsg), fmt.Sprintf("Changed Verb %v should have generated an invalid hash.\n", sMsg.MessageInfo.Verb))
}

func TestInvalidMessage(t *testing.T) {
	sMsg := getTestMessage("http://www.jonwear.com", "get", "Get item 42.")
	sMsg.SetHash([]byte(privkey))
	sMsg.MessageInfo.Body = "That's not true!  That's imPOSSible!"
	jsonmsg, _ := json.Marshal(sMsg)
	checkTest(t, false, []byte(privkey), string(jsonmsg), fmt.Sprintf("Changed Message %v should have generated an invalid hash.\n", sMsg.MessageInfo.Verb))
}

func TestInvalidNonce(t *testing.T) {
	sMsg := getTestMessage("http://www.jonwear.com", "get", "Get item 42.")
	sMsg.SetHash([]byte(privkey))
	_, sMsg.MessageInfo.Nonce = GenerateKey(32)
	jsonmsg, _ := json.Marshal(sMsg)
	checkTest(t, false, []byte(privkey), string(jsonmsg), fmt.Sprintf("Changed nonce (%v) should have generated an invalid hash.", sMsg.MessageInfo.Nonce))
}

func TestInvalidPrivateKey(t *testing.T) {
	sMsg := getTestMessage("http://www.jonwear.com", "get", "Get item 42.")
	newprivkey, newkeystr := GenerateKey(64)
	sMsg.SetHash(newprivkey)
	jsonmsg, _ := json.Marshal(sMsg)
	checkTest(t, false, []byte(privkey), string(jsonmsg), fmt.Sprintf("Changed private key (%v) should have generated an invalid hash.", newkeystr))
}

func TestInvalidDateTime(t *testing.T) {
	sMsg := getTestMessage("http://www.jonwear.com", "get", "Get item 42.")
	sMsg.SetHash([]byte(privkey))
	sMsg.MessageInfo.TimeStamp = sMsg.MessageInfo.TimeStamp.Add(1 * time.Second)
	jsonmsg, _ := json.Marshal(sMsg)
	checkTest(t, false, []byte(privkey), string(jsonmsg), fmt.Sprintf("Changed datetime (%v) should have generated an invalid hash.", sMsg.MessageInfo.TimeStamp))
}

func checkTest(t *testing.T, desired bool, privkey []byte, jsonmsg string, failmessage string) {
	validtest, err := IsValid(privkey, jsonmsg)
	if validtest != desired {
		if err != nil {
			t.Error(err)
		}
		t.Errorf(failmessage)
	}
}

func getTestMessage(url string, verb string, messageBody string) SignedMessage {
	_, nonce := GenerateKey(32)
	sMsg := SignedMessage{}
	sMsg.MessageInfo.PublicKey = pubkey
	sMsg.MessageInfo.TimeStamp = time.Now().Local()
	sMsg.MessageInfo.Verb = verb
	sMsg.MessageInfo.URL = url
	sMsg.MessageInfo.Nonce = nonce
	sMsg.MessageInfo.Body = messageBody
	return sMsg
}
