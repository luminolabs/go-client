package cmd

import (
	"lumino/utils"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/pflag"
)

var utilsInterface = utils.UtilsInterface

// This function initializes the utils
func InitializeUtils() {
	utilsInterface = &utils.UtilsStruct{}
	utils.FlagSetInterface = &utils.FLagSetStruct{}
}

// GetRootInt32Buffer retrieves the buffer value from root command flags.
// It returns the buffer as an int32 and an error if retrieval fails.
func (FlagSetUtils FlagSetUtils) GetRootInt32Buffer() (int32, error) {
	return rootCmd.PersistentFlags().GetInt32("buffer")
}

// This function returns the wait of root in Int32
func (FlagSetUtils FlagSetUtils) GetRootInt32Wait() (int32, error) {
	return rootCmd.PersistentFlags().GetInt32("wait")
}

// This function returns the gas price of root in Int32
func (FlagSetUtils FlagSetUtils) GetRootInt32GasPrice() (int32, error) {
	return rootCmd.PersistentFlags().GetInt32("gasprice")
}

// This function returns the log level of root in string
func (FlagSetUtils FlagSetUtils) GetRootStringLogLevel() (string, error) {
	return rootCmd.PersistentFlags().GetString("logLevel")
}

// This function returns the gas limit of root in Float32
func (FlagSetUtils FlagSetUtils) GetRootFloat32GasLimit() (float32, error) {
	return rootCmd.PersistentFlags().GetFloat32("gasLimit")
}

// This function returns the gas limit of root in Float32
func (FlagSetUtils FlagSetUtils) GetRootInt64RPCTimeout() (int64, error) {
	return rootCmd.PersistentFlags().GetInt64("rpcTimeout")
}

// This function returns the provider in string
func (FlagSetUtils FlagSetUtils) GetStringProvider(flagSet *pflag.FlagSet) (string, error) {
	return flagSet.GetString("provider")
}

// This function returns gas multiplier in float 32
func (FlagSetUtils FlagSetUtils) GetFloat32GasMultiplier(flagSet *pflag.FlagSet) (float32, error) {
	return flagSet.GetFloat32("gasmultiplier")
}

// This function returns Buffer in Int32
func (FlagSetUtils FlagSetUtils) GetInt32Buffer(flagSet *pflag.FlagSet) (int32, error) {
	return flagSet.GetInt32("buffer")
}

// This function returns Wait in Int32
func (FlagSetUtils FlagSetUtils) GetInt32Wait(flagSet *pflag.FlagSet) (int32, error) {
	return flagSet.GetInt32("wait")
}

// This function returns GasPrice in Int32
func (FlagSetUtils FlagSetUtils) GetInt32GasPrice(flagSet *pflag.FlagSet) (int32, error) {
	return flagSet.GetInt32("gasprice")
}

// This function returns Log Level in string
func (FlagSetUtils FlagSetUtils) GetStringLogLevel(flagSet *pflag.FlagSet) (string, error) {
	return flagSet.GetString("logLevel")
}

func (FlagSetUtils FlagSetUtils) GetInt64RPCTimeout(flagSet *pflag.FlagSet) (int64, error) {
	return flagSet.GetInt64("rpcTimeout")
}

// This function returns Gas Limit in Float32
func (FlagSetUtils FlagSetUtils) GetFloat32GasLimit(flagSet *pflag.FlagSet) (float32, error) {
	return flagSet.GetFloat32("gasLimit")
}

// GetDelayedState calculates the delayed state based on the current block and buffer.
// It returns the delayed state as an int64 and an error if calculation fails.
func (u Utils) GetDelayedState(client *ethclient.Client, buffer int32) (int64, error) {
	return utilsInterface.GetDelayedState(client, buffer)
}

// This function returns the amount in wei
func (u Utils) GetAmountInWei(amount *big.Int) *big.Int {
	return utils.GetAmountInWei(amount)
}

// This function returns the epoch
func (u Utils) GetEpoch(client *ethclient.Client) (uint32, error) {
	return utilsInterface.GetEpoch(client)
}

// This function connects to the client
func (u Utils) ConnectToEthClient(provider string) *ethclient.Client {
	returnedValues := utils.InvokeFunctionWithTimeout(utilsInterface, "ConnectToEthClient", provider)
	returnedError := utils.CheckIfAnyError(returnedValues)
	if returnedError != nil {
		return nil
	}
	return returnedValues[0].Interface().(*ethclient.Client)
}
