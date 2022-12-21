package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Args struct {
	A, B string
}

type Words int

func (t *Words) Multiply(args *Args, reply *string) error {
	*reply = args.B + " " + args.A + "!"
	log.Printf(`Hello recrut: %s `, *reply)
	return nil
}

type rpcRequest struct {
	r    io.Reader
	rw   io.ReadWriter
	done chan bool
}

// NewRPCRequest returns a new rpcRequest.
func NewRPCRequest(r io.Reader) *rpcRequest {
	var buf bytes.Buffer
	done := make(chan bool)
	return &rpcRequest{r, &buf, done}
}

// Read implements the io.ReadWriteCloser Read method.
func (r *rpcRequest) Read(p []byte) (n int, err error) {
	return r.r.Read(p)
}

// Write implements the io.ReadWriteCloser Write method.
func (r *rpcRequest) Write(p []byte) (n int, err error) {
	return r.rw.Write(p)
}

// Close implements the io.ReadWriteCloser Close method.
func (r *rpcRequest) Close() error {
	r.done <- true
	return nil
}

// Call invokes the RPC request, waits for it to complete, and returns the results.
func (r *rpcRequest) Call() io.Reader {
	go jsonrpc.ServeConn(r)
	<-r.done
	return r.rw
}

func main() {
	words := new(Words)
	rpc.Register(words)
	http.HandleFunc("/concat", func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		res := NewRPCRequest(req.Body).Call()
		io.Copy(w, res)
	})
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
