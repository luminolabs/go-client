package main

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"internal/blockchain"
	"internal/keymanager"
	"internal/statesync"
	"internal/transactionmanager"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Global variables for client components
var (
	log          = logrus.New()
	client       *blockchain.Client
	keyManager   *keymanager.Manager
	stateManager *statesync.Manager
	txManager    *transactionmanager.Manager
)

func init() {
	// Configure logging
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
}

func main() {
	// Set up the root command
	var rootCmd = &cobra.Command{Use: "luminocli"}

	// Define the stake command
	var stakeCmd = &cobra.Command{
		Use:   "stake [amount]",
		Short: "Stake tokens in the Lumino network",
		Args:  cobra.ExactArgs(1),
		Run:   runStake,
	}

	// Define the unstake command
	var unstakeCmd = &cobra.Command{
		Use:   "unstake [amount]",
		Short: "Unstake tokens from the Lumino network",
		Args:  cobra.ExactArgs(1),
		Run:   runUnstake,
	}

	// Define the init command
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize the Lumino client",
		Run:   runInit,
	}

	// Add subcommands to the root command
	rootCmd.AddCommand(stakeCmd, unstakeCmd, initCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// runInit initializes the Lumino client components
func runInit(cmd *cobra.Command, args []string) {
	var err error

	// Initialize blockchain client
	client, err = blockchain.NewClient("http://localhost:8545") // Replace with actual RPC URL
	if err != nil {
		log.Fatalf("Failed to create blockchain client: %v", err)
	}

	// Initialize key manager
	keyManager, err = keymanager.NewManager("your_private_key_here") // Replace with actual private key
	if err != nil {
		log.Fatalf("Failed to create key manager: %v", err)
	}

	// Initialize state manager
	stateManager = statesync.NewManager(client.EthClient(), 300, 4) // 5-minute epochs, 4 states

	// Initialize transaction manager
	txManager = transactionmanager.NewManager(client.EthClient())

	log.Info("Lumino client initialized successfully")
}

// runStake handles the staking process
func runStake(cmd *cobra.Command, args []string) {
	// Parse stake amount
	amount, ok := new(big.Int).SetString(args[0], 10)
	if !ok {
		log.Fatalf("Invalid stake amount: %s", args[0])
	}

	// Create context
	ctx := context.Background()

	// Initiate staking transaction
	tx, err := client.Stake(ctx, amount)
	if err != nil {
		log.Fatalf("Failed to stake: %v", err)
	}

	// Wait for transaction receipt
	receipt, err := txManager.WaitForReceipt(ctx, tx)
	if err != nil {
		log.Fatalf("Failed to get transaction receipt: %v", err)
	}

	// Log successful staking
	log.WithFields(logrus.Fields{
		"txHash": receipt.TxHash.Hex(),
		"amount": amount.String(),
	}).Info("Staking successful")
}

// runUnstake handles the unstaking process
func runUnstake(cmd *cobra.Command, args []string) {
	// Parse unstake amount
	amount, ok := new(big.Int).SetString(args[0], 10)
	if !ok {
		log.Fatalf("Invalid unstake amount: %s", args[0])
	}

	// Create context
	ctx := context.Background()

	// Initiate unstaking transaction
	tx, err := client.Unstake(ctx, amount)
	if err != nil {
		log.Fatalf("Failed to unstake: %v", err)
	}

	// Wait for transaction receipt
	receipt, err := txManager.WaitForReceipt(ctx, tx)
	if err != nil {
		log.Fatalf("Failed to get transaction receipt: %v", err)
	}

	// Log successful unstaking
	log.WithFields(logrus.Fields{
		"txHash": receipt.TxHash.Hex(),
		"amount": amount.String(),
	}).Info("Unstaking successful")
}
