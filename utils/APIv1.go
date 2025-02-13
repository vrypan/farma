package utils

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

var endpoints = map[string]struct{}{
	"notification/send":   {},
	"notification/status": {},
}

type Api struct {
	PubKeys map[string]int // strings are hex, 0x prefixed
}

func NewApi() *Api {
	return &Api{
		PubKeys: make(map[string]int),
	}
}

func (a *Api) AddKey(k string) {
	if !strings.HasPrefix(k, "0x") {
		k = "0x" + k
	}
	a.PubKeys[k] = 1
}

func (a *Api) IsValid(endpoint string, payload string) error {
	endpoint = strings.TrimSuffix(endpoint, "/")

	if _, ok := endpoints[endpoint]; !ok {
		return fmt.Errorf("Invalid endpoint: %s", endpoint)
	}

	var jsonBody, jsonHeader map[string]interface{}

	if err := json.Unmarshal([]byte(payload), &jsonBody); err != nil {
		return err
	}

	header, _ := jsonBody["header"].(string)
	headerBytes, err := base64.RawURLEncoding.DecodeString(header)

	if err != nil {
		return err
	}

	if err := json.Unmarshal(headerBytes, &jsonHeader); err != nil {
		return err
	}

	requestKey, _ := jsonHeader["key"].(string)

	if !strings.HasPrefix(requestKey, "0x") {
		requestKey = "0x" + requestKey
	}

	if _, ok := a.PubKeys[requestKey]; !ok {
		return fmt.Errorf("Invalid key: %s", requestKey)
	}

	signedData, _ := jsonBody["payload"].(string)
	signedDataBytes := []byte(header + "." + signedData)

	signature, _ := jsonBody["signature"].(string)
	signatureBytes, _ := base64.RawURLEncoding.DecodeString(signature)

	if !ed25519.Verify(common.FromHex(requestKey), signedDataBytes, signatureBytes) {
		return fmt.Errorf("Invalid signature")
	}
	return nil
}
