package main

import (
	"os"
	"io"
	"log"
	"NativeMessagingHost"
	"fmt"
	"errors"
)

func main() {
	logger := log.New(os.Stderr, "GnuPGHostApp: ", log.Lshortfile)

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
			err =  r.(error)
		}
	}()

	request, err := NativeMessagingHost.ReadRequest(stdin)
	handleError(err)

	response := NativeMessagingHost.Response{}

	switch request.Action {
	case "decrypt":
		response.Data.Decrypt, err = request.Data.Decrypt.Execute()
	case "encrypt":
		response.Data.Encrypt, err = request.Data.Encrypt.Execute()
	default:
		err = errors.New(fmt.Sprintf("unknown action type '%s'", request.Action))
	}


	handleError(err)

	err = response.Send(stdout)
	handleError(err)

	return nil;
}
