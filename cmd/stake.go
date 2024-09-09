package cmd

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"lumino/core"
	"lumino/core/types"
	"lumino/logger"
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

func (*UtilsStruct) ExecuteStake(flagSet *pflag.FlagSet) {
	config, err := cmdUtils.GetConfigData()
	utils.CheckError("Error in getting config: ", err)
	log.Debugf("ExecuteStake: Config: %+v", config)

	client := protoUtils.ConnectToEthClient(config.Provider)
	if client == nil {
		log.Fatal("Failed to connect to Ethereum client")
		return
	}
	logger.SetLoggerParameters(client, "")

	// Retrieve required parameters from the flags
	stakeArgs, err := cmdUtils.GetStakeArgs(flagSet, client)
	utils.CheckError("Error in getting stake arguments: ", err)

	// Execute the staking logic
	err = executeStake(context.Background(), stakeArgs)
	utils.CheckError("Error during staking process: ", err)
}

func (*UtilsStruct) GetStakeArgs(flagSet *pflag.FlagSet, client *ethclient.Client) (types.StakeArgs, error) {
	// Get the stake amount
	stakeAmount, _ := flagSet.GetString("amount")
	amount, err := utils.ParseBigInt(stakeAmount)
	if err != nil {
		return types.StakeArgs{}, err
	}

	// Get the stake address
	stakeAddress, _ := flagSet.GetString("address")
	address := common.HexToAddress(stakeAddress)

	// Get the password
	password, _ := flagSet.GetString("password")

	return types.StakeArgs{
		Client:   client,
		Address:  address,
		Amount:   amount,
		Password: password,
	}, nil
}

func executeStake(ctx context.Context, args types.StakeArgs) error {
	if err := validateStakeArgs(ctx, args); err != nil {
		logger.Error("Validation of stake arguments failed:", err)
		return err
	}

	// Directly stake tokens
	return stakeTokens(ctx, args)
}

func validateStakeArgs(ctx context.Context, args types.StakeArgs) error {
	// Validate the provided password
	if err := core.ValidatePassword(args.Address, args.Password); err != nil {
		logger.Error("Invalid password:", err)
		return err
	}

	// Check the LUMINO balance of the staker
	balance, err := GetLuminoBalanceForStaker(ctx, args.Client, args.Address)
	if err != nil {
		logger.Error("Failed to get LUMINO balance:", err)
		return err
	}

	// Ensure the staker has sufficient balance
	if balance.Cmp(args.Amount) < 0 {
		logger.Error("Insufficient LUMINO balance. Have ", balance.String(), " need ", args.Amount.String())
		return nil
	}

	// Use minstake from constants file
	minStakeBigInt := big.NewInt(int64(core.MinimumStake))

	// Ensure the stake amount meets the minimum requirement
	if args.Amount.Cmp(minStakeBigInt) < 0 {
		logger.Fatal("Stake amount", args.Amount.String(), "is below minimum required", minStakeBigInt.String())
	}

	return nil
}

func stakeTokens(ctx context.Context, args types.StakeArgs) error {
	logger.Info("Preparing to stake LUMINO tokens...")

	transactOpts, err := utils.PrepareStakeTransaction(ctx, args.Client, args.Address, args.Amount, args.Password)
	if err != nil {
		logger.Error("Failed to prepare stake transaction:", err)
		return err
	}

	logger.Debug("TransactOpts:", transactOpts)

	// Get the StakeManager Contract Instance
	utilsInterface := utils.UtilsStruct{}
	stakeManager, err := utilsInterface.GetStakeManager(args.Client)
	if err != nil {
		logger.Error("Failed to get stake manager:", err)
		return err
	}

	// Stake the Tokens using the Contract Instance
	logger.Info("Staking LUMINO tokens...")
	epoch, err := protoUtils.GetEpoch(args.Client)
	if err != nil {
		logger.Error("Failed to get epoch:", err)
		return err
	}

	transaction, err := stakeManager.Stake(transactOpts, epoch, args.Amount, "")
	if err != nil {
		logger.Error("Failed to stake tokens:", err)
		return err
	}

	// Wait for the Transaction to be Mined
	logger.Info("Waiting for stake transaction to be mined...")
	receipt, err := bind.WaitMined(ctx, args.Client, transaction)
	if err != nil {
		logger.Error("Failed waiting for stake transaction:", err)
		return err
	}

	if receipt.Status == 0 {
		logger.Fatal("stake transaction failed")
	}

	logger.Info("Successfully staked ", args.Amount, " LUMINO tokens")
	return nil
}

func GetLuminoBalanceForStaker(ctx context.Context, client *ethclient.Client, address common.Address) (*big.Int, error) {
	// Get the balance of the address in Wei (smallest unit of Ether)
	balance, err := client.BalanceAt(ctx, address, nil)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func init() {
	rootCmd.AddCommand(stakeCmd)

	var (
		stakeAmount  string // Amount of LUMINO tokens to stake
		stakeAddress string // Address of the staker
		password     string // Password for the staker's account
	)

	stakeCmd.Flags().StringVar(&stakeAmount, "amount", "", "Amount of LUMINO tokens to stake")
	stakeCmd.Flags().StringVar(&stakeAddress, "address", "", "Address of the staker")
	stakeCmd.Flags().StringVar(&password, "password", "", "Password for the staker's account")

	stakeCmd.MarkFlagRequired("amount")
	stakeCmd.MarkFlagRequired("address")
	stakeCmd.MarkFlagRequired("password")
}
