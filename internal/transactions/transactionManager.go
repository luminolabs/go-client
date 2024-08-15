package transactionmanager

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Manager struct {
	ethClient *ethclient.Client
}

func NewManager(ethClient *ethclient.Client) *Manager {
	return &Manager{
		ethClient: ethClient,
	}
}

func (m *Manager) SendTransaction(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	err := m.ethClient.SendTransaction(ctx, tx)
	if err != nil {
		return nil, err
	}

	return m.WaitForReceipt(ctx, tx)
}

func (m *Manager) WaitForReceipt(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()

	for {
		receipt, err := m.ethClient.TransactionReceipt(ctx, tx.Hash())
		if err == nil {
			return receipt, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-queryTicker.C:
		}
	}
}

// func (m *Manager) EstimateGas(ctx context.Context, tx *types.Transaction) (uint64, error) {
// 	return m.ethClient.EstimateGas(ctx, ethereum.CallMsg{
// 		From:     tx.From(),
// 		To:       tx.To(),
// 		Gas:      tx.Gas(),
// 		GasPrice: tx.GasPrice(),
// 		Value:    tx.Value(),
// 		Data:     tx.Data(),
// 	})
// }
