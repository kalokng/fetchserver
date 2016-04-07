package main

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
)

func hello(w http.ResponseWriter) {
	fmt.Println("Hello!")
	// just say hello...
	const avgSize = 2000
	const avgDev = 500
	len := int(rand.NormFloat64()*avgDev + avgSize)

	enc := base64.NewEncoder(base64.StdEncoding, w)

	const bufsize = 7 * 9

	// make a 7*9 bytes buffer
	buf := make([]byte, bufsize)
	for i := 0; i < len; i += bufsize {
		rand.Read(buf)
		enc.Write(buf)
	}
	enc.Close()
}
