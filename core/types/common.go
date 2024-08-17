package types

import (
	"fmt"
	"math/big"
)

// Address represents an Ethereum address
type Address string

// Hash represents a 32-byte Keccak256 hash
type Hash [32]byte

// Amount represents a big integer amount
type Amount struct {
	*big.Int
}

// NewAmount creates a new Amount from a string
func NewAmount(s string) (Amount, error) {
	i, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return Amount{}, fmt.Errorf("invalid amount: %s", s)
	}
	return Amount{i}, nil
}
