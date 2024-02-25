package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

func HashObject(object interface{}) (string, error) {
	hasher := sha256.New()
	serializedObject, err := json.Marshal(object)
	if err != nil {
		return "", err
	}
	hasher.Write(serializedObject)
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
