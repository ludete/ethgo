package keystore

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/crypto"
)

func RecoverKey(privateKey string) (*ecdsa.PrivateKey, error) {
	if len(privateKey) == 0 {
		return nil, fmt.Errorf("private key is empty")
	}
	privKey, err := crypto.HexToECDSA(privateKey)
	return privKey, err
}

func GetAddrFromKey(key *ecdsa.PrivateKey) common.Address {
	return crypto.PubkeyToAddress(key.PublicKey)
}
