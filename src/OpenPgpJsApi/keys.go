package OpenPgpJsApi

import (
	"github.com/proglottis/gpgme"
	"time"
	"strings"
)

type Key struct {
	Revoked         bool `json:"revoked"`
	Expired         bool `json:"expired"`
	Disabled        bool `json:"disabled"`
	Secret          bool `json:"secret"`
	CanEncrypt      bool `json:"canEncrypt"`
	CanSign         bool `json:"canSign"`
	CanCertify      bool `json:"canCertify"`
	CanAuthenticate bool `json:"canAuthenticate"`
	OwnerTrust      gpgme.Validity `json:"ownerTrust"`
	SubKeys         []SubKey `json:"subKeys"`
	UserIDs         []UserID  `json:"userIDs"`
}

func (key Key) getFingerPrint() string {
	if len(key.SubKeys) == 0 {
		return ""
	}
	return key.SubKeys[0].FingerPrint
}

type SubKey struct {
	Revoked     bool `json:"revoked"`
	Expired     bool `json:"expired"`
	Disabled    bool `json:"disabled"`
	Invalid     bool `json:"invalid"`
	Secret      bool `json:"secret"`
	KeyID       string `json:"keyID"`
	FingerPrint string `json:"fingerPrint"`
	Created     time.Time `json:"created"`
	Expires     time.Time `json:"expires"`
	CardNumber  string `json:"cardNumber"`
}

type UserID struct {
	Revoked  bool `json:"revoked"`
	Invalid  bool `json:"invalid"`
	Validity gpgme.Validity `json:"validity"`
	UID      string `json:"UID"`
	Name     string `json:"name"`
	Comment  string `json:"comment"`
	Email    string `json:"email"`
}

type FindKeyResult struct {
	Keys map[string]Key `json:"keys"`
}

func newKey(key *gpgme.Key) (ret Key) {
	ret = Key{
		Revoked:         key.Revoked(),
		Expired:         key.Expired(),
		Disabled:        key.Disabled(),
		Secret:          key.Secret(),
		CanEncrypt:      key.CanEncrypt(),
		CanSign:         key.CanSign(),
		CanCertify:      key.CanCertify(),
		CanAuthenticate: key.CanAuthenticate(),
		OwnerTrust:      key.OwnerTrust(),
	}

	for subKey := key.SubKeys(); subKey != nil; subKey = subKey.Next() {
		ret.SubKeys = append(ret.SubKeys, SubKey{
			Revoked:     subKey.Revoked(),
			Expired:     subKey.Expired(),
			Disabled:    subKey.Disabled(),
			Invalid:     subKey.Invalid(),
			Secret:      subKey.Secret(),
			KeyID:       subKey.KeyID(),
			FingerPrint: subKey.Fingerprint(),
			Created:     subKey.Created(),
			Expires:     subKey.Expires(),
			CardNumber:  subKey.CardNumber(),
		})
	}

	for userID := key.UserIDs(); userID != nil; userID = userID.Next() {
		ret.UserIDs = append(ret.UserIDs, UserID{
			Invalid:  userID.Invalid(),
			Revoked:  userID.Revoked(),
			Validity: userID.Validity(),
			UID:      userID.UID(),
			Name:     userID.Name(),
			Comment:  userID.Comment(),
			Email:    userID.Email(),
		})
	}

	return
}

func (r FindKeyRequest) Execute() (result FindKeyResult, err error) {
	candidates := make(map[string]*gpgme.Key)
	result.Keys = make(map[string]Key)

	for _, searchValue := range []string{""} {
		foundKeys, err := gpgme.FindKeys(searchValue, r.SecretOnly)
		handleErr(err)
		for _, key := range foundKeys {
			candidates[key.SubKeys().Fingerprint()] = key
		}
	}
	for _, origKey := range candidates {
		key := newKey(origKey)
		if r.matches(key) {
			result.Keys[key.getFingerPrint()] = key
		}
	}
	return
}

type FindKeyRequest struct {
	KeyID       string `json:"keyID"`
	FingerPrint string `json:"fingerPrint"`
	UID         string `json:"UID"`
	Name        string `json:"name"`
	Comment     string `json:"comment"`
	Email       string `json:"email"`
	SecretOnly  bool `json:"secretOnly"`
}

func strMatch(haystack, needle string) bool {
	return needle == "" || strings.Contains(haystack, needle)
}

func strMatchExact(haystack, needle string) bool {
	return needle == "" || haystack == needle
}

func (r FindKeyRequest) matches(key Key) bool {
	if r.SecretOnly && !key.Secret {
		return false
	}
	subKeyMatches := r.KeyID == "" && r.FingerPrint == ""
	if !subKeyMatches {
		for _, subKey := range key.SubKeys {
			if strMatch(subKey.KeyID, r.KeyID) && strMatchExact(subKey.FingerPrint, r.FingerPrint) {
				subKeyMatches = true
				break
			}
		}
	}
	if !subKeyMatches {
		return false
	}
	userIDMatches := r.UID == "" && r.Name == "" && r.Comment == "" && r.Email == ""
	if !userIDMatches {
		for _, userID := range key.UserIDs {
			if strMatchExact(userID.UID, r.UID) && strMatch(userID.Name, r.Name) && strMatch(userID.Comment, r.Comment) && strMatch(userID.Email, r.Email) {
				userIDMatches = true
				break
			}
		}
	}
	if !userIDMatches {
		return false
	}
	return true

}
