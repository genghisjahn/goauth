package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/process/", processHandler)
	http.ListenAndServe(":8090", nil)
}

func processHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Processing your order....")
}
