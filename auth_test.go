package main

import (
	"encoding/json"
	"testing"
	"time"
)

func TestValidAuth(t *testing.T) {
	_, nonce := GenerateKey(32)

	sMsg := SignedMessage{}
	sMsg.MessageInfo.PublicKey = pubkey
	sMsg.MessageInfo.TimeStamp = time.Now().Local()
	sMsg.MessageInfo.Verb = "get"
	sMsg.MessageInfo.Nonce = nonce
	sMsg.MessageInfo.Body = "Get item id 42"
	sMsg.SetHash([]byte(privkey))

	jsonmsg, _ := json.Marshal(sMsg)
	if !IsValid([]byte(privkey), string(jsonmsg)) {
		t.Errorf("Supplied hash %v does not equal computed hash.\n", sMsg.Hash)
	}
}

func GetTestMessage(messageBody string) SignedMessage {
	_, nonce := GenerateKey(32)

	sMsg := SignedMessage{}
	sMsg.MessageInfo.PublicKey = pubkey
	sMsg.MessageInfo.TimeStamp = time.Now().Local()
	sMsg.MessageInfo.Verb = "get"
	sMsg.MessageInfo.Nonce = nonce
	sMsg.MessageInfo.Body = messageBody
}

func TestInvalidNonce(t *testing.T) {
	sMsg := GetTestMessage("Get item 42.")
	sMsg.SetHash([]byte(privkey))
	_, sMsg.MessageInfo.Nonce = GenerateKey(32)
	jsonmsg, _ := json.Marshal(sMsg)
	if IsValid([]byte(privkey), string(jsonmsg)) {
		t.Errorf("Changed nonce (%v) should have generated an invalid hash.", sMsg.MessageInfo.Nonce)
	}
}

func TestInvalidPrivateKey(t *testing.T) {
	_, nonce := GenerateKey(32)

	sMsg := SignedMessage{}
	sMsg.MessageInfo.PublicKey = pubkey
	sMsg.MessageInfo.TimeStamp = time.Now().Local()
	sMsg.MessageInfo.Verb = "get"
	sMsg.MessageInfo.Nonce = nonce
	sMsg.MessageInfo.Body = "Get item id 42"
	newprivkey, _ := GenerateKey(64)
	sMsg.SetHash(newprivkey)
	jsonmsg, _ := json.Marshal(sMsg)
	if IsValid([]byte(privkey), string(jsonmsg)) {
		t.Errorf("Changed nonce (%v) should have generated an invalid hash.", sMsg.MessageInfo.Nonce)
	}
}
