package utils

import "github.com/spf13/pflag"

func IntiliaseLuminoUtils(optionsPackageStruct OptionsPackageStruct) Utils {
	UtilsInterface = optionsPackageStruct.UtilsInterface
	FlagSetInterface = optionsPackageStruct.FlagSetInterface
	return &UtilsStruct{}
}

func (u UtilsStruct) GetUint32(flagSet *pflag.FlagSet, name string) (uint32, error) {
	return flagSet.GetUint32(name)
}

func (f FLagSetStruct) GetLogFileName(flagSet *pflag.FlagSet) (string, error) {
	return flagSet.GetString("logFile")
}
