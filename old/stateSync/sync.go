package statesync

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type State int

const (
	StateAssign State = iota
	StateAccept
	StateConfirm
	StateBuffer
)

type Manager struct {
	ethClient    *ethclient.Client
	epochLength  uint64
	numStates    uint8
	startTime    time.Time
	currentEpoch uint32
	currentState State
	mu           sync.RWMutex
}

func NewManager(ethClient *ethclient.Client, epochLength uint64, numStates uint8) *Manager {
	return &Manager{
		ethClient:   ethClient,
		epochLength: epochLength,
		numStates:   numStates,
		startTime:   time.Now(),
	}
}

func (m *Manager) Start(ctx context.Context) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := m.sync(); err != nil {
				log.WithError(err).Error("Failed to sync state")
			}
		}
	}
}

func (m *Manager) sync() error {
	header, err := m.ethClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to get latest block header: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.currentEpoch = uint32(header.Time / m.epochLength)
	timeInEpoch := header.Time % m.epochLength
	stateLength := m.epochLength / uint64(m.numStates)
	m.currentState = State(timeInEpoch / stateLength)

	log.WithFields(logrus.Fields{
		"epoch": m.currentEpoch,
		"state": m.currentState,
	}).Debug("State synced")

	return nil
}

func (m *Manager) GetCurrentEpoch() uint32 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentEpoch
}

func (m *Manager) GetCurrentState() State {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentState
}
