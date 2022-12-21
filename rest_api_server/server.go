package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type BStruct struct {
	Jsonrpc float64 `json:"jsonrpc"`
	Method  string  `json:"method"`
	Params  []struct {
		A string `json:"A"`
		B string `json:"B"`
	} `json:"params"`
}

type BodyResponse struct {
	Jsonrpc float64 `json:"jsonrpc"`
	Result  string  `json:"result"`
	Id      string  `json:"id"`
}

func doSomething(rw http.ResponseWriter, r *http.Request) {
	var vars map[string]string = mux.Vars(r)
	var occupation string = vars["occupation"]
	log.Printf(`So you want to be in a %s`, occupation)
	resp := make(map[string]string)
	occupation = occupation + " :) "
	log.Printf(`Your parents will honour if you will be in a %s`, occupation)

	respFrom2server := startServer(occupation)

	log.Printf(`Response from JSON-RPS SERVER: %s`, respFrom2server)

	respfromRest := respFrom2server + "That is a very hard decision of you!"
	resp["message"] = respfromRest
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	_, err = rw.Write(jsonResp)
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
}

func startServer(occupation string) string {
	c := http.Client{Timeout: time.Second}
	req, err := http.NewRequest(`POST`, `http://localhost:8080/concat`, nil)
	if err != nil {
		fmt.Printf("Error: %s\\n", err)
		return "Something wrong with your request"
	}
	body := `{"jsonrpc": 2.0,"method": "Words.Multiply","params": [{"A": "num", "B": "num"}]}`
	bodytobyte := []byte(body)
	var bodyStruct BStruct
	json.Unmarshal(bodytobyte, &bodyStruct)
	bodyStruct.Params[0].A = occupation
	bodyStruct.Params[0].B = "man"
	resultBody, _ := json.Marshal(&BStruct{
		Jsonrpc: 2.0,
		Method:  "Words.Multiply",
		Params:  bodyStruct.Params,
	})
	req.Header.Add(`Content-Type`, `application/json`)
	neededBody := io.NopCloser(bytes.NewReader(resultBody))
	req.Body = neededBody

	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("Error: %s\\n", err)
		return "Error with request from REST server to JSON_RPC server"
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var bodyResponse BodyResponse
	json.Unmarshal(bodyBytes, &bodyResponse)
	defer resp.Body.Close()
	return string(bodyResponse.Result)
}

func main() {
	ctx := context.Background()
	router := mux.NewRouter()
	srv := &http.Server{
		Addr:              `0.0.0.0:5000`,
		ReadTimeout:       time.Millisecond * 200,
		WriteTimeout:      time.Millisecond * 200,
		IdleTimeout:       time.Second * 10,
		ReadHeaderTimeout: time.Millisecond * 200,
		Handler:           router,
	}
	// "localhost:8081/{military}"
	router.HandleFunc("/{occupation}", doSomething)

	go func() {
		log.Println(`Web Server started`)
		err := srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	<-done

	log.Println(`Web Server is shutting down`)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(ctx, err)
	}

}
