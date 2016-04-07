package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
)

func read(p []byte) (n int, err error) {
	for i := 0; i < len(p); i += 7 {
		val := rand.Int63()
		for j := 0; i+j < len(p) && j < 7; j++ {
			p[i+j] = byte(val)
			val >>= 8
		}
	}
	return len(p), nil
}

func testHeader(w http.ResponseWriter, req *http.Request) {
	b, err := httputil.DumpRequest(req, true)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s\n%v", err, req), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func hello(w http.ResponseWriter, req *http.Request) {
	log.Println("Hello!")

	if req.URL.Path == "/testheader" {
		testHeader(w, req)
		return
	}

	// just say hello...
	const avgSize = 2000
	const avgDev = 500
	len := int(rand.NormFloat64()*avgDev + avgSize)

	enc := base64.NewEncoder(base64.StdEncoding, w)

	const bufsize = 7 * 9

	// make a 7*9 bytes buffer
	buf := make([]byte, bufsize)
	for i := 0; i < len; i += bufsize {
		read(buf)
		enc.Write(buf)
	}
	enc.Close()
}
