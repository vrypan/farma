package utils

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/vrypan/farma/fctools"
)

// Implement FIP-208 https://github.com/farcasterxyz/protocol/discussions/208
//

type JfsVerification struct {
	UserFid  int
	AppFid   int
	AppKey   string
	Verified bool
}

func JfsVerify(hub *fctools.FarcasterHub, Header, Payload, Signature string) (JfsVerification, error) {
	verificationInfo, err := JfsVerifySig(Header, Payload, Signature)

	if err != nil {
		return JfsVerification{}, err
	}

	verificationInfo.AppFid = int(fctools.AppIdFromFidSigner(hub, uint64(verificationInfo.UserFid), common.FromHex(verificationInfo.AppKey)))
	if verificationInfo.AppFid == 0 {
		return verificationInfo, fmt.Errorf("Unable to verify signer")
	}

	return verificationInfo, nil
}

func JfsVerifySig(Header, Payload, Signature string) (JfsVerification, error) {
	verificationInfo := JfsVerification{0, 0, "", false}

	signatureBytes, _ := base64.RawURLEncoding.DecodeString(Signature)
	headerBytes, _ := base64.RawURLEncoding.DecodeString(Header)

	headerData := make(map[string]interface{})
	if err := json.Unmarshal(headerBytes, &headerData); err != nil {
		return verificationInfo, err
	}

	verificationInfo.UserFid = int(headerData["fid"].(float64))
	verificationInfo.AppKey = headerData["key"].(string)

	signedData := []byte(Header + "." + Payload)
	publicKey := common.FromHex(verificationInfo.AppKey)

	if isValidSig := ed25519.Verify(publicKey, signedData, signatureBytes); !isValidSig {
		return verificationInfo, nil
	}

	verificationInfo.Verified = true
	return verificationInfo, nil
}
