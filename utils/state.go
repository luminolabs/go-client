package utils

import (
	"lumino/pkg/bindings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

func (*UtilsStruct) GetStateManagerWithOpts(client *ethclient.Client) (*bindings.StateManager, bind.CallOpts) {
	return UtilsInterface.GetStateManager(client), UtilsInterface.GetOptions()
}
