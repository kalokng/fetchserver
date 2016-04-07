package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kalokng/fetch/obj"

	"golang.org/x/net/websocket"
)

func Reply(ws *websocket.Conn, resp *http.Response) {
	enc := gob.NewEncoder(ws)
	re := obj.NewResponse(resp)
	err := enc.Encode(re)
	if err != nil {
		log.Println(err)
	}
}

// Echo the data received on the WebSocket.
func EchoServer(ws *websocket.Conn) {
	dec := gob.NewDecoder(ws)
	var req obj.Request
	err := dec.Decode(&req)
	if err != nil {
		log.Println("fail decode", err)
		Reply(ws, &http.Response{})
		return
	}
	resp, err := http.DefaultClient.Do(req.HTTPRequest())
	if err != nil {
		log.Println("fail request", err)
		Reply(ws, &http.Response{})
		return
	}
	Reply(ws, resp)
	//n, err := io.Copy(ws, ws)
	//fmt.Println(n, err)
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	header := req.Header
	proto := header.Get("x-forwarded-proto")
	switch proto {
	case "ws", "wss":
	default:
		hello(w, req)
		return
	}
	ws := header.Get("hello")
	if ws != "world" {
		hello(w, req)
		return
	}
	websocket.Handler(EchoServer).ServeHTTP(w, req)
}

func getIP() string {
	return os.Getenv("OPENSHIFT_GO_IP")
}

func getPort() string {
	s := os.Getenv("OPENSHIFT_GO_PORT")
	if s == "" {
		return "8080"
	}
	return s
}

// This example demonstrates a trivial echo server.
func main() {
	bind := fmt.Sprintf("%s:%s", getIP(), getPort())

	err := http.ListenAndServe(bind, http.HandlerFunc(mainHandle))
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
