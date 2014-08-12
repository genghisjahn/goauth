package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	inHttpAddr  = flag.String("inhttp", "192.168.1.50:8090", "Accept requests at this address.")
	outHttpAddr = flag.String("outhttp", "192.168.1.7:8090", "Send requests to this address")
)

func main() {
	flag.Parse()
	http.HandleFunc("/process", processHandler)
	log.Printf("Listening on: %v\n", *inHttpAddr)
	log.Printf("Sending to: %v\n", *outHttpAddr)
	http.ListenAndServe(*inHttpAddr, nil)
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	oMsg := GetOriginalMessage(w, r)
	_ = oMsg
	result, status := ProcessNormalOrder(oMsg)
	//result, status := ProcessChangedOrder(oMsg)
	//result, status := ProcessRepeatOrder(oMsg)
	//result, status := ProcessDelayOrder(oMsg)
	//result, status := ProcessDelayRepeatOrder(oMsg)
	//result, status := ProcessInvalidJson()

	w.Header().Set("Content-Type", "application/json")
	if status < 200 || status > 299 {
		http.Error(w, string(result), status)
		return
	}
	w.Write(result)
}

func ProcessNormalOrder(oMsg SignedMessage) ([]byte, int) {
	jsonNewMsg, _ := json.Marshal(oMsg)
	return processOrder(string(jsonNewMsg))
}

func ProcessInvalidJson() ([]byte, int) {
	return processOrder("INVALID JSON")
}

func ProcessChangedOrder(oMsg SignedMessage) ([]byte, int) {
	newMaxPrice := 1000
	newNumShares := 500
	log.Printf("Changing maxPrice from %v to %v\n", oMsg.Order.MaxPrice, newMaxPrice)
	log.Printf("Changing numShares from %v to %v\n", oMsg.Order.NumShares, newNumShares)
	oMsg.Order.MaxPrice = newMaxPrice
	oMsg.Order.NumShares = newNumShares
	jsonNewMsg, _ := json.Marshal(oMsg)
	return processOrder(string(jsonNewMsg))
}

func ProcessRepeatOrder(oMsg SignedMessage) ([]byte, int) {
	jsonNewMsg, _ := json.Marshal(oMsg)
	data1, code1 := processOrder(string(jsonNewMsg))
	time.Sleep(100 * time.Millisecond)
	processOrder(string(jsonNewMsg))
	return data1, code1

}

func ProcessDelayRepeatOrder(oMsg SignedMessage) ([]byte, int) {
	jsonNewMsg, _ := json.Marshal(oMsg)
	data1, code1 := processOrder(string(jsonNewMsg))
	go func() {
		time.Sleep(7000 * time.Millisecond)
		processOrder(string(jsonNewMsg))

	}()
	return data1, code1
}

func ProcessDelayOrder(oMsg SignedMessage) ([]byte, int) {
	jsonNewMsg, _ := json.Marshal(oMsg)
	time.Sleep(7000 * time.Millisecond)
	return processOrder(string(jsonNewMsg))
}

func GetOriginalMessage(w http.ResponseWriter, r *http.Request) SignedMessage {
	originalMsg := SignedMessage{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	json.Unmarshal([]byte(body), &originalMsg)
	return originalMsg
}

func processOrder(jsonNewMsg string) ([]byte, int) {
	client := &http.Client{}
	remoteUrl := fmt.Sprintf("http://%v/process", *outHttpAddr)
	req, _ := http.NewRequest("POST", remoteUrl, bytes.NewBufferString(jsonNewMsg))
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body)
	log.Println("Response : ", string(contents))
	return contents, resp.StatusCode
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

type ReturnMessage struct {
	Message  string
	DateTime time.Time
	Success  bool
}
