// Package core implements the fundamental blockchain operations and data structures
// for the Lumino network, including block processing and chain management.
package core

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

// latestBlock stores the most recent block header retrieved from the chain.
var latestBlock *types.Header

// Access is protected by a mutex to ensure thread safety.
// mu provides mutual exclusion for latestBlock access.
var mu = sync.Mutex{}

// GetLatestBlock safely retrieves the current latest block header.
// Uses mutex locking to ensure thread-safe access to the block data.
func GetLatestBlock() *types.Header {
	mu.Lock()
	defer mu.Unlock()
	return latestBlock
}

// SetLatestBlock safely updates the latest block header.
// Uses mutex locking to ensure thread-safe updates to the block data.
func SetLatestBlock(block *types.Header) {
	mu.Lock()
	latestBlock = block
	mu.Unlock()
}

// CalculateLatestBlock continuously updates the latest block information.
// Runs in an infinite loop, fetching new block headers at regular intervals
// defined by BlockNumberInterval. Handles connection errors gracefully by
// continuing to retry on failures.
func CalculateLatestBlock(client *ethclient.Client) {
	for {
		if client != nil {
			latestHeader, err := client.HeaderByNumber(context.Background(), nil)
			if err != nil {
				logrus.Error("CalculateBlockNumber: Error in fetching block: ", err)
				continue
			}
			SetLatestBlock(latestHeader)
		}
		time.Sleep(time.Second * time.Duration(BlockNumberInterval))
	}
}
