// Package cmd provides all functions related to command line
package cmd

import (
	luminoAccounts "lumino/accounts"
	"lumino/utils"
	"path"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create command can be used to create new accounts",
	Long: `For a new user to start doing anything, an account is required. This command helps the user to create a new account secured by a password so that only that user would be able to use the account

Example: 
  ./lumino create --logFile createLogs`,
	Run: initialiseCreate,
}

// This function initialises the ExecuteCreate function
func initialiseCreate(cmd *cobra.Command, args []string) {
	cmdUtils.ExecuteCreate(cmd.Flags())
}

// ExecuteCreate: Main entry point for the account creation process. Handles the flag parsing, sets up logging,
// creates the account and reports the results. Takes flags as input, initializes necessary
// components and executes the creation flow. Returns early with error if account creation fails.
func (*UtilsStruct) ExecuteCreate(flagSet *pflag.FlagSet) {
	log.Debug("Checking to assign log file...")
	protoUtils.AssignLogFile(flagSet)
	log.Info("The password should be of minimum 8 characters containing least 1 uppercase, lowercase, digit and special character.")
	password := protoUtils.AssignPassword(flagSet)
	log.Debug("ExecuteCreate: Calling Create() with argument as input password")
	account, err := cmdUtils.Create(password)
	utils.CheckError("Create error: ", err)
	log.Info("ExecuteCreate: Account address: ", account.Address)
	log.Info("ExecuteCreate: Keystore Path: ", account.URL)
}

// Create: Creates a new account in the keystore with the given password. It returns the created account
// and any error encountered in the process. The account is stored in the keystore directory
// under ~/.lumino/keystore_files.
func (*UtilsStruct) Create(password string) (accounts.Account, error) {
	luminoPath, err := protoUtils.GetDefaultPath()
	if err != nil {
		log.Error("Error in fetching .lumino directory")
		return accounts.Account{Address: common.Address{0x00}}, err
	}
	log.Debug("Create: .lumino directory path: ", luminoPath)
	keystorePath := path.Join(luminoPath, "keystore_files")
	account := luminoAccounts.AccountUtilsInterface.CreateAccount(keystorePath, password)
	return account, nil
}

// Initializes the cobra command for account creation by configuring flags and help text.
// Configures required flags for address and password.
func init() {
	rootCmd.AddCommand(createCmd)

	var (
		Password string
	)

	createCmd.Flags().StringVarP(&Password, "password", "", "", "password file path to protect the keystore")
}
