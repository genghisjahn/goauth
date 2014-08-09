package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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
	UUID     string
	DateTime time.Time
}

func main() {
	http.HandleFunc("/process", processHandler)
	http.ListenAndServe(":8090", nil)
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	sm := SignedMessage{}
	body, _ := ioutil.ReadAll(r.Body)
	jsonBody, _ := url.QueryUnescape(string(body))
	json.Unmarshal([]byte(jsonBody), &sm)

	rm := ReturnMessage{}
	rm.Message = "The reqeust was processed."
	rm.DateTime = time.Now().Local()
	rm.UUID = GenerateUUID()
	rmJson, _ := json.Marshal(rm)
	/*
		http.Error(w, "JW ERROR!", http.StatusInternalServerError)
		return
	*/
	w.Header().Set("Content-Type", "application/json")
	w.Write(rmJson)

}

func GenerateUUID() string {
	f, _ := os.Open("/dev/urandom")
	defer f.Close()

	b := make([]byte, 16)
	f.Read(b)

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
