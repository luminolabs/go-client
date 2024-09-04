// Package cmd provides all functions related to command line
package cmd

import (
	"lumino/path"
	"lumino/utils"
	pathPkg "path"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import can be used to import existing accounts into lumino-go",
	Long: `If the user has their private key of an account, they can import that account into lumino-go to perform further operations with lumino-go.
Example:
  ./lumino import --logFile importLogs`,
	Run: initialiseImport,
}

// This function initialises the ExecuteImport function
func initialiseImport(cmd *cobra.Command, args []string) {
	cmdUtils.ExecuteImport(cmd.Flags())
}

// This function sets the flags appropriately and executes the ImportAccount function
func (*UtilsStruct) ExecuteImport(flagSet *pflag.FlagSet) {
	log.Debug("Checking to assign log file...")
	protoUtils.AssignLogFile(flagSet)
	log.Debug("Calling ImportAccount()...")
	account, err := cmdUtils.ImportAccount()
	utils.CheckError("Import error: ", err)
	log.Info("ExecuteImport: Account Address: ", account.Address)
	log.Info("ExecuteImport: Keystore Path: ", account.URL)
}

// This function is used to import existing accounts into
func (*UtilsStruct) ImportAccount() (accounts.Account, error) {
	log.Info("Enter the private key for the account that you want to import")
	privateKey := protoUtils.PrivateKeyPrompt()
	// Remove 0x from the private key
	privateKey = strings.TrimPrefix(privateKey, "0x")
	log.Info("Enter password to protect keystore file")
	log.Info("The password should be of minimum 8 characters containing least 1 uppercase, lowercase, digit and special character.")
	password := protoUtils.PasswordPrompt()
	luminoPath, err := protoUtils.GetDefaultPath()
	if err != nil {
		log.Error("Error in fetching .lumino directory")
		return accounts.Account{Address: common.Address{0x00}}, err
	}
	log.Debug("ImportAccount: .lumino directory path: ", luminoPath)
	priv, err := cryptoUtils.HexToECDSA(privateKey)
	if err != nil {
		log.Error("Error in parsing private key")
		return accounts.Account{Address: common.Address{0x00}}, err
	}
	keystoreDir := pathPkg.Join(luminoPath, "keystore_files")
	if _, err := path.OSUtilsInterface.Stat(keystoreDir); path.OSUtilsInterface.IsNotExist(err) {
		mkdirErr := path.OSUtilsInterface.Mkdir(keystoreDir, 0700)
		if mkdirErr != nil {
			return accounts.Account{Address: common.Address{0x00}}, mkdirErr
		}
	}
	log.Debug("ImportAccount: Keystore directory path: ", keystoreDir)
	log.Debug("Importing the account...")
	account, err := keystoreUtils.ImportECDSA(keystoreDir, priv, password)
	if err != nil {
		log.Error("Error in importing account")
		return accounts.Account{Address: common.Address{0x00}}, err
	}
	log.Info("Account imported...")
	return account, nil
}

func init() {
	rootCmd.AddCommand(importCmd)

}
