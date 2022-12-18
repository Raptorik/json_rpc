package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {

	c := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(`GET`, `http://localhost:8081/rpc`, nil)
	if err != nil {
		fmt.Printf("Error: %s\\n", err)
		return
	}
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("Error: %s\\n", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("Body: %s\n", body)

}
