package cmd

import (
	"context"
	"fmt"
	"lumino/core"
	"lumino/core/types"
	"lumino/logger"
	"lumino/utils"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
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

// initializeStake is the entry point for the stake command
// It prepares the necessary arguments and calls executeStake
func initializeStake(cmd *cobra.Command, args []string) {
	// Create a background context
	ctx := context.Background()

	// Connect to the Ethereum client
	client := protoUtils.ConnectToEthClient(core.DefaultRPCProvider)

	// Get the stake amount from the command flags
	stakeAmount, _ := cmd.Flags().GetString("amount")

	// Parse the stake amount from string to big.Int
	amount, err := utils.ParseBigInt(stakeAmount)
	if err != nil {
		logger.Fatal("Invalid stake amount:", err)
	}

	// Get the stake address from the command flags
	stakeAddress, _ := cmd.Flags().GetString("address")

	// Convert the address string to Ethereum address type
	address := common.HexToAddress(stakeAddress)

	// Get the password from the command flags
	password, _ := cmd.Flags().GetString("password")

	// Prepare the arguments for staking
	stakeArgs := types.StakeArgs{ //check inside core - > staker.go and define there
		Client:   client,
		Address:  address,
		Amount:   amount,
		Password: password,
	}

	// Execute the staking process
	if err := executeStake(ctx, stakeArgs); err != nil {

		logger.Fatal("Stake operation failed:", err)
	}
}

// executeStake performs the main staking logic
func executeStake(ctx context.Context, args types.StakeArgs) error {
	if err := validateStakeArgs(ctx, args); err != nil {
		return err
	}

	// Directly stake tokens since no approval is needed
	return stakeTokens(ctx, args)
}

// validateStakeArgs checks if the provided arguments are valid
func validateStakeArgs(ctx context.Context, args types.StakeArgs) error {
	// Validate the provided password
	if err := core.ValidatePassword(args.Address, args.Password); err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	// Check the LUMINO balance of the staker
	balance, err := core.GetLuminoBalanceForStaker(ctx, args.Client, args.Address)
	if err != nil {
		return fmt.Errorf("failed to get LUMINO balance: %w", err)
	}

	// Ensure the staker has sufficient balance
	if balance.Cmp(args.Amount) < 0 {
		return fmt.Errorf("insufficient LUMINO balance. Have %s, need %s", balance.String(), args.Amount.String())
	}

	// Check the minimum required stake amount
	//minStake, err := core.GetMinimumStake(ctx, args.Client)
	//if err != nil {
	//	return fmt.Errorf("failed to get minimum stake amount: %w", err)
	//}

	minStakeBigInt := big.NewInt(int64(core.MinimumStake))

	// Ensure the stake amount meets the minimum requirement
	if args.Amount.Cmp(minStakeBigInt) < 0 {
		return fmt.Errorf("stake amount (%s) is below minimum required (%s)", args.Amount.String(), minStakeBigInt.String())
	}

	return nil
}

func stakeTokens(ctx context.Context, args types.StakeArgs) error {
	// Logging the start of the staking process
	logger.Info("Preparing to stake LUMINO tokens...")

	// Step 1: Prepare the Transaction
	transactOpts, err := utils.PrepareStakeTransaction(ctx, args.Client, args.Address, args.Amount, args.Password)
	if err != nil {
		return fmt.Errorf("failed to prepare stake transaction: %w", err)
	}

	// Step 2: Get the StakeManager Contract Instance
	stakeManager := utilsInterface.GetStakeManager(args.Client)

	//stakeManager, err := core.GetStakeManagerContract(args.Client)
	//if err != nil {
	//	return fmt.Errorf("failed to get stake manager contract: %w", err)
	//}

	// Step 3: Stake the Tokens using the Contract Instance
	logger.Info("Staking LUMINO tokens...")
	//arrayOfValues := utils.InvokeFunctionWithTimeout(stakeManager, "Stake", transactOpts, args.Amount)

	epoch, err := protoUtils.GetEpoch(args.Client)

	transaction, err := stakeManager.Stake(transactOpts, epoch, args.Amount, "")
	//tx, err := stakeManager.Stake(_)

	if err != nil {
		return fmt.Errorf("failed to stake tokens: %w", err)
	}

	// Step 4: Wait for the Transaction to be Mined
	logger.Info("Waiting for stake transaction to be mined...")
	receipt, err := bind.WaitMined(ctx, args.Client, transaction)
	if err != nil {
		return fmt.Errorf("failed waiting for stake transaction: %w", err)
	}

	// Step 5: Check Transaction Status
	if receipt.Status == 0 {
		return fmt.Errorf("stake transaction failed")
	}

	// Successful staking operation
	logger.Info("Successfully staked", args.Amount, "LUMINO tokens")
	return nil
}

func (*UtilsStruct) GetOptions() bind.CallOpts {
	block, _ := new(big.Int).SetString("", 10)
	return bind.CallOpts{
		Pending:     false,
		BlockNumber: block,
		Context:     context.Background(),
	}
}

// init function is called when the package is initialized
func init() {
	// Add the stake command to the root command
	rootCmd.AddCommand(stakeCmd)

	var (
		stakeAmount  string // Amount of LUMINO tokens to stake
		stakeAddress string // Address of the staker
		password     string // Password for the staker's account
	)

	// Define flags for the stake command
	stakeCmd.Flags().StringVar(&stakeAmount, "amount", "", "Amount of LUMINO tokens to stake")
	stakeCmd.Flags().StringVar(&stakeAddress, "address", "", "Address of the staker")
	stakeCmd.Flags().StringVar(&password, "password", "", "Password for the staker's account")

	// Mark flags as required
	stakeCmd.MarkFlagRequired("amount")
	stakeCmd.MarkFlagRequired("address")
	stakeCmd.MarkFlagRequired("password")
}
