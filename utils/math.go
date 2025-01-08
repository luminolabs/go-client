package utils

import (
	"fmt"
	"math/big"
)

// GetAmountInWei converts token amounts to Wei denomination.
// Multiplies the input amount by 10^18 to get Wei equivalent.
func GetAmountInWei(amount *big.Int) *big.Int {
	amountInWei := big.NewInt(1).Mul(amount, big.NewInt(1e18))
	return amountInWei
}

// ParseBigInt safely converts string to big.Int with base 10.
// Returns error if string cannot be parsed as valid integer.
func ParseBigInt(s string) (*big.Int, error) {
	n := new(big.Int)
	n, ok := n.SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("invalid big integer: %s", s)
	}
	return n, nil
}

// CheckAmountAndBalance verifies if transaction amount exceeds available balance.
// Terminates with fatal error if amount exceeds balance.
// Returns amount in Wei if check passes.
func CheckAmountAndBalance(amountInWei *big.Int, balance *big.Int) *big.Int {
	if amountInWei.Cmp(balance) > 0 {
		log.Fatal("Not enough Lumino token balance")
	}
	return amountInWei
}

// MultiplyFloatAndBigInt performs safe multiplication between float and big.Int.
// Handles potential overflow cases and nil values.
// Essential for gas price calculations and token conversions.
func (*UtilsStruct) MultiplyFloatAndBigInt(bigIntVal *big.Int, floatingVal float64) *big.Int {
	if bigIntVal == nil || floatingVal == 0 {
		return big.NewInt(0)
	}
	value := new(big.Float)
	value.SetFloat64(floatingVal)
	conversionInt := new(big.Float)
	conversionInt.SetInt(bigIntVal)
	value.Mul(value, conversionInt)
	result := new(big.Int)
	value.Int(result)
	return result
}
