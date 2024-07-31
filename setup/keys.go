package setup

import (
	"encoding/json"
	"os"

	"github.com/genovatix/vectorlib/utils"
	"github.com/genovatix/vectorlib/utils/oqs"
)

func SetRootKeys() {
	k, s, err := oqs.Generate()
	if err != nil {
		panic(utils.RevertNoParams(err.Error()))
	}
	WriteRootKeys(k, s)
}

type RootKeysEncrypted struct {
	EncryptedKey1 []byte
	EncryptedKey2 []byte
}

func WriteRootKeys(k, s []byte) {
	rootEncryptionKeyFile := os.Getenv("root_encryption_key")
	e1, err := utils.Encrypt(k, []byte(rootEncryptionKeyFile)[:32])
	if err != nil {
		panic(utils.RevertNoParams(err.Error()))
	}
	e2, err := utils.Encrypt(s, []byte(rootEncryptionKeyFile)[:32])
	if err != nil {
		panic(utils.RevertNoParams(err.Error()))
	}
	rke := &RootKeysEncrypted{EncryptedKey1: []byte(e1), EncryptedKey2: []byte(e2)}
	rkej, _ := json.Marshal(rke)
	os.WriteFile(".rootkeys", rkej, 0600)
}

func init() {
	SetRootKeys()
}
