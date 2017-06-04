package OpenPgpJsApi

import (
	"github.com/proglottis/gpgme"
	"strings"
)

type FindKeyRequest struct {
	KeyID       string `json:"keyID"`
	FingerPrint string `json:"fingerPrint"`
	UID         string `json:"UID"`
	Name        string `json:"name"`
	Comment     string `json:"comment"`
	Email       string `json:"email"`
	SecretOnly  bool `json:"secretOnly"`
}

type FindKeyResult struct {
	Keys map[string]Key `json:"keys"`
}

func (r FindKeyRequest) Execute() (result FindKeyResult, err error) {
	candidates := make(map[string]*gpgme.Key)
	result.Keys = make(map[string]Key)

	for _, searchValue := range []string{r.KeyID, r.FingerPrint, r.UID, r.Name, r.Comment, r.Email} {
		if searchValue == "" {
			continue
		}
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
