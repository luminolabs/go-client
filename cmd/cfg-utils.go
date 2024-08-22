package cmd

import (
	"lumino/core"
	"lumino/core/types"
	"strings"

	"github.com/spf13/viper"
)

// GetBufferPercent retrieves the buffer percent value from flags or config.
// It returns the buffer percent as an int32 and an error if retrieval fails.
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

func (*UtilsStruct) GetConfigData() (types.Configurations, error) {
	config := types.Configurations{
		Provider:      "",
		BufferPercent: 0,
	}

	provider, err := cmdUtils.GetRPCProvider()
	if err != nil {
		return config, err
	}
	bufferPercent, err := cmdUtils.GetBufferPercent()
	if err != nil {
		return config, err
	}
	config.Provider = provider
	config.BufferPercent = bufferPercent

	return config, nil
}

// This function returns the provider
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
