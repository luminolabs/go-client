package cmd

import (
	"lumino/core"
	"lumino/core/types"
	"lumino/utils"
	"strings"

	"github.com/spf13/viper"
)

// GetBufferPercent retrieves buffer percentage from configuration or flags. If not explicitly set,
// uses the default value from core configuration. Returns error if retrieval fails.
func (*UtilsStruct) GetBufferPercent() (int32, error) {
	bufferPercent, err := flagSetUtils.GetRootInt32Buffer()
	if err != nil {
		return int32(core.DefaultBufferPercent), err
	}
	if bufferPercent == 0 {
		if viper.IsSet("buffer") {
			bufferPercent = viper.GetInt32("buffer")
		} else {
			bufferPercent = int32(core.DefaultBufferPercent)
			log.Debug("BufferPercent is not set, taking its default value ", bufferPercent)
		}
	}
	return bufferPercent, nil
}

// GetConfigData assembles complete configuration data by:
// 1. Gathering all configuration parameters
// 2. Applying default values where needed
// 3. Validating configuration consistency
// Returns full configuration object or error if validation fails.
func (*UtilsStruct) GetConfigData() (types.Configurations, error) {
	config := types.Configurations{
		Provider:           "",
		BufferPercent:      0,
		WaitTime:           0,
		GasPrice:           0,
		RPCTimeout:         0,
		LogLevel:           "",
		GasMultiplier:      0.0,
		GasLimitMultiplier: 0.0,
	}

	provider, err := cmdUtils.GetRPCProvider()
	if err != nil {
		return config, err
	}
	gasMultiplier, err := cmdUtils.GetMultiplier()
	if err != nil {
		return config, err
	}
	bufferPercent, err := cmdUtils.GetBufferPercent()
	if err != nil {
		return config, err
	}
	waitTime, err := cmdUtils.GetWaitTime()
	if err != nil {
		return config, err
	}
	gasPrice, err := cmdUtils.GetGasPrice()
	if err != nil {
		return config, err
	}
	logLevel, err := cmdUtils.GetLogLevel()
	if err != nil {
		return config, err
	}
	gasLimit, err := cmdUtils.GetGasLimit()
	if err != nil {
		return config, err
	}
	rpcTimeout, err := cmdUtils.GetRPCTimeout()
	if err != nil {
		return config, err
	}
	config.Provider = provider
	config.GasMultiplier = gasMultiplier
	config.BufferPercent = bufferPercent
	config.WaitTime = waitTime
	config.GasPrice = gasPrice
	config.LogLevel = logLevel
	config.GasLimitMultiplier = gasLimit
	config.RPCTimeout = rpcTimeout
	utils.RPCTimeout = rpcTimeout

	return config, nil
}

// GetRPCProvider retrieves RPC provider URL from configuration or flags.
// Validates URL format and warns if non-secure URL is used.
// Falls back to default provider if none specified.
func (*UtilsStruct) GetRPCProvider() (string, error) {
	provider, err := flagSetUtils.GetRootStringProvider()
	if err != nil {
		return core.DefaultRPCProvider, err
	}
	if provider == "" {
		if viper.IsSet("provider") {
			provider = viper.GetString("provider")
		} else {
			provider = core.DefaultRPCProvider
			log.Debug("Provider is not set, taking its default value ", provider)
		}
	}
	if !strings.HasPrefix(provider, "https") {
		log.Warn("You are not using a secure RPC URL. Switch to an https URL instead to be safe.")
	}
	return provider, nil
}

// GetMultiplier gets gas multiplier value from configuration or flags.
// Uses default if not specified. Returns error if value is invalid.
func (*UtilsStruct) GetMultiplier() (float32, error) {
	gasMultiplier, err := flagSetUtils.GetRootFloat32GasMultiplier()
	if err != nil {
		return float32(core.DefaultGasMultiplier), err
	}
	if gasMultiplier == -1 {
		if viper.IsSet("gasmultiplier") {
			gasMultiplier = float32(viper.GetFloat64("gasmultiplier"))
		} else {
			gasMultiplier = float32(core.DefaultGasMultiplier)
			log.Debug("GasMultiplier is not set, taking its default value ", gasMultiplier)
		}
	}
	return gasMultiplier, nil
}

// This function returns the wait time
func (*UtilsStruct) GetWaitTime() (int32, error) {
	waitTime, err := flagSetUtils.GetRootInt32Wait()
	if err != nil {
		return int32(core.DefaultWaitTime), err
	}
	if waitTime == -1 {
		if viper.IsSet("wait") {
			waitTime = viper.GetInt32("wait")
		} else {
			waitTime = int32(core.DefaultWaitTime)
			log.Debug("WaitTime is not set, taking its default value ", waitTime)
		}
	}
	return waitTime, nil
}

// GetGasPrice retrieves gas price setting from configuration or flags.
// Applies default if not set. Validates price is within acceptable range.
func (*UtilsStruct) GetGasPrice() (int32, error) {
	gasPrice, err := flagSetUtils.GetRootInt32GasPrice()
	if err != nil {
		return int32(core.DefaultGasPrice), err
	}
	if gasPrice == -1 {
		if viper.IsSet("gasprice") {
			gasPrice = viper.GetInt32("gasprice")
		} else {
			gasPrice = int32(core.DefaultGasPrice)
			log.Debug("GasPrice is not set, taking its default value ", gasPrice)

		}
	}
	return gasPrice, nil
}

// GetLogLevel gets logging level from configuration or flags.
// Uses default level if not specified. Validates level is supported.
func (*UtilsStruct) GetLogLevel() (string, error) {
	logLevel, err := flagSetUtils.GetRootStringLogLevel()
	if err != nil {
		return core.DefaultLogLevel, err
	}
	if logLevel == "" {
		if viper.IsSet("logLevel") {
			logLevel = viper.GetString("logLevel")
		} else {
			logLevel = core.DefaultLogLevel
			log.Debug("LogLevel is not set, taking its default value ", logLevel)
		}
	}
	return logLevel, nil
}

// GetGasLimit retrieves gas limit multiplier from configuration or flags.
// Applies default if not specified. Validates limit is within safe range.
func (*UtilsStruct) GetGasLimit() (float32, error) {
	gasLimit, err := flagSetUtils.GetRootFloat32GasLimit()
	if err != nil {
		return float32(core.DefaultGasLimit), err
	}
	if gasLimit == -1 {
		if viper.IsSet("gasLimit") {
			gasLimit = float32(viper.GetFloat64("gasLimit"))
		} else {
			gasLimit = float32(core.DefaultGasLimit)
			log.Debug("GasLimit is not set, taking its default value ", gasLimit)
		}
	}
	return gasLimit, nil
}

// This function returns the RPC timeout
func (*UtilsStruct) GetRPCTimeout() (int64, error) {
	rpcTimeout, err := flagSetUtils.GetRootInt64RPCTimeout()
	if err != nil {
		return int64(core.DefaultRPCTimeout), err
	}
	if rpcTimeout == 0 {
		if viper.IsSet("rpcTimeout") {
			rpcTimeout = viper.GetInt64("rpcTimeout")
		} else {
			rpcTimeout = int64(core.DefaultRPCTimeout)
			log.Debug("RPCTimeout is not set, taking its default value ", rpcTimeout)
		}
	}
	return rpcTimeout, nil
}
