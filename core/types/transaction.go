package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

type TransactionOptions struct {
	Client          *ethclient.Client
	Password        string
	EtherValue      *big.Int
	Amount          *big.Int
	AccountAddress  string
	ChainId         *big.Int
	Config          Configurations
	ContractAddress string
	MethodName      string
	Parameters      []interface{}
	ABI             string
}
