package cmd

import (
	"context"
	"lumino/cmd/systemspecs"
	"lumino/core"
	"lumino/core/types"
	"lumino/logger"
	"lumino/pkg/bindings"
	"lumino/utils"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// stakeCmd represents the stake command
var stakeCmd = &cobra.Command{
	Use:   "stake",
	Short: "Stake LUMINO tokens",
	Long: `Stake LUMINO tokens in the Lumino network.

This command allows users to stake their LUMINO tokens in the network.
Staking is required to participate in network operations and earn rewards.

Example:
  ./lumino stake --address 0x1234567890123456789012345678901234567890 --amount 1000 --password mySecurePassword`,
	Run: initializeStake,
}

func initializeStake(cmd *cobra.Command, args []string) {
	cmdUtils.ExecuteStake(cmd.Flags())
}

// ExecuteStake manages the staking workflow from user input to transaction submission.
// This function orchestrates the entire staking process:
// 1. Validates configuration and connects to network
// 2. Checks account balance against stake amount
// 3. Collects and validates machine specifications
// 4. Executes the staking transaction
// Returns early if validation fails or if transaction encounters errors.
func (*UtilsStruct) ExecuteStake(flagSet *pflag.FlagSet) {
	config, err := cmdUtils.GetConfigData()
	utils.CheckError("Error in getting config: ", err)
	log.Debugf("ExecuteStake: Config: %+v", config)

	client := protoUtils.ConnectToEthClient(config.Provider)

	address, err := flagSetUtils.GetStringAddress(flagSet)
	utils.CheckError("Error in getting address: ", err)
	log.Debug("ExecuteStake: Address: ", address)

	logger.SetLoggerParameters(client, address)
	log.Debug("Checking to assign log file...")
	protoUtils.AssignLogFile(flagSet)

	log.Debug("Getting password...")
	password := protoUtils.AssignPassword(flagSet)

	// Check the LUMINO balance of the staker
	balance, err := protoUtils.FetchBalance(context.Background(), client, common.HexToAddress(address))
	utils.CheckError("Failed to get LUMINO balance:"+address, err)

	log.Debug("Getting amount in wei...")
	valueInWei, err := cmdUtils.AssignAmountInWei(flagSet)
	utils.CheckError("Error in getting amount: ", err)
	log.Debug("ExecuteStake: Amount in wei: ", valueInWei)

	log.Debug("Checking for sufficient balance...")
	protoUtils.CheckAmountAndBalance(valueInWei, balance)

	// TODO: fetch minStake from contracts in Future
	minStakeBigInt := big.NewInt(int64(core.MinimumStake))

	// Ensure the stake amount meets the minimum requirement
	if valueInWei.Cmp(minStakeBigInt) < 0 {
		logger.Fatal("Stake amount", valueInWei.String(), "is below minimum required", minStakeBigInt.String())
	}

	stakerId, err := protoUtils.GetStakerId(client, address)
	utils.CheckError("Error in getting stakerId: ", err)
	log.Debug("ExecuteStake: Staker Id: ", stakerId)

	txnArgs := types.TransactionOptions{
		Client:         client,
		AccountAddress: address,
		Password:       password,
		Amount:         valueInWei,
		ChainId:        core.ChainID,
		Config:         config,
	}

	// Get system specifications
	machineSpecs, err := systemspecs.GetSystemSpecs()
	log.Info("Machine Specs in JSON : ", machineSpecs)
	if err != nil {
		log.Error("Failed to get system specifications: ", err)
		machineSpecs = "{}" // Use empty JSON object if specs couldn't be retrieved
	}

	log.Debug("ExecuteStake: Calling StakeTokens() for amount: ", txnArgs.Amount)
	stakeTxnHash, err := cmdUtils.StakeTokens(txnArgs, machineSpecs)
	utils.CheckError("Stake error: ", err)

	err = protoUtils.WaitForBlockCompletion(txnArgs.Client, stakeTxnHash.String())
	utils.CheckError("Error in WaitForBlockCompletion for stake: ", err)
}

// StakeTokens stakes tokens in the Lumino network for a compute provider. This function:
// 1. Validates the staking amount and account balance
// 2. Retrieves current epoch and network state
// 3. Submits staking transaction with machine specifications
// 4. Monitors transaction confirmation
// Returns transaction hash upon successful staking.
func (*UtilsStruct) StakeTokens(txnArgs types.TransactionOptions, machineSpecs string) (common.Hash, error) {
	epoch, err := protoUtils.GetEpoch(txnArgs.Client)
	if err != nil {
		return common.Hash{0x00}, err
	}
	log.Debug("StakeCoins: Epoch: ", epoch)

	txnArgs.ContractAddress = core.StakeManagerAddress
	txnArgs.MethodName = "stake"
	txnArgs.Parameters = []interface{}{epoch, txnArgs.Amount, machineSpecs}
	txnArgs.ABI = bindings.StakeManagerABI
	txnArgs.EtherValue = txnArgs.Amount
	txnOpts := protoUtils.GetTransactionOpts(txnArgs)
	log.Debugf("Executing Stake transaction with epoch = %d, amount = %d, machineSpecs = %s", epoch, txnArgs.Amount, machineSpecs)
	tx, err := stakeManagerUtils.Stake(txnArgs.Client, txnOpts, epoch, txnArgs.Amount, machineSpecs)
	if err != nil {
		return common.Hash{0x00}, err
	}
	log.Info("Txn Hash: ", transactionUtils.Hash(tx).Hex())
	return transactionUtils.Hash(tx), nil
}

// Initializes the staking command with required flags for address, stake value
// and optional flags for password and wei denomination specification.
func init() {
	rootCmd.AddCommand(stakeCmd)

	var (
		stakeValue    string // Amount of LUMINO tokens to stake
		stakerAddress string // Address of the staker
		password      string // Password for the staker's account
		IsWei         bool
	)

	stakeCmd.Flags().StringVarP(&stakeValue, "value", "v", "0", "Amount of LUMINO tokens to stake")
	stakeCmd.Flags().StringVarP(&stakerAddress, "address", "a", "", "Address of the staker")
	stakeCmd.Flags().StringVarP(&password, "password", "", "", "Password for the staker's account")
	stakeCmd.Flags().BoolVarP(&IsWei, "weiValue", "", false, "value passed in wei")

	stakeAmountErr := stakeCmd.MarkFlagRequired("value")
	utils.CheckError("Value error: ", stakeAmountErr)
	stakerAddrErr := stakeCmd.MarkFlagRequired("address")
	utils.CheckError("Address error: ", stakerAddrErr)
}
