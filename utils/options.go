package utils

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (*UtilsStruct) GetOptions() bind.CallOpts {
	block, _ := new(big.Int).SetString("", 10)
	return bind.CallOpts{
		Pending:     false,
		BlockNumber: block,
		Context:     context.Background(),
	}
}
