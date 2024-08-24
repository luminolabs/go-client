package core

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"lumino/core/contracts"
	"math/big"
)

func GetLuminoBalance(ctx context.Context, client *ethclient.Client, address common.Address) (*big.Int, error) {
	// Replace this with your actual LUMINO token contract address
	tokenAddress := common.HexToAddress("0x1234567890123456789012345678901234567890")

	// Create a new instance of the token contract
	tokenContract, err := contracts.NewLuminoToken(tokenAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate token contract: %w", err)
	}

	// Call the balanceOf function of the token contract
	balance, err := tokenContract.BalanceOf(&bind.CallOpts{Context: ctx}, address)
	if err != nil {
		return nil, fmt.Errorf("failed to get token balance: %w", err)
	}

	return balance, nil
}
