package OpenPgpJsApi

import (
	"github.com/proglottis/gpgme"
	"io"
	"strings"
	"bytes"
	"fmt"
	"errors"
)

// @see https://openpgpjs.org/openpgpjs/doc/module-openpgp.html#.encrypt

// unused options from the openpgp.js API:
//passwords   //String | Array.<String> 	(optional) array of passwords or a single password to encrypt the Message
//filename    //String 	(optional) a filename for the literal Data packet
type EncryptRequest struct {
	DataString  string `json:"dataString"`      //String | Uint8Array 	text/Data to be encrypted as JavaScript binary string or Uint8Array
	DataBytes   []byte  `json:"dataBytes"`      //String | Uint8Array 	text/Data to be encrypted as JavaScript binary string or Uint8Array
	PublicKeys  []string `json:"publicKeys"`   //Key | Array.<Key> 	(optional) array of keys or single key, used to encrypt the Message
	PrivateKeys []string  `json:"privateKeys"` //Key | Array.<Key> 	(optional) private keys for signing. If omitted Message will not be signed TODO
	Armor       bool `json:"armor"`             //Boolean 	(optional) if the return values should be ascii armored or the Message/Signature objects
	Detached    bool `json:"detached"`          //Boolean 	(optional) if the Signature should be Detached (if true, Signature will be added to returned object) TODO
	Signature   interface{} `json:"signature"`  //Signature 	(optional) a Detached Signature to add to the encrypted Message TODO
}

func (r EncryptRequest) String() string {
	return fmt.Sprintf("DataString: %s, DataBytes: %s, PublicKeys: %s, PrivateKeys: %s, Armor: %b, Detached: %b, Signature: %s",
	r.DataString,
	r.DataBytes,
	r.PublicKeys,
	r.PrivateKeys,
	r.Armor,
	r.Detached,
	r.Signature)
}

type EncryptResult struct {
	Data      string `json:"data"`      // ASCII armored Message if 'Armor' is true
	Message   []byte `json:"message"`   // full Message object if 'Armor' is false
	Signature []byte `json:"signature"` //Detached Signature if 'Detached' is true
}

func (r EncryptRequest) Execute() (result EncryptResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			result = EncryptResult{}
			err = r.(error)
		}
	}()

	recipients := make([]*gpgme.Key, 0)
	for _, pubKey := range r.PublicKeys {
		foundKeys, err := gpgme.FindKeys(pubKey, false)
		handleErr(err)
		for _, foundKey := range foundKeys {
			recipients = append(recipients, foundKey)
		}
	}

	if len(recipients) == 0 {
		handleErr(errors.New("no recipient key found"))
	}

	cipher, err := gpgme.NewData()
	handleErr(err)
	defer cipher.Close()

	var plain *gpgme.Data
	if r.DataString != "" {
		plain, err = gpgme.NewDataReader(strings.NewReader(r.DataString))
	} else {
		plain, err = gpgme.NewDataReader(bytes.NewReader(r.DataBytes))
	}
	handleErr(err)
	defer plain.Close()

	ctx, err := gpgme.New()
	handleErr(err)
	defer ctx.Release()

	ctx.SetArmor(r.Armor)
	err = ctx.Encrypt(recipients, 0, plain, cipher)
	handleErr(err)

	cipher.Seek(0, gpgme.SeekSet)
	buf := new(bytes.Buffer)
	io.Copy(buf, cipher)

	if r.Armor {
		result.Data = buf.String()
	} else {
		result.Message = buf.Bytes()
	}

	return
}
