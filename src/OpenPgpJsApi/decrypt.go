package OpenPgpJsApi

import (
	"github.com/proglottis/gpgme"
	"io"
	"strings"
	"bytes"
)

// unused options from the openpgp.js API:
// privateKey  //Key 	(optional) private key with decrypted secret key data or session key
// sessionKey //Object 	(optional) session key in the form: { data:Uint8Array, algorithm:String }
// password   string //String 	(optional) single password to decrypt the message
// publicKeys []string //Key | Array.<Key> 	(optional) array of public keys or single key, to verify signatures

type OpenPgpJsDecryptRequest struct {
	message string  `json:"message"`
	//Message 	the message object with the encrypted data
	// passed as Armored String
	format string `json:"format"`
	//String 	(optional) return data format either as 'utf8' or 'binary'
	// one of 'utf8' or 'binary'
	signature string  `json:"signature"`
	//Signature 	(optional) detached signature for verification
	// passed as Armored String
}

type OpenPgpJsDecryptSignature struct {
	keyid string
	valid bool
}

type OpenPgpJsDecryptResult struct {
	signatures []OpenPgpJsDecryptSignature
	dataString string
	dataBytes  []uint8
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (r OpenPgpJsDecryptRequest) Execute() (result OpenPgpJsDecryptResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			result = OpenPgpJsDecryptResult{}
			err = r.(error)
		}
	}()

	ctx, err := gpgme.New()
	signature, err := gpgme.NewDataReader(strings.NewReader(r.signature))
	handleErr(err)
	message, err := gpgme.NewDataReader(strings.NewReader(r.message))
	handleErr(err)


	plain, err := gpgme.NewData()
	handleErr(err)

	var signatures []gpgme.Signature
	if (r.signature != ""){
		_, signatures, err = ctx.Verify(signature, message, nil)
		handleErr(err)

		err = ctx.Decrypt(message, plain)
		handleErr(err)
	} else {
		_, signatures, err = ctx.Verify(message, nil, plain)
		handleErr(err)
	}

	for _, signature := range signatures {
		validity := (signature.Summary & gpgme.SigSumGreen) != 0
		result.signatures = append(result.signatures, OpenPgpJsDecryptSignature{keyid: signature.Fingerprint, valid: validity})
	}


	plain.Seek(0, gpgme.SeekSet)

	buf := new(bytes.Buffer)
	io.Copy(buf, plain)
	result.dataString = buf.String()
	return;
}
