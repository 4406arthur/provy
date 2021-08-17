// Simple server listening on a Unix socket and echoing back each client.
//
// With netcat, connect as follows:
//
// $ nc -U /tmp/echo.sock
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

const SockAddr = "/tmp/echo.sock"

func main() {
	if err := os.RemoveAll(SockAddr); err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("unix", SockAddr)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer l.Close()

	http.HandleFunc("/", HelloServer)
	if err := http.Serve(l, nil); err != nil {
		log.Fatal(err)
	}
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
