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
	outHttpAddr = flag.String("outhttp", "192.168.1.7:8090", "Send requests to this address.")
	attack      = flag.String("attack", "none", "What type of MiM attack to run.  Default is none.")

	MimAttacks = map[string]fnMimAttack{
		"none":        ProcessNormalOrder,
		"changeorder": ProcessChangedOrder,
		"repeat":      ProcessRepeatOrder,
		"delay":       ProcessDelayOrder,
		"delayrepeat": ProcessDelayRepeatOrder,
		"changeurl":   ProcessChangedURL,
		"changeverb":  ProcessChangedVerb,
		"invalid":     ProcessInvalidJson,
	}
)

type fnMimAttack func(oMsg SignedMessage) ([]byte, int)

func init() {
	flag.Parse()
	if MimAttacks[*attack] == nil {
		log.Printf("Valid Values are:\n")
		for key, _ := range MimAttacks {
			log.Println(key)
		}
		log.Fatalf("MiMAttack %v does not exist.", *attack)
	}
}
func main() {

	http.HandleFunc("/process", processHandler)
	log.Printf("Listening on: %v\n", *inHttpAddr)
	log.Printf("Sending to: %v\n", *outHttpAddr)
	http.ListenAndServe(*inHttpAddr, nil)
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	oMsg := GetOriginalMessage(w, r)

	fnAttack := MimAttacks[*attack]

	result, status := ProcessAttack(oMsg, fnAttack)
	w.Header().Set("Content-Type", "application/json")
	if status < 200 || status > 299 {
		http.Error(w, string(result), status)
		return
	}
	w.Write(result)
}

func ProcessAttack(oMsg SignedMessage, fnAttack fnMimAttack) ([]byte, int) {
	return (fnAttack(oMsg))
}

func ProcessNormalOrder(oMsg SignedMessage) ([]byte, int) {
	jsonNewMsg, _ := json.Marshal(oMsg)
	return processOrder(string(jsonNewMsg), oMsg.Order.Verb)
}

func ProcessInvalidJson(oMsg SignedMessage) ([]byte, int) {
	return processOrder("INVALID JSON", oMsg.Order.Verb)
}

func ProcessChangedOrder(oMsg SignedMessage) ([]byte, int) {
	newMaxPrice := 1000
	newNumShares := 500
	log.Printf("Changing maxPrice from %v to %v\n", oMsg.Order.MaxPrice, newMaxPrice)
	log.Printf("Changing numShares from %v to %v\n", oMsg.Order.NumShares, newNumShares)
	oMsg.Order.MaxPrice = newMaxPrice
	oMsg.Order.NumShares = newNumShares
	jsonNewMsg, _ := json.Marshal(oMsg)
	return processOrder(string(jsonNewMsg), oMsg.Order.Verb)
}

func ProcessChangedURL(oMsg SignedMessage) ([]byte, int) {
	newURL := "www.order-live.com:8090"
	log.Printf("Changing URL from %v to %v\n", oMsg.Order.URL, newURL)
	oMsg.Order.URL = newURL
	jsonNewMsg, _ := json.Marshal(oMsg)
	return processOrder(string(jsonNewMsg), oMsg.Order.Verb)
}

func ProcessChangedVerb(oMsg SignedMessage) ([]byte, int) {
	newVerb := "DELETE"
	log.Printf("Changing Method/Verb from %v to %v\n", oMsg.Order.Verb, newVerb)
	oMsg.Order.Verb = newVerb
	jsonNewMsg, _ := json.Marshal(oMsg)
	return processOrder(string(jsonNewMsg), oMsg.Order.Verb)
}

func ProcessRepeatOrder(oMsg SignedMessage) ([]byte, int) {
	jsonNewMsg, _ := json.Marshal(oMsg)
	data1, code1 := processOrder(string(jsonNewMsg), oMsg.Order.Verb)
	time.Sleep(100 * time.Millisecond)
	processOrder(string(jsonNewMsg), oMsg.Order.Verb)
	return data1, code1

}

func ProcessDelayRepeatOrder(oMsg SignedMessage) ([]byte, int) {
	jsonNewMsg, _ := json.Marshal(oMsg)
	data1, code1 := processOrder(string(jsonNewMsg), oMsg.Order.Verb)
	go func() {
		time.Sleep(7000 * time.Millisecond)
		processOrder(string(jsonNewMsg), oMsg.Order.Verb)

	}()
	return data1, code1
}

func ProcessDelayOrder(oMsg SignedMessage) ([]byte, int) {
	jsonNewMsg, _ := json.Marshal(oMsg)
	time.Sleep(7000 * time.Millisecond)
	return processOrder(string(jsonNewMsg), oMsg.Order.Verb)
}

func GetOriginalMessage(w http.ResponseWriter, r *http.Request) SignedMessage {
	originalMsg := SignedMessage{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	json.Unmarshal([]byte(body), &originalMsg)
	return originalMsg
}

func processOrder(jsonNewMsg string, verb string) ([]byte, int) {
	client := &http.Client{}
	remoteUrl := fmt.Sprintf("http://%v/process", *outHttpAddr)
	req, _ := http.NewRequest(verb, remoteUrl, bytes.NewBufferString(jsonNewMsg))
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
