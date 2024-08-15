package contractsabi

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

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

func (m *Manager) LoadABIFromFile(contractName, filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read ABI file: %w", err)
	}

	parsedABI, err := abi.JSON(data)
	if err != nil {
		return fmt.Errorf("failed to parse ABI JSON: %w", err)
	}

	m.abis[contractName] = &parsedABI
	log.WithField("contract", contractName).Info("ABI loaded successfully")
	return nil
}

func (m *Manager) LoadABIsFromDirectory(dirPath string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			contractName := filepath.Base(file.Name()[:len(file.Name())-5])
			err := m.LoadABIFromFile(contractName, filepath.Join(dirPath, file.Name()))
			if err != nil {
				return fmt.Errorf("failed to load ABI for %s: %w", contractName, err)
			}
		}
	}

	return nil
}

func (m *Manager) GetABI(contractName string) (*abi.ABI, error) {
	abi, ok := m.abis[contractName]
	if !ok {
		return nil, fmt.Errorf("ABI for contract %s not found", contractName)
	}
	return abi, nil
}

func (m *Manager) EncodeFunction(contractName, functionName string, args ...interface{}) ([]byte, error) {
	abi, err := m.GetABI(contractName)
	if err != nil {
		return nil, err
	}

	return abi.Pack(functionName, args...)
}

func (m *Manager) DecodeLog(contractName string, log abi.Log) (map[string]interface{}, error) {
	abi, err := m.GetABI(contractName)
	if err != nil {
		return nil, err
	}

	event, err := abi.EventByID(log.Topics[0])
	if err != nil {
		return nil, fmt.Errorf("event not found: %w", err)
	}

	decoded, err := event.Inputs.Unpack(log.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack log data: %w", err)
	}

	return decoded, nil
}
