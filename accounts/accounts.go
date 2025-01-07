// Package accounts implements Ethereum account management functionality including
// account creation, private key handling, and signing operations.
package accounts

import (
	"crypto/ecdsa"
	"errors"
	"lumino/core/types"
	"lumino/logger"
	"lumino/path"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
)

var log = logger.NewLogger()

// CreateAccount generates a new Ethereum account with the provided keystore path and password.
// It ensures the keystore directory exists, creating it if necessary, then generates
// a new account with the specified parameters. Returns the newly created account.
// If directory creation or account generation fails, it logs a fatal error.
func (AccountUtils) CreateAccount(keystorePath string, password string) accounts.Account {
	if _, err := path.OSUtilsInterface.Stat(keystorePath); path.OSUtilsInterface.IsNotExist(err) {
		mkdirErr := path.OSUtilsInterface.Mkdir(keystorePath, 0700)
		if mkdirErr != nil {
			log.Fatal("Error in creating directory: ", mkdirErr)
		}
	}
	newAcc, err := AccountUtilsInterface.NewAccount(keystorePath, password)
	if err != nil {
		log.Fatal("Error in creating account: ", err)
	}
	return newAcc
}

// GetPrivateKeyFromKeystore extracts the private key from a keystore file.
// It reads the encrypted keystore JSON, decrypts it using the provided password,
// and returns the decrypted private key. Returns an error if file reading or
// decryption fails.
func (AccountUtils) GetPrivateKeyFromKeystore(keystorePath string, password string) (*ecdsa.PrivateKey, error) {
	jsonBytes, err := AccountUtilsInterface.ReadFile(keystorePath)
	if err != nil {
		log.Error("Error in reading keystore: ", err)
		return nil, err
	}
	key, err := AccountUtilsInterface.DecryptKey(jsonBytes, password)
	if err != nil {
		log.Error("Error in fetching private key: ", err)
		return nil, err
	}
	return key.PrivateKey, nil
}

// GetPrivateKey retrieves the private key for a given account address.
// It searches through all accounts in the keystore directory to find a matching address,
// then extracts the private key from the corresponding keystore file.
// Returns an error if no matching account is found or key extraction fails.
func (AccountUtils) GetPrivateKey(address string, password string, keystorePath string) (*ecdsa.PrivateKey, error) {
	allAccounts := AccountUtilsInterface.Accounts(keystorePath)
	for _, account := range allAccounts {
		if strings.EqualFold(account.Address.Hex(), address) {
			return AccountUtilsInterface.GetPrivateKeyFromKeystore(account.URL.Path, password)
		}
	}
	return nil, errors.New("no keystore file found")
}

// SignData signs the provided hash using the account's private key.
// It first retrieves the private key for the account, then uses it to generate
// a cryptographic signature of the input hash. Returns the signature as a byte array
// or an error if key retrieval or signing fails.
func (AccountUtils) SignData(hash []byte, account types.Account, defaultPath string) ([]byte, error) {
	privateKey, err := AccountUtilsInterface.GetPrivateKey(account.Address, account.Password, defaultPath)
	if err != nil {
		return nil, err
	}
	return AccountUtilsInterface.Sign(hash, privateKey)
}
