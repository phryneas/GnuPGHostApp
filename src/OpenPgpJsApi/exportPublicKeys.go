package OpenPgpJsApi

import (
	"github.com/proglottis/gpgme"
	"bytes"
	"io"
	"strings"
)

type ExportPublicKeysRequest struct {
	Pattern string `json:"pattern"`
}

type ExportPublicKeysResult struct {
	KeyBlock string `json:"keyBlock"`
}

func (r ExportPublicKeysRequest) Execute() (result ExportPublicKeysResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			result = ExportPublicKeysResult{}
			err = r.(error)
		}
	}()

	ctx, err := gpgme.New()
	handleErr(err)
	defer ctx.Release()

	err = ctx.SetEngineInfo(gpgme.ProtocolOpenPGP, "", "")
	handleErr(err)

	ctx.SetArmor(true)

	plain, err := gpgme.NewData()
	handleErr(err)
	defer plain.Close()

	err = ctx.Export(r.Pattern, gpgme.ExportModeMinimal, plain)
	handleErr(err)

	plain.Seek(0, gpgme.SeekSet)
	buf := new(bytes.Buffer)
	io.Copy(buf, plain)
	result.KeyBlock = normalizeLines(buf.String())

	return
}

func normalizeLines(str string) string {
	return strings.Replace(strings.Replace(str, "\n\r", "\n", -1), "\r\n", "\n", -1);
}
