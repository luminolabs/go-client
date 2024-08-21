package utils

// import (
// 	"math/big"

// 	"github.com/ethereum/go-ethereum/common"
// )

// func ToWei(ether float64) *big.Int {
// 	ethWei := new(big.Float).Mul(big.NewFloat(ether), big.NewFloat(1e18))
// 	wei, _ := ethWei.Int(nil)
// 	return wei
// }

// func FromWei(wei *big.Int) *big.Float {
// 	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18))
// }

// func IsValidAddress(address string) bool {
// 	return common.IsHexAddress(address)
// }
