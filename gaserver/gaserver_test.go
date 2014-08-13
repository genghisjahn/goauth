package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestValidPublicKey(t *testing.T) {
	sMsg := GetMessage(5, 5)
	sMsg.SetHash([]byte(privkey))
	jsonmsg, _ := json.Marshal(sMsg)
	request, _ := http.NewRequest("POST", "/", bytes.NewBufferString(string(jsonmsg)))
	response := httptest.NewRecorder()

	processHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("Response body did not contain expected %v:\n\tbody: %v", "200", response.Code)
	}
}

func TestInvalidPublicKey(t *testing.T) {
	sMsg := GetMessage(0, 0)
	sMsg.Order.PublicKey, _ = GenerateKey(64)
	sMsg.SetHash([]byte(privkey))
	jsonmsg, _ := json.Marshal(sMsg)
	request, _ := http.NewRequest("POST", "/", bytes.NewBufferString(string(jsonmsg)))
	response := httptest.NewRecorder()

	processHandler(response, request)

	if response.Code == http.StatusOK {
		t.Fatalf("Response body did not contain expected %v:\n\tbody: %v", "400", response.Code)
	}
}

func TestInvalidOrder(t *testing.T) {
	sMsg := GetMessage(0, 0)
	sMsg.SetHash([]byte(privkey))
	jsonmsg, _ := json.Marshal(sMsg)
	request, _ := http.NewRequest("POST", "/", bytes.NewBufferString(string(jsonmsg)))
	response := httptest.NewRecorder()

	processHandler(response, request)

	if response.Code == http.StatusOK {
		t.Fatalf("Response body did not contain expected %v:\n\tbody: %v", "400", response.Code)
	}
}

func GetMessage(numshares int, maxprice int) SignedMessage {
	nonce, _ := GenerateKey(32)
	result := SignedMessage{}
	result.Order.NumShares = numshares
	result.Order.MaxPrice = maxprice
	result.Order.PublicKey = []byte(pubkey)
	result.Order.Nonce = nonce
	result.Order.OrderDateTime = time.Now().Local()
	result.Order.Verb = "POST"
	result.Order.URL = "http://www.order-demo.com:8090/process"
	return result
}

func GenerateKey(keylength int) ([]byte, string) {
	key := make([]byte, keylength)
	rand.Read(key)
	base64str := base64.StdEncoding.EncodeToString(key)
	return key, base64str
}
