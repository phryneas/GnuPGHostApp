package OpenPgpJsApi

import (
	"github.com/pkg/errors"
)


func handleErr(err error) {
	if err != nil {
		panic(errors.WithStack(err))
	}
}

