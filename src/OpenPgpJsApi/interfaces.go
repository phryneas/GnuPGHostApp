package OpenPgpJsApi

import "io"

type ActionRequest interface {
	Execute() (result RequestResult, err error)
}

type RequestResult interface {

}


type ReaderSeeker interface {
	io.Reader
	io.Seeker
}