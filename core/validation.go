package core

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"os"
	"path/filepath"
)

func ValidatePassword(address common.Address, password string) error {
	// TODO (1): Implement the ValidatePassword function
	// Get the keystore directory
	keystoreDir := filepath.Join(os.Getenv("HOME"), ".ethereum", "keystore")

	// Create a new keystore instance
	ks := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)

	// Find the account
	account := accounts.Account{Address: address}

	// Try to unlock the account with the provided password
	err := ks.Unlock(account, password)
	if err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	return nil
}
