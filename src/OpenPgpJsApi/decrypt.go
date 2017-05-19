package OpenPgpJsApi

import (
	"github.com/proglottis/gpgme"
	"io"
	"strings"
	"bytes"
)

// @see https://openpgpjs.org/openpgpjs/doc/module-openpgp.html#.decrypt

// unused options from the openpgp.js API:
// privateKey  //Key 	(optional) private key with decrypted secret key Data or session key
// sessionKey //Object 	(optional) session key in the form: { Data:Uint8Array, algorithm:String }
// password   string //String 	(optional) single password to decrypt the Message
//

type DecryptRequest struct {
	Message string  `json:"message"`
	//Message 	the Message object with the encrypted Data
	// passed as Armored String
	// TODO: allow for byte[] data
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

	messageReader := strings.NewReader(r.Message)
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
		validity := (signature.Summary & gpgme.SigSumGreen) != 0
		result.Signatures = append(result.Signatures, DecryptSignature{Keyid: signature.Fingerprint, Valid: validity})
	}

	plain.Seek(0, gpgme.SeekSet)
	buf := new(bytes.Buffer)
	io.Copy(buf, plain)
	result.DataString = buf.String()

	return
}
