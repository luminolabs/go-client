package cmd

import (
	"lumino/core"

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
