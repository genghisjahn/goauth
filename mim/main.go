package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/process", processHandler)
	http.ListenAndServe("192.168.1.50:8090", nil)
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Har!  I'm the Man in the Middle!")
}
