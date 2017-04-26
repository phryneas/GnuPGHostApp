package OpenPgpJsApi

import "github.com/go-errors/errors"

func handleErr(err error) {
	if err != nil {
		panic(errors.Wrap(err.(error), 1))
	}
}
