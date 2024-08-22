package utils

import (
	"lumino/core"
	"lumino/pkg/bindings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func (*UtilsStruct) GetStateManager(client *ethclient.Client) *bindings.StateManager {
	stateManagerContract, err := BindingsInterface.NewStateManager(common.HexToAddress(core.StateManagerAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	return stateManagerContract
}
