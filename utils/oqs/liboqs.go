package oqs

import "io"

type INT_TO_KEYGEN map[int]*KeyGen

const (
	NONE       = 0
	KYBER_1024 = 1
	KYBER_768  = 2
	KYBER_512  = 3
	FALCON_SIG = 4
	DILITHIUM3 = 5
	SCHNORR    = 6
	ED25519    = 7
	RSA        = 8
	ECDSA      = 9
)

var KeyGenRepo INT_TO_KEYGEN

type KeyGen interface {
	Generate() (sk interface{}, pub interface{}, err error)
	GenerateFromSeed(seed []byte) (sk interface{}, pub interface{}, err error)
}

type Key interface {
	Type() int
	Features() ([]string, error)
	Marshal() (k []byte, err error)
	Unmarshal(data []byte, to interface{}) error
	Raw() (rawData interface{}, err error)
	Signature() []byte
	Hex() []byte
	Base64() string
	Store() error
	Load() error
}

type PrivateKey interface {
	Key
	GetPublic() (pub []byte, err error)
	Sign(to io.Writer, data ...[]byte) error
	Decrypt() []byte
}

type PublicKey interface {
	Key
	Derive(priv *key, to io.Writer) error
	Encrypt() []byte
}

type DecentralizedKey interface {
	Key
}

type key struct{}

type keygen struct{}

func init() {

}
