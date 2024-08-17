package utils

import (
	"github.com/ethereum/go-ethereum/crypto"
)

func VerifySignature(publicKeyBytes, signature, message []byte) bool {
	return crypto.VerifySignature(publicKeyBytes, message, signature[:64])
}
