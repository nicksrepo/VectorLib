package oqs

import (
	"encoding/base64"
	"encoding/hex"
	"io"

	oqs2 "github.com/open-quantum-safe/liboqs-go/oqs"
)

type PQCPrivateKey PrivateKey
type PQCPublicKey PublicKey

var poqs oqs2.KeyEncapsulation

type Kyber1024_Key []byte

// Base64 implements PrivateKey.
func (kyb *Kyber1024_Key) Base64() string {
	r, _ := kyb.Raw()
	kstr := base64.StdEncoding.EncodeToString(r.([]byte))
	return kstr
}

// Decrypt implements PrivateKey.
func (kyb *Kyber1024_Key) Decrypt() []byte {
	panic("unimplemented")
}

// Features implements PrivateKey.
func (kyb *Kyber1024_Key) Features() ([]string, error) {
	panic("unimplemented")
}

// Hex implements PrivateKey.
func (kyb *Kyber1024_Key) Hex() []byte {
	return []byte(hex.EncodeToString([]byte(kyb.Base64())))
}

// Load implements PrivateKey.
func (kyb *Kyber1024_Key) Load() error {
	panic("unimplemented")
}

// Marshal implements PrivateKey.
func (kyb *Kyber1024_Key) Marshal() (k []byte, err error) {
	panic("unimplemented")
}

// Raw implements PrivateKey.
func (kyb *Kyber1024_Key) Raw() (rawData interface{}, err error) {
	panic("unimplemented")
}

// Sign implements PrivateKey.
func (kyb *Kyber1024_Key) Sign(to io.Writer, data ...[]byte) error {
	panic("unimplemented")
}

// Signature implements PrivateKey.
func (kyb *Kyber1024_Key) Signature() []byte {
	panic("unimplemented")
}

// Store implements PrivateKey.
func (kyb *Kyber1024_Key) Store() error {
	panic("unimplemented")
}

// Type implements PrivateKey.
func (kyb *Kyber1024_Key) Type() int {
	panic("unimplemented")
}

// Unmarshal implements PrivateKey.
func (kyb *Kyber1024_Key) Unmarshal(data []byte, to interface{}) error {
	panic("unimplemented")
}

func Generate() ([]byte, []byte, error) {
	poqs = oqs2.KeyEncapsulation{}
	poqs.Init("Kyber1024", nil)
	defer poqs.Clean()
	pub, err := poqs.GenerateKeyPair()
	if err != nil {
		return nil, nil, err
	}
	sk := poqs.ExportSecretKey()
	return pub, sk, nil
}

func NewPQCPrivateKey(key []byte) PrivateKey {
	return &Kyber1024_Key{}
}

func (kyb Kyber1024_Key) GetPublic() (pub []byte, err error) {
	k := oqs2.KeyEncapsulation{}
	k.Init("Kyber1024", kyb)
	pub, err = k.GenerateKeyPair()
	return
}
