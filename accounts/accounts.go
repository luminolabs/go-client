// Package account provides all account related functions
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

// This function takes path and password as input and returns new account
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

// This function takes and path of keystore and password as input and returns private key of account
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

// This function takes address of account, password and keystore path as input and returns private key of account
func (AccountUtils) GetPrivateKey(address string, password string, keystorePath string) (*ecdsa.PrivateKey, error) {
	allAccounts := AccountUtilsInterface.Accounts(keystorePath)
	for _, account := range allAccounts {
		if strings.EqualFold(account.Address.Hex(), address) {
			return AccountUtilsInterface.GetPrivateKeyFromKeystore(account.URL.Path, password)
		}
	}
	return nil, errors.New("no keystore file found")
}

// This function takes hash, account and path as input and returns the signed data as array of byte
func (AccountUtils) SignData(hash []byte, account types.Account, defaultPath string) ([]byte, error) {
	privateKey, err := AccountUtilsInterface.GetPrivateKey(account.Address, account.Password, defaultPath)
	if err != nil {
		return nil, err
	}
	return AccountUtilsInterface.Sign(hash, privateKey)
}
