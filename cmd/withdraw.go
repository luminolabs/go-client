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

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// WithdrawCmd represents the withdraw command
var withdrawCmd = &cobra.Command{
	Use:   "withdraw",
	Short: "withdraw withdraws your lumino Tokens once unstake lock has passed",
	Long:  `withdraw has to be called once the unstake lock period is over to get back all the lumino tokens into your account`,
	Run:   initializeWithdraw,
}

// This function initialises the ExecuteWithdraw function
func initializeWithdraw(cmd *cobra.Command, args []string) {
	cmdUtils.ExecuteWithdraw(cmd.Flags())
}

// This function sets the flag appropriately and executes the Withdraw function
func (*UtilsStruct) ExecuteWithdraw(flagSet *pflag.FlagSet) {
	config, err := cmdUtils.GetConfigData()
	utils.CheckError("Error in getting config: ", err)
	log.Debugf("ExecuteWithdraw: Config: %+v", config)

	client := protoUtils.ConnectToEthClient(config.Provider)

	address, err := flagSetUtils.GetStringAddress(flagSet)
	utils.CheckError("Error in getting address: ", err)
	log.Debug("ExecuteWithdraw: Address: ", address)

	logger.SetLoggerParameters(client, address)
	log.Debug("Checking to assign log file...")
	protoUtils.AssignLogFile(flagSet)

	log.Debug("Getting password...")
	password := protoUtils.AssignPassword(flagSet)

	// TODO: might be needed in future
	// protoUtils.CheckEthBalanceIsZero(client, address)

	stakerId, err := protoUtils.AssignStakerId(flagSet, client, address)
	utils.CheckError("Error in fetching stakerId:  ", err)
	log.Debug("ExecuteWithdraw: StakerId: ", stakerId)

	log.Debugf("ExecuteWithdraw: Calling HandleWithdrawLock with arguments account address = %s, stakerId = %d", address, stakerId)
	txn, err := cmdUtils.HandleUnstakeLock(client, types.Account{
		Address:  address,
		Password: password,
	}, config, stakerId)

	utils.CheckError("Withdraw error: ", err)
	if txn != core.NilHash {
		err = protoUtils.WaitForBlockCompletion(client, txn.String())
		utils.CheckError("Error in WaitForBlockCompletion for withdraw: ", err)
	}
}

// This function handles the Unstake lock
func (*UtilsStruct) HandleUnstakeLock(client *ethclient.Client, account types.Account, configurations types.Configurations, stakerId uint32) (common.Hash, error) {
	// _, err := cmdUtils.WaitForAppropriateState(client, "initiateWithdraw", 0, 1, 4)
	// if err != nil {
	// 	log.Error("Error in fetching epoch: ", err)
	// 	return core.NilHash, err
	// }

	unstakeLock, err := protoUtils.GetLock(client, account.Address)
	if err != nil {
		return core.NilHash, err
	}
	log.Debugf("HandleWithdrawLock: Withdraw lock: %+v", unstakeLock)

	if unstakeLock.UnlockAfter.Cmp(big.NewInt(0)) == 0 {
		log.Error("unstake command not called before withdrawing lumino tokens!")
		return core.NilHash, errors.New("unstake Lumino Tokens before withdrawing")
	}

	epoch, err := protoUtils.GetEpoch(client)
	if err != nil {
		log.Error("Error in fetching epoch")
		return core.NilHash, err
	}
	log.Debug("HandleWithdrawLock: Epoch: ", epoch)

	if big.NewInt(int64(epoch)).Cmp(unstakeLock.UnlockAfter) >= 0 {
		txnArgs := types.TransactionOptions{
			Client:          client,
			Password:        account.Password,
			AccountAddress:  account.Address,
			ChainId:         core.ChainID,
			Config:          configurations,
			ContractAddress: core.StakeManagerAddress,
			MethodName:      "withdraw",
			ABI:             bindings.StakeManagerABI,
			Parameters:      []interface{}{stakerId},
		}
		txnOpts := protoUtils.GetTransactionOpts(txnArgs)
		log.Debug("HandleWithdrawLock: Calling Withdraw() with arguments stakerId = ", stakerId)
		return cmdUtils.Withdraw(client, txnOpts, stakerId)
	}
	return core.NilHash, errors.New("unstakeLock period not over yet! Please try after some time")
}

// This function withdraws your proto once withdraw lock has passed
func (*UtilsStruct) Withdraw(client *ethclient.Client, txnOpts *bind.TransactOpts, stakerId uint32) (common.Hash, error) {
	log.Info("Unlocking funds...")

	log.Debug("Executing Withdraw transaction with stakerId = ", stakerId)
	txn, err := stakeManagerUtils.Withdraw(client, txnOpts, stakerId)
	if err != nil {
		log.Error("Error in unlocking funds")
		return core.NilHash, err
	}

	log.Info("Txn Hash: ", transactionUtils.Hash(txn))

	return transactionUtils.Hash(txn), nil
}

func init() {
	rootCmd.AddCommand(withdrawCmd)
	var (
		Address  string
		Password string
		StakerId uint32
	)

	withdrawCmd.Flags().StringVarP(&Address, "address", "a", "", "address of the user")
	withdrawCmd.Flags().StringVarP(&Password, "password", "", "", "password path of user to protect the keystore")
	withdrawCmd.Flags().Uint32VarP(&StakerId, "stakerId", "", 0, "password path of user to protect the keystore")

	addrErr := withdrawCmd.MarkFlagRequired("address")
	utils.CheckError("Address error: ", addrErr)
}
