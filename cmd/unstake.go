// Package cmd provides all functions related to command line
package cmd

import (
	"errors"
	"lumino/core"
	"lumino/core/types"
	"lumino/logger"
	"lumino/pkg/bindings"
	"lumino/utils"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/pflag"

	"github.com/spf13/cobra"
)

var unstakeCmd = &cobra.Command{
	Use:   "unstake",
	Short: "Unstake your luminos",
	Long: `unstake allows user to unstake their sRzrs in the lumino network

Example:	
  ./lumino unstake --address 0x5a0b54d5dc17e0aadc383d2db43b0a0d3e029c4c --value 1000
	`,
	Run: initialiseUnstake,
}

// This function initialises the ExecuteUnstake function
func initialiseUnstake(cmd *cobra.Command, args []string) {
	cmdUtils.ExecuteUnstake(cmd.Flags())
}

// ExecuteUnstake is the entry point for token unstaking process. This function:
// 1. Loads and validates network configuration
// 2. Verifies account and token amount
// 3. Checks for existing unstake locks
// 4. Executes the unstaking transaction
// 5. Monitors transaction confirmation
// Returns early if any validation step fails.
func (*UtilsStruct) ExecuteUnstake(flagSet *pflag.FlagSet) {
	config, err := cmdUtils.GetConfigData()
	utils.CheckError("Error in getting config: ", err)
	log.Debugf("ExecuteUnstake: Config: %+v", config)

	client := protoUtils.ConnectToEthClient(config.Provider)

	address, err := flagSetUtils.GetStringAddress(flagSet)
	utils.CheckError("Error in getting address: ", err)

	logger.SetLoggerParameters(client, address)
	log.Debug("Checking to assign log file...")
	protoUtils.AssignLogFile(flagSet)

	log.Debug("Getting password...")
	password := protoUtils.AssignPassword(flagSet)

	log.Debug("Getting amount in wei...")
	valueInWei, err := cmdUtils.AssignAmountInWei(flagSet)
	utils.CheckError("Error in getting amountInWei: ", err)

	// TODO: might be needed in future
	// protoUtils.CheckEthBalanceIsZero(client, address)

	stakerId, err := protoUtils.AssignStakerId(flagSet, client, address)
	utils.CheckError("StakerId error: ", err)

	unstakeInput := types.UnstakeInput{
		Address:    address,
		Password:   password,
		ValueInWei: valueInWei,
		StakerId:   stakerId,
	}

	log.Debugf("ExecuteUnstake: Calling Unstake() with arguments unstakeInput: %+v", unstakeInput)
	txnHash, err := cmdUtils.Unstake(config, client, unstakeInput)
	utils.CheckError("Unstake Error: ", err)
	if txnHash != core.NilHash {
		err = protoUtils.WaitForBlockCompletion(client, txnHash.String())
		utils.CheckError("Error in WaitForBlockCompletion for unstake: ", err)
	}
}

// Unstake initiates the unstaking process for staked tokens in the Lumino network.
// This critical function:
// 1. Verifies staker existence and status
// 2. Checks for existing unstake locks
// 3. Constructs and submits unstaking transaction
// 4. Returns transaction hash upon successful submission
// Returns error if unstaking conditions are not met.
func (*UtilsStruct) Unstake(config types.Configurations, client *ethclient.Client, input types.UnstakeInput) (common.Hash, error) {
	txnArgs := types.TransactionOptions{
		Client:         client,
		Password:       input.Password,
		AccountAddress: input.Address,
		Amount:         input.ValueInWei,
		ChainId:        core.ChainID,
		Config:         config,
	}
	stakerId := input.StakerId
	staker, err := protoUtils.GetStaker(client, stakerId)
	if err != nil {
		log.Error("Error in getting staker: ", err)
		return core.NilHash, err
	}
	log.Debugf("Unstake: Staker info: %+v", staker)

	txnArgs.ContractAddress = core.StakeManagerAddress
	txnArgs.MethodName = "unstake"
	txnArgs.ABI = bindings.StakeManagerABI

	unstakeLock, err := protoUtils.GetLock(txnArgs.Client, txnArgs.AccountAddress)
	if err != nil {
		log.Error("Error in getting unstakeLock: ", err)
		return core.NilHash, err
	}
	log.Debugf("Unstake: Unstake lock: %+v", unstakeLock)

	if unstakeLock.Amount.Cmp(big.NewInt(0)) != 0 {
		err := errors.New("existing unstake lock")
		log.Error(err)
		return core.NilHash, err
	}

	txnArgs.Parameters = []interface{}{stakerId, txnArgs.Amount}
	txnOpts := protoUtils.GetTransactionOpts(txnArgs)
	log.Info("Unstaking tokens")
	log.Debugf("Executing Unstake transaction with stakerId = %d, amount = %s", stakerId, txnArgs.Amount)
	txn, err := stakeManagerUtils.Unstake(txnArgs.Client, txnOpts, stakerId, txnArgs.Amount)
	if err != nil {
		log.Error("Error in un-staking: ", err)
		return core.NilHash, err
	}
	log.Info("Transaction hash: ", transactionUtils.Hash(txn))
	return transactionUtils.Hash(txn), nil
}

// Sets up the unstake command with required flags for address and value,
// plus optional flags for password and wei denomination specification.
func init() {
	rootCmd.AddCommand(unstakeCmd)

	var (
		Address         string
		AmountToUnStake string
		Password        string
		WeiLumino       bool
		StakerId        uint32
	)

	unstakeCmd.Flags().StringVarP(&Address, "address", "a", "", "user's address")
	unstakeCmd.Flags().StringVarP(&AmountToUnStake, "value", "v", "0", "value of lumino tokens to un-stake")
	unstakeCmd.Flags().StringVarP(&Password, "password", "", "", "password path to protect the keystore")
	unstakeCmd.Flags().BoolVarP(&WeiLumino, "weiLumino", "", false, "value can be passed in wei")
	unstakeCmd.Flags().Uint32VarP(&StakerId, "stakerId", "", 0, "staker id")

	addrErr := unstakeCmd.MarkFlagRequired("address")
	utils.CheckError("Address error: ", addrErr)
	valueErr := unstakeCmd.MarkFlagRequired("value")
	utils.CheckError("Value error: ", valueErr)

}
