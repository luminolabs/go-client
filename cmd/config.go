package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	RPCUrl     string `json:"rpc_url"`
	PrivateKey string `json:"private_key"`
	LogLevel   string `json:"log_level"`
}

func loadConfig(cfgFile string) (*Config, error) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".luminocli")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	var config Config
	err := viper.Unmarshal(&config)
	return &config, err
}

func saveConfig(config *Config, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}
