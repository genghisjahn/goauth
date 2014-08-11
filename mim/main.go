package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/process", processHandler)
	http.ListenAndServe("192.168.1.50:8090", nil)
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	originalMsg := SignedMessage{}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	json.Unmarshal([]byte(body), &originalMsg)

	client := &http.Client{}

	remoteUrl := "http://192.168.1.7:8090/process"

	newMsg, _ := json.Marshal(originalMsg)

	req, _ := http.NewRequest("POST", remoteUrl, bytes.NewBufferString(string(newMsg)))
	resp, _ := client.Do(req)

	defer resp.Body.Close()
	contents, _ := ioutil.ReadAll(resp.Body)
	w.Header().Set("Content-Type", "application/json")
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		http.Error(w, string(contents), resp.StatusCode)
		return
	}
	w.Write(contents)

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
