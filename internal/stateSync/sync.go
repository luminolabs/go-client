package sync

import (
	"sync"
	"time"

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
	epochLength uint64
	numStates   uint8
	startTime   time.Time
	mu          sync.RWMutex
}

func NewManager(epochLength uint64, numStates uint8) *Manager {
	return &Manager{
		epochLength: epochLength,
		numStates:   numStates,
		startTime:   time.Now(),
	}
}

func (m *Manager) CurrentEpoch() uint32 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return uint32(time.Since(m.startTime).Seconds() / float64(m.epochLength))
}

func (m *Manager) CurrentState() State {
	m.mu.RLock()
	defer m.mu.RUnlock()
	timeInEpoch := uint64(time.Since(m.startTime).Seconds()) % m.epochLength
	stateLength := m.epochLength / uint64(m.numStates)
	return State(timeInEpoch / stateLength)
}
