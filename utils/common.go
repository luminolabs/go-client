package utils

import (
	"math/big"

	"lumino/core"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetStateName(stateNumber int64) string {
	var stateName string
	switch stateNumber {
	case 0:
		stateName = "Assign"
	case 1:
		stateName = "Accept"
	case 2:
		stateName = "Confirm"
	default:
		stateName = "Buffer"
	}
	return stateName
}

func ToWei(ether float64) *big.Int {
	ethWei := new(big.Float).Mul(big.NewFloat(ether), big.NewFloat(1e18))
	wei, _ := ethWei.Int(nil)
	return wei
}

func FromWei(wei *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18))
}

func IsValidAddress(address string) bool {
	return common.IsHexAddress(address)
}

func (*UtilsStruct) ConnectToClient(provider string) *ethclient.Client {
	client, err := EthClient.Dial(provider)
	if err != nil {
		log.Fatal("Error in connecting...", err)
	}
	log.Info("Connected to: ", provider)
	return client
}

// func (*UtilsStruct) GetDelayedState(client *ethclient.Client, buffer int32) (int64, error) {
// 	block, err := UtilsInterface.GetLatestBlockWithRetry(client)
// 	if err != nil {
// 		return -1, err
// 	}
// 	stateBuffer, err := UtilsInterface.GetStateBuffer(client)
// 	if err != nil {
// 		return -1, err
// 	}
// 	blockTime := uint64(block.Time)
// 	lowerLimit := (core.StateLength * uint64(buffer)) / 100
// 	upperLimit := core.StateLength - (core.StateLength*uint64(buffer))/100
// 	if blockTime%(core.StateLength) > upperLimit-stateBuffer || blockTime%(core.StateLength) < lowerLimit+stateBuffer {
// 		return -1, nil
// 	}
// 	state := blockTime / core.StateLength
// 	return int64(state) % core.NumberOfStates, nil
// }

func (*UtilsStruct) GetEpoch(client *ethclient.Client) (uint32, error) {
	latestHeader, err := UtilsInterface.GetLatestBlockWithRetry(client)
	if err != nil {
		log.Error("Error in fetching block: ", err)
		return 0, err
	}
	epoch := uint64(latestHeader.Time) / uint64(core.EpochLength)
	return uint32(epoch), nil
}
