package utils

import (
	"math/big"

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
