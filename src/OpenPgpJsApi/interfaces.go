package OpenPgpJsApi

type RequestPackage interface {
	execute() (result RequestResult, err error)
}

type RequestResult interface {

}
