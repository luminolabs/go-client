package utils

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// PrepareStakeTransaction prepares the transaction options for staking
func PrepareStakeTransaction(ctx context.Context, client *ethclient.Client, from common.Address, amount *big.Int, password string) (*bind.TransactOpts, error) {
	// Step 1: Retrieve the account's private key from the keystore
	keyStore := keystore.NewKeyStore("/path/to/keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	account := accounts.Account{Address: from}

	// Step 2: Unlock the account using the password
	err := keyStore.Unlock(account, password)
	if err != nil {
		return nil, fmt.Errorf("failed to unlock account: %w", err)
	}

	// Step 3: Retrieve the private key
	key, err := keyStore.Export(account, password, password)
	if err != nil {
		return nil, fmt.Errorf("failed to export account: %w", err)
	}

	// Step 4: Convert the key to an ecdsa private key
	privateKey, err := crypto.ToECDSA(key)
	if err != nil {
		return nil, fmt.Errorf("failed to convert key to ECDSA: %w", err)
	}

	// Step 5: Retrieve the chain ID
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve chain ID: %w", err)
	}

	// Step 6: Create the TransactOpts
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction signer: %w", err)
	}

	// Step 7: Set the sender address and value (amount)
	auth.From = from
	auth.Value = amount // amount of ETH or token to send
	auth.Context = ctx

	// Step 8: Estimate the gas limit
	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
		From:     from,
		To:       nil, // set to the contract address if needed
		GasPrice: nil, // set to a specific gas price if needed
		Value:    amount,
		Data:     nil, // set if you are calling a function on a contract
	})
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas limit: %w", err)
	}
	auth.GasLimit = gasLimit

	// Step 9: Get the suggested gas price
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price: %w", err)
	}
	auth.GasPrice = gasPrice

	return auth, nil
}
