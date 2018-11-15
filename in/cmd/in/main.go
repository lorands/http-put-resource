package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	resource "github.com/lorands/http-resource"
	"github.com/lorands/http-resource/in"
)

func main() {
	var request in.Request

	fmt.Println("Request:", request)

	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		fatal("reading request from stdin", err)
	}

	fmt.Println("Request:", request)

	timestamp := request.Version.Timestamp
	if timestamp.IsZero() {
		timestamp = time.Now()
	}

	fmt.Println("TS:", timestamp)

	response := in.Response{
		Version: resource.Version{
			Timestamp: timestamp,
		},
	}
	fmt.Println("Response:", response)

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		fatal("writing response", err)
	}
	fmt.Println("Done.")
}

func fatal(message string, err error) {
	fmt.Fprintf(os.Stderr, "error %s: %s\n", message, err)
	os.Exit(1)
}
