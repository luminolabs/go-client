// root.go implements the root command for the Lumino CLI.
// It handles configuration initialization, command registration,
// and global flag setup for all subcommands.
package cmd

import (
	"fmt"
	"os"

	"lumino/core"
	"lumino/logger"
	"lumino/path"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Provider           string
	BufferPercent      int32
	WaitTime           int32
	GasPrice           int32
	RPCTimeout         int64
	LogLevel           string
	LogFile            string
	GasMultiplier      float32
	GasLimitMultiplier float32
)

// log is the package-level logger instance
var log = logger.NewLogger()

// Root command for the Lumino CLI. Represents the base command
// when called without any subcommands. Handles version information
// and global help documentation.
var rootCmd = &cobra.Command{
	Version: core.VersionWithMeta,
	Use:     "luminocli",
	Short:   "Lumino CLI is a command line interface for interacting with the Lumino network",
	Long:    "Lumino CLI can be used by the computerProvider to register and participate in the Lumino Protocol. Compute Provider and perform jobs and stay active to earn rewards.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to lumino-cli.")
		err := cmd.Help()
		if err != nil {
			panic(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Initializes the root command by:
// 1. Setting up global flags
// 2. Configuring default values
// 3. Setting up configuration file handling
// 4. Initializing logging
// Must be called before any subcommands are executed.
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&Provider, "provider", "p", "", "provider name")
	rootCmd.PersistentFlags().Float32VarP(&GasMultiplier, "gasmultiplier", "g", -1, "gas multiplier value")
	rootCmd.PersistentFlags().Int32VarP(&BufferPercent, "buffer", "b", 0, "buffer percent")
	rootCmd.PersistentFlags().Int32VarP(&WaitTime, "wait", "w", -1, "wait time")
	rootCmd.PersistentFlags().Int32VarP(&GasPrice, "gasprice", "", -1, "gas price")
	rootCmd.PersistentFlags().StringVarP(&LogLevel, "logLevel", "", "", "log level")
	rootCmd.PersistentFlags().Float32VarP(&GasLimitMultiplier, "gasLimit", "", -1, "gas limit percentage increase")
	rootCmd.PersistentFlags().StringVarP(&LogFile, "logFile", "", "", "name of log file")
	rootCmd.PersistentFlags().Int64VarP(&RPCTimeout, "rpcTimeout", "", 0, "RPC timeout if its not responding")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// Handles configuration initialization including:
// 1. Locating configuration file
// 2. Loading and parsing configuration
// 3. Setting up environment variables
// 4. Initializing logging levels
// Called automatically during startup.
func initConfig() {
	fmt.Println("Entering initConfig")

	if path.PathUtilsInterface == nil {
		fmt.Println("PathUtilsInterface is nil")
		return
	}

	home, err := path.PathUtilsInterface.GetDefaultPath()
	if err != nil {
		fmt.Printf("Error in fetching .lumino directory: %v\n", err)
		return
	}

	fmt.Printf("Home directory: %s\n", home)

	viper.AddConfigPath(home)
	viper.SetConfigName("lumino")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found")
		} else {
			fmt.Printf("Error in reading config: %v\n", err)
		}
	} else {
		fmt.Println("Config file read successfully")
	}

	fmt.Println("About to call setLogLevel")
	setLogLevel()
	fmt.Println("setLogLevel completed")
}

// Configures logging settings based on configuration.
// Updates log levels and outputs detailed configuration
// information in debug mode.
func setLogLevel() {
	fmt.Println("Entering setLogLevel")

	if cmdUtils == nil {
		fmt.Println("cmdUtils is nil in setLogLevel")
		return
	}

	config, err := cmdUtils.GetConfigData()
	if err != nil {
		fmt.Printf("Error getting config data: %v\n", err)
		return
	}

	if log == nil {
		fmt.Println("log is nil in setLogLevel")
		return
	}

	fmt.Printf("LogLevel from config: %s\n", config.LogLevel)

	if config.LogLevel == "debug" {
		log.SetLevel(logrus.DebugLevel)
	}

	fmt.Println("Exiting setLogLevel")

	log.Debug("Config details: ")
	log.Debugf("Provider: %s", config.Provider)
	log.Debugf("Gas Multiplier: %.2f", config.GasMultiplier)
	log.Debugf("Buffer Percent: %d", config.BufferPercent)
	log.Debugf("Wait Time: %d", config.WaitTime)
	log.Debugf("Gas Price: %d", config.GasPrice)
	log.Debugf("Log Level: %s", config.LogLevel)
	log.Debugf("Gas Limit: %.2f", config.GasLimitMultiplier)
	log.Debugf("RPC Timeout: %d", config.RPCTimeout)
}
