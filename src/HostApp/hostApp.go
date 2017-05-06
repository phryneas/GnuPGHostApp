package main

import (
	"os"
	"io"
	"log"
	"NativeMessagingHost"
	"fmt"
	"time"
	"github.com/pkg/errors"
)

var logger *log.Logger

func main() {
	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		os.Stderr.WriteString("could not open logfile\n")
		os.Exit(1)
	}
	//logger := log.New(os.Stderr, "GnuPGHostApp: ", log.Lshortfile)
	logger = log.New(f, "GnuPGHostApp: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Printf("initialized with PATH %s", os.Getenv("PATH"))

	mainLoop()
}

func mainLoop() {
	for {
		e := LoopExecution(os.Stdin, os.Stdout)
		if errors.Cause(e) == io.EOF {
			time.Sleep(time.Millisecond * 100)
			continue;
		} else if e != nil {
			logger.Fatalf("encountered error: %s", e)
		}
	}
}

func LoopExecution(stdin io.Reader, stdout io.Writer) (err error) {
	request, err := NativeMessagingHost.ReadRequest(stdin)
	if err != nil {
		return;
	}
	logger.Printf("received request: '%s'", request)

	response := NativeMessagingHost.Response{}

	switch request.Action {
	case "decrypt":
		response.Data.Decrypt, err = request.Data.Decrypt.Execute()
	case "encrypt":
		response.Data.Encrypt, err = request.Data.Encrypt.Execute()
	case "findKeys":
		response.Data.FindKeys, err = request.Data.FindKeys.Execute()
	case "test":
		response.Message = "test"
	default:
		err = errors.New(fmt.Sprintf("unknown action type '%s'", request.Action))
	}
	if response.Status == "" && err == nil {
		response.Status = "ok"
	} else if err != nil {
		response.Status = "error"
		response.Message = fmt.Sprintf("%+v", err)
	}

	logger.Printf("sending response: '%s'", response)
	err = response.Send(stdout)
	return err;
}
