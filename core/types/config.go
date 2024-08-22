package types

type Configurations struct {
	BufferPercent      int32
	WaitTime           int32
	GasPrice           int32
	RPCTimeout         int64
	Provider           string
	LogLevel           string
	GasMultiplier      float32
	GasLimitMultiplier float32
}
