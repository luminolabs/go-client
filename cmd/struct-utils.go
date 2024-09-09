package cmd

import (
	"crypto/ecdsa"
	"math/big"
	"os"
	"time"

	"lumino/core/types"
	"lumino/path"
	"lumino/pkg/bindings"
	"lumino/utils"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var utilsInterface = utils.UtilsInterface

// This function initializes the utils
func InitializeUtils() {
	utilsInterface = &utils.UtilsStruct{}
	utils.UtilsInterface = &utils.UtilsStruct{}
	utils.FlagSetInterface = &utils.FLagSetStruct{}
	utils.EthClient = &utils.EthClientStruct{}
	utils.ClientInterface = &utils.ClientStruct{}
	utils.Time = &utils.TimeStruct{}
	utils.OS = &utils.OSStruct{}
	utils.PathInterface = &utils.PathStruct{}
	utils.BindInterface = &utils.BindStruct{}
	utils.BlockManagerInterface = &utils.BlockManagerStruct{}
	utils.BindingsInterface = &utils.BindingsStruct{}
	utils.RetryInterface = &utils.RetryStruct{}
}

// This function returns the gas multiplier of root in float32
func (flagSetUtils FlagSetUtils) GetRootFloat32GasMultiplier() (float32, error) {
	return rootCmd.PersistentFlags().GetFloat32("gasmultiplier")
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

// This function returns the provider of root in string
func (flagSetUtils FlagSetUtils) GetRootStringProvider() (string, error) {
	return rootCmd.PersistentFlags().GetString("provider")
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

// This function returns the options
func (u Utils) GetOptions() bind.CallOpts {
	return utilsInterface.GetOptions()
}

// This function returns the default path
func (u Utils) GetDefaultPath() (string, error) {
	return path.PathUtilsInterface.GetDefaultPath()
}

// This function returns the config file path
func (u Utils) GetConfigFilePath() (string, error) {
	return path.PathUtilsInterface.GetConfigFilePath()
}

// This function retrns the block manager
func (u Utils) GetBlockManager(client *ethclient.Client) *bindings.BlockManager {
	return utilsInterface.GetBlockManager(client)
}

// This function assigns the log file
func (u Utils) AssignLogFile(flagSet *pflag.FlagSet) {
	utilsInterface.AssignLogFile(flagSet)
}

// This function checks if the flag is passed
func (u Utils) IsFlagPassed(name string) bool {
	return utilsInterface.IsFlagPassed(name)
}

// This function assigns the password
func (u Utils) AssignPassword(flagSet *pflag.FlagSet) string {
	return utils.AssignPassword(flagSet)
}

// This function prompts the password
func (u Utils) PasswordPrompt() string {
	return utils.PasswordPrompt()
}

// This function prompts the private key
func (u Utils) PrivateKeyPrompt() string {
	return utils.PrivateKeyPrompt()
}

// This function connects to the client
func (u Utils) ConnectToEthClient(provider string) *ethclient.Client {
	log.Debug("Attempting to connect to Ethereum client at: ", provider)
	client, err := ethclient.Dial(provider)
	if err != nil {
		log.Debug("Error in connecting: ", err)
		return nil
	}
	log.Info("Connected to: ", provider)
	return client
}

func (keystoreUtils KeystoreUtils) ImportECDSA(path string, priv *ecdsa.PrivateKey, passphrase string) (accounts.Account, error) {
	ks := keystore.NewKeyStore(path, keystore.StandardScryptN, keystore.StandardScryptP)
	return ks.ImportECDSA(priv, passphrase)
}

// This function is used to write config as
func (v ViperUtils) ViperWriteConfigAs(path string) error {
	return viper.WriteConfigAs(path)
}

// This function is used to convert from Hex to ECDSA
func (c CryptoUtils) HexToECDSA(hexKey string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(hexKey)
}

// This function is used for sleep
func (t TimeUtils) Sleep(duration time.Duration) {
	utils.Time.Sleep(duration)
}

// This function returns the staker Info
func (stateManagerUtils StateManagerUtils) NetworkInfo(client *ethclient.Client, opts *bind.CallOpts) (types.NetworkInfo, error) {

	stateManager := utilsInterface.GetStateManager(client)
	epoch := utils.InvokeFunctionWithTimeout(stateManager, "GetEpoch", opts)
	epochError := utils.CheckIfAnyError(epoch)
	if epochError != nil {
		return types.NetworkInfo{}, epochError
	}
	epochVal := epoch[0].Interface().(uint32)

	state := utils.InvokeFunctionWithTimeout(stateManager, "GetState", opts, uint8(20))
	stateError := utils.CheckIfAnyError(state)
	if stateError != nil {
		return types.NetworkInfo{}, stateError
	}
	stateVal := state[0].Interface().(uint8)

	return types.NetworkInfo{
		EpochNumber: epochVal, State: types.EpochState(stateVal), Timestamp: time.Now()}, nil
}

// This function is used for exiting the code
func (o OSUtils) Exit(code int) {
	os.Exit(code)
}
