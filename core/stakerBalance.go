package core

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

func GetLuminoBalanceForStaker(ctx context.Context, client *ethclient.Client, address common.Address) (*big.Int, error) {
	// Get the balance of the address in Wei (smallest unit of Ether)
	balance, err := client.BalanceAt(ctx, address, nil)
	if err != nil {
		return nil, err
	}

	return balance, nil
}
