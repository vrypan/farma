package fctools

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mitchellh/mapstructure"
)

type SignedKeyRequestMetadata struct {
	RequestFid    *big.Int
	RequestSigner common.Address
	Signature     []byte
	Deadline      *big.Int
}

func (skr *SignedKeyRequestMetadata) FromPayload(data []byte) error {

	arguments, _ := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{
			Name: "requestFid",
			Type: "uint256",
		},
		{
			Name: "requestSigner",
			Type: "address",
		},
		{
			Name: "signature",
			Type: "bytes",
		},
		{
			Name: "deadline",
			Type: "uint256",
		},
	})
	argumentsTuple := abi.Arguments{
		{
			Type: arguments,
		},
	}

	decoded, err := argumentsTuple.Unpack(data)
	if err != nil {
		return fmt.Errorf("Error unpacking data: %w", err)
	}
	if err := mapstructure.Decode(decoded[0], skr); err != nil {
		return fmt.Errorf("Error converting decode[0] to SignedKeyRequestMetadata. %w\n", err)
	}
	return nil
}

func AppIdFromFidSigner(hub *FarcasterHub, fid uint64, signer []byte) uint64 {
	evt, err := hub.GetSigner(fid, signer)
	if err != nil {
		panic(err) // For prod: Just return s
	}
	skr := SignedKeyRequestMetadata{}
	if skr.FromPayload(evt.GetSignerEventBody().Metadata) != nil {
		return 0
	}
	return skr.RequestFid.Uint64()
}
