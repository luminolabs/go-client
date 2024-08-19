package keymanager

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Manager struct {
	privateKey *ecdsa.PrivateKey
}

func NewManager(privateKeyHex string) (*Manager, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}
	return &Manager{privateKey: privateKey}, nil
}

func (m *Manager) Sign(data []byte) ([]byte, error) {
	return crypto.Sign(crypto.Keccak256(data), m.privateKey)
}

func (m *Manager) GetAddress() string {
	publicKey := m.privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return ""
	}
	return crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
}
