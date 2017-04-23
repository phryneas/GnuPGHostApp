package main

import (
	"os"
	"io"
	"log"
	"NativeMessagingHost"
)

func main() {
	// for testing purposes, paste this string to STDIN:
	// XXXX{"action":"qwe","data":"ads"}'

	logger := log.New(os.Stderr, "GnuPGHostApp: ", log.Lshortfile)

	for {
		request, err := NativeMessagingHost.ReadRequest(os.Stdin)
		if err == io.EOF {
			logger.Println("encountered EOF, exiting gracefully")
			break // exit loop, cleanup might follow
		} else if err != nil {
			logger.Fatalf("encountered error: %s", err)
		}

		response := &NativeMessagingHost.Response{Status: "ok", Data: request.Data}
		err = response.Send(os.Stdout)
		if err != nil {
			logger.Fatalf("encountered error: %s", err)
		}
	}
}
