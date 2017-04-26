package main

import (
	"os"
	"io"
	"log"
	"NativeMessagingHost"
	"github.com/proglottis/gpgme"
	"OpenPgpJsApi"
	"github.com/go-errors/errors"
	"fmt"
)

func main() {
	// for testing purposes, paste this string to STDIN:
	// XXXX{"action":"qwe","data":"ads"}'

	logger := log.New(os.Stderr, "GnuPGHostApp: ", log.Lshortfile)

	recipients, err := gpgme.FindKeys("e", false)
	if err != nil {
		logger.Fatalf("encountered error: %s", err)
	}
	for _, val := range recipients {
		logger.Printf("%s", val.UserIDs().Name())
	}

	mainLoop(logger)

}

func handleErrorNotEOF(err error) {
	if err == io.EOF  {
		return
	}
	handleError(err)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func mainLoop(logger *log.Logger) {
	for {
		e := LoopExecution(os.Stdin, os.Stdout)
		if e == io.EOF {
			logger.Println("encountered EOF, exiting gracefully")
		} else if e != nil {
			logger.Fatalf("encountered error: %s", e)
		}
	}
}

func LoopExecution(stdin io.Reader, stdout io.Writer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err =  errors.Wrap(r.(error), 1)
		}
	}()

	var request NativeMessagingHost.Request
	decoder, err := NativeMessagingHost.PrepareDecoder(stdin)
	handleError(err)

	decoder.Decode(&request)
	handleErrorNotEOF(err)

	decoder, err = NativeMessagingHost.PrepareDecoder(stdin)
	handleError(err)

	var actionRequest OpenPgpJsApi.ActionRequest
	switch request.Action {
	case "decrypt":
		var x OpenPgpJsApi.OpenPgpJsDecryptRequest
		err = decoder.Decode(&x)
		actionRequest = x
	case "encrypt":
		var x OpenPgpJsApi.OpenPgpJsEncryptRequest
		err = decoder.Decode(&x)
		actionRequest = x
	default:
		err = errors.New(fmt.Sprintf("unknown action type '%s'", request.Action))
	}
	handleErrorNotEOF(err)

	response, err := actionRequest.Execute()
	handleError(err)


	err = NativeMessagingHost.SendResponse(response, stdout)
	handleError(err)

	return nil;
}
