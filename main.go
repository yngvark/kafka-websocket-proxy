package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yngvark/go-rest-api-template/pkg/helloworld"
)

func main() {
	fmt.Println(helloworld.Hello())

	port := 8080
	fmt.Printf("Running web server on port %v\n", port)

	handleRequests(port)
}

func handleRequests(port int) {
	http.HandleFunc("/", root)
	http.HandleFunc("/health", health)

	address := fmt.Sprintf("%v%v", ":", port)
	log.Fatal(http.ListenAndServe(address, nil))
}

func root(w http.ResponseWriter, req *http.Request) {
	// The "/" pattern matches everything, so we need to check
	// that we're at the root here.
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}

	fmt.Println("Request to /")
	_, _ = fmt.Fprintln(w, helloworld.Hello())
}

func health(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "OK")
}
