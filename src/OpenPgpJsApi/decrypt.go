package OpenPgpJsApi

import (
	"github.com/proglottis/gpgme"
	"io"
	"strings"
	"bytes"
	"github.com/pkg/errors"
)

// @see https://openpgpjs.org/openpgpjs/doc/module-openpgp.html#.decrypt

// unused options from the openpgp.js API:
// privateKey  //Key 	(optional) private key with decrypted secret key Data or session key
// sessionKey //Object 	(optional) session key in the form: { Data:Uint8Array, algorithm:String }
// password   string //String 	(optional) single password to decrypt the Message
//
const (
	UTF8 = "utf8"
	BINARY = "binary"
)

type DecryptRequest struct {
	DataString  string `json:"dataString"`
	DataBytes   []byte  `json:"dataBytes"`
	Format string `json:"format"`
	//String 	(optional) return Data Format either as 'utf8' or 'binary'
	// one of 'utf8' or 'binary'
	Signature string  `json:"signature"`
	//Signature 	(optional) Detached Signature for verification
	// passed as Armored String
	PublicKeys []string `json:"publicKeys"`
	//Key | Array.<Key> 	(optional) array of public keys or single key, to verify Signatures
}

type DecryptSignature struct {
	Keyid string `json:"keyid"`
	Valid bool `json:"valid"`
}

type DecryptResult struct {
	Signatures []DecryptSignature `json:"signatures"`
	DataString string `json:"dataString"`
	DataBytes  []uint8 `json:"dataBytes"`
}

func (r DecryptRequest) Execute() (result DecryptResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			result = DecryptResult{}
			err = r.(error)
		}
	}()

	ctx, err := gpgme.New()
	handleErr(err)
	defer ctx.Release()


	var messageReader ReaderSeeker
	if (r.DataBytes != nil){
		messageReader = bytes.NewReader(r.DataBytes)
	} else {
		messageReader = strings.NewReader(r.DataString)
	}
	message, err := gpgme.NewDataReader(messageReader)
	handleErr(err)
	defer message.Close()

	plain, err := gpgme.NewData()
	handleErr(err)
	defer plain.Close()

	var signatures []gpgme.Signature
	if r.Signature != "" {
		signature, err := gpgme.NewDataReader(strings.NewReader(r.Signature))
		handleErr(err)
		defer signature.Close()

		_, signatures, err = ctx.Verify(signature, message, nil)
		handleErr(err)

		messageReader.Seek(0, io.SeekStart)
		err = ctx.Decrypt(message, plain)
		handleErr(err)
	} else {
		_, signatures, err = ctx.Verify(message, nil, plain)
		handleErr(err)
	}

	for _, signature := range signatures {
		// TODO: hand out full signature object signatures and let javascript
		validity := (signature.Summary & gpgme.SigSumGreen) != 0
		result.Signatures = append(result.Signatures, DecryptSignature{Keyid: signature.Fingerprint, Valid: validity})
	}

	plain.Seek(0, gpgme.SeekSet)
	buf := new(bytes.Buffer)
	io.Copy(buf, plain)
	if r.Format == BINARY {
		result.DataBytes = buf.Bytes()
	} else if r.Format == UTF8 {
		result.DataString = buf.String()
	} else {
		handleErr(errors.New("unknown format"))
	}

	return
}
