package models

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	db "github.com/vrypan/farma/localdb"
)

// Genaretes a new keypair. k is updated with the public key and frameId.
// It returns the private key and error (nil = no error).
func (k *PubKey) GenerateKey(frameId string) ([]byte, error) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("Error generating Frame key pair: %v\n", err)
	}
	k.FrameId = frameId
	k.Key = pubKey
	return privKey, nil
}

func (k *PubKey) Prefix() string {
	return "f:pk:"
}
func (k *PubKey) DbKey() string {
	return fmt.Sprintf("%s%s:%s", k.Prefix(), k.FrameId, base64.StdEncoding.EncodeToString(k.Key))
}
func (k *PubKey) DbKeyBytes() []byte {
	return []byte(k.DbKey())
}

// Saves a PubKey to the database.
// There is no payload saved, just the dbkey, because
// we are only interested in ckecking if it exists using InDb().
func (k *PubKey) Save() error {
	err := db.Set(k.DbKeyBytes(), []byte(""))
	return err
}
func (k *PubKey) InDb() bool {
	f := NewFrame().FromId(k.FrameId)
	if f != nil && bytes.Equal(f.PublicKey.Key, k.Key) {
		return true
	}
	return false
}
func (k *PubKey) Delete() error {
	return db.Delete(k.DbKeyBytes())
}
func (k *PubKey) Decode(header string) error {
	var err error
	parts := strings.Split(header, ":")
	if len(parts) != 2 {
		return fmt.Errorf("Invalid X-Public-Key format")
	}
	k.FrameId = parts[0]
	k.Key, err = base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		k = nil
		return fmt.Errorf("Invalid Key in X-Public-Key, should be a base64 encoded string")
	}
	return nil
}
func (k *PubKey) Encode() string {
	return fmt.Sprintf("%s:%s", k.FrameId, base64.StdEncoding.EncodeToString(k.Key))
}

func (k *PubKey) FromPrivateKey(frameId string, privKey []byte) *PubKey {
	pubKey := ed25519.PrivateKey(privKey).Public().(ed25519.PublicKey)
	k.FrameId = frameId
	k.Key = pubKey
	return k
}
