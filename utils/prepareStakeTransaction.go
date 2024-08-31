package utils

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"lumino/logger" // Import the logger package
)

// PrepareStakeTransaction prepares the transaction options for staking
func PrepareStakeTransaction(ctx context.Context, client *ethclient.Client, from common.Address, amount *big.Int, password string) (*bind.TransactOpts, error) {
	// Step 1: Get the keystore directory
	keystoreDir := filepath.Join(os.Getenv("HOME"), ".ethereum", "keystore")
	keyStore := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)

	// Step 2: Find the account in the keystore
	account := accounts.Account{Address: from}
	found := false
	for _, acc := range keyStore.Accounts() {
		if acc.Address == account.Address {
			account = acc
			found = true
			break
		}
	}

	if !found {
		err := fmt.Errorf("account with address %s not found in keystore", from.Hex())
		logger.Error("Account with address not found in keystore: ", err)
		return nil, err
	}

	// Step 3: Unlock the account using the password
	err := keyStore.Unlock(account, password)
	if err != nil {
		logger.Error("Failed to unlock account: ", err)
		return nil, err
	}

	// Step 4: Retrieve the chain ID
	chainID, err := client.ChainID(ctx)
	if err != nil {
		logger.Error("Failed to retrieve chain ID: ", err)
		return nil, err
	}

	// Step 5: Create a TransactOpts using the keystore's Signer
	txOpts, err := bind.NewKeyStoreTransactorWithChainID(keyStore, account, chainID)
	if err != nil {
		logger.Error("Failed to create transaction signer: ", err)
		return nil, err
	}

	// Step 6: Set the sender address and value (amount)
	txOpts.From = from
	txOpts.Value = amount // amount of ETH or token to send
	txOpts.Context = ctx

	// Set the nonce manually
	nonce, err := client.PendingNonceAt(ctx, from)
	if err != nil {
		logger.Error("Failed to get nonce: ", err)
		return nil, err
	}
	txOpts.Nonce = big.NewInt(int64(nonce)) // Set nonce

	// Step 7: Estimate the gas limit (commented out but available if needed)
	// gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
	// 	From:     from,
	// 	To:       nil, // set to the contract address if needed
	// 	GasPrice: nil, // set to a specific gas price if needed
	// 	Value:    amount,
	// 	Data:     nil, // set if you are calling a function on a contract
	// })
	// if err != nil {
	// 	logger.Error("Failed to estimate gas limit: ", err)
	// 	return nil, err
	// }
	// txOpts.GasLimit = gasLimit

	// Set the gas limit to a fixed value
	txOpts.GasLimit = 300000

	// Step 8: Get the suggested gas price
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		logger.Error("Failed to suggest gas price: ", err)
		return nil, err
	}
	txOpts.GasPrice = gasPrice

	return txOpts, nil
}
