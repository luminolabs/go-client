// Package accounts provides the core account management functionality
// including account creation, key management, and signing operations
package accounts

import (
	"crypto/ecdsa"
	"lumino/core/types"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

//go:generate mockery --name AccountInterface --output ./mocks/ --case=underscore

// AccountInterface defines the contract for account management operations.
// Implementations must provide methods for account creation, key management,
// and cryptographic signing functions.
var AccountUtilsInterface AccountInterface

type AccountInterface interface {
	CreateAccount(path string, password string) accounts.Account
	GetPrivateKeyFromKeystore(keystorePath string, password string) (*ecdsa.PrivateKey, error)
	GetPrivateKey(address string, password string, keystorePath string) (*ecdsa.PrivateKey, error)
	SignData(hash []byte, account types.Account, defaultPath string) ([]byte, error)
	Accounts(path string) []accounts.Account
	NewAccount(path string, passphrase string) (accounts.Account, error)
	DecryptKey(jsonBytes []byte, password string) (*keystore.Key, error)
	Sign(digestHash []byte, prv *ecdsa.PrivateKey) ([]byte, error)
	ReadFile(filename string) ([]byte, error)
}

// Accounts returns all Ethereum accounts found in the specified keystore directory.
// Creates a new keystore with standard scrypt parameters and returns all accounts.
type AccountUtils struct{}

// Accounts returns all Ethereum accounts found in the specified keystore directory.
// Creates a new keystore with standard scrypt parameters and returns all accounts.
func (accountUtils AccountUtils) Accounts(path string) []accounts.Account {
	ks := keystore.NewKeyStore(path, keystore.StandardScryptN, keystore.StandardScryptP)
	return ks.Accounts()
}

// NewAccount creates a new Ethereum account in the specified keystore directory.
// Uses standard scrypt parameters for key derivation and returns the new account
// along with any error that occurred during creation.
func (accountUtils AccountUtils) NewAccount(path string, passphrase string) (accounts.Account, error) {
	ks := keystore.NewKeyStore(path, keystore.StandardScryptN, keystore.StandardScryptP)
	accounts.NewManager(&accounts.Config{InsecureUnlockAllowed: false}, ks)
	return ks.NewAccount(passphrase)
}

// DecryptKey decrypts an encrypted keystore JSON file using the provided password.
// Returns the decrypted keystore key or an error if decryption fails.
func (accountUtils AccountUtils) DecryptKey(jsonBytes []byte, password string) (*keystore.Key, error) {
	return keystore.DecryptKey(jsonBytes, password)
}

// Sign generates a cryptographic signature for the provided digest hash using the
// specified private key. Returns the signature as a byte array or an error if
// signing fails.
func (accountUtils AccountUtils) Sign(digestHash []byte, prv *ecdsa.PrivateKey) (sig []byte, err error) {
	return crypto.Sign(digestHash, prv)
}

// ReadFile reads and returns the contents of a file at the specified path.
// Returns the file contents as a byte array or an error if reading fails.
func (accountUtils AccountUtils) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}
