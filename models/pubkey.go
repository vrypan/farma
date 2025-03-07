package models

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	db "github.com/vrypan/farma/localdb"
)

type PublicKey struct {
	Key   []byte
	Frame uint64
}

func NewPublicKey(pk []byte, frameId uint64) PublicKey {
	return PublicKey{
		Key:   pk,
		Frame: frameId,
	}
}

func (k PublicKey) Prefix() string {
	return "f:pk:"
}

func (k PublicKey) String() string {
	b64 := base64.StdEncoding.EncodeToString(k.Key)
	return fmt.Sprintf("%s%s:%d", k.Prefix(), b64, k.Frame)
}
func (k PublicKey) Bytes() []byte {
	return []byte(k.String())
}
func (k *PublicKey) DecodeBytes(b []byte) *PublicKey {
	return k.DecodeString(string(b))
}
func (k *PublicKey) DecodeString(s string) *PublicKey {
	s = strings.TrimPrefix(s, k.Prefix())
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return nil
	}
	pk, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return nil
	}
	frameId, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return nil
	}
	k.Key = pk
	k.Frame = frameId
	return k
}

func (k PublicKey) Set(frameKey []byte) error {
	err := db.Set(k.Bytes(), frameKey)
	return err
}

func (k PublicKey) Get() ([]byte, error) {
	return db.Get(k.Bytes())
}

func (k PublicKey) Delete() error {
	return db.Delete(k.Bytes())
}
