package keymanager

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Manager struct {
	privateKey *ecdsa.PrivateKey
	keystore   *keystore.KeyStore
}

func NewManager(privateKeyHex string) (*Manager, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	ks := keystore.NewKeyStore("./keystore", keystore.StandardScryptN, keystore.StandardScryptP)

	return &Manager{
		privateKey: privateKey,
		keystore:   ks,
	}, nil
}

func (m *Manager) Sign(data []byte) ([]byte, error) {
	return crypto.Sign(crypto.Keccak256(data), m.privateKey)
}

func (m *Manager) GetAddress() common.Address {
	return crypto.PubkeyToAddress(m.privateKey.PublicKey)
}

func (m *Manager) ImportPrivateKey(privateKeyHex string, passphrase string) (common.Address, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return common.Address{}, err
	}

	account, err := m.keystore.ImportECDSA(privateKey, passphrase)
	if err != nil {
		return common.Address{}, err
	}

	return account.Address, nil
}

func (m *Manager) UnlockAccount(address common.Address, passphrase string) error {
	return m.keystore.Unlock(address, passphrase)
}
