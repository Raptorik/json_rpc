package main

import (
	"bytes"
	"github.com/gorilla/rpc/json"
	"log"
	"net/http"
)

type Args struct {
	A, B int
}
type Result int

func main() {
	client := new(http.Client)
	url := "http://localhost:8081/rpc"
	args := Args{A: 2, B: 3}

	message, err := json.EncodeClientRequest("Arith.Multiply", args)
	if err != nil {
		log.Fatalf("%s", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(message))
	if err != nil {
		log.Fatalf("%s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error in sending request to %s. %s", url, err)
	}
	defer resp.Body.Close()

	var result Result
	err = json.DecodeClientResponse(resp.Body, &result)
	if err != nil {
		log.Fatalf("Couldn't decode response. %s", err)
	}
	log.Printf("%d*%d=%d\n", args.A, args.B, result)
}
