package eventlistener

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Listener struct {
	ethClient *ethclient.Client
}

func NewListener(ethClient *ethclient.Client) *Listener {
	return &Listener{
		ethClient: ethClient,
	}
}

func (l *Listener) SubscribeEvents(ctx context.Context, contractAddress common.Address, eventSignature common.Hash) (chan types.Log, ethereum.Subscription, error) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{eventSignature}},
	}

	logs := make(chan types.Log)
	sub, err := l.ethClient.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to subscribe to event logs: %w", err)
	}

	return logs, sub, nil
}

func (l *Listener) ProcessEvent(event types.Log) error {
	// Implement event processing logic here
	log.WithFields(logrus.Fields{
		"blockNumber": event.BlockNumber,
		"txHash":      event.TxHash.Hex(),
	}).Info("Received event")

	return nil
}
