package OpenPgpJsApi

type ActionRequest interface {
	Execute() (result RequestResult, err error)
}

type RequestResult interface {

}
