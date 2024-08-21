package utils

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func CheckTransactionReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	ctx := context.Background()
	receipt, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}
