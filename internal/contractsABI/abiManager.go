package contractsABI

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Manager struct {
	abis map[string]*abi.ABI
}

func NewManager() *Manager {
	return &Manager{
		abis: make(map[string]*abi.ABI),
	}
}

func (m *Manager) LoadABI(contractName, abiJSON string) error {
	parsedABI, err := abi.JSON([]byte(abiJSON))
	if err != nil {
		return err
	}
	m.abis[contractName] = &parsedABI
	return nil
}

func (m *Manager) GetABI(contractName string) (*abi.ABI, bool) {
	abi, ok := m.abis[contractName]
	return abi, ok
}
