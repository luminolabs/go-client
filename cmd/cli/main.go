package main

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"github.com/luminolabs/lumino-go-client/internal/blockchain"
	"github.com/luminolabs/lumino-go-client/internal/keymanager"
	"github.com/luminolabs/lumino-go-client/internal/statesync"
	"github.com/luminolabs/lumino-go-client/internal/transactionmanager"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var log = logrus.New()
var client *blockchain.Client
var keyManager *keymanager.Manager
var stateManager *statesync.Manager
var txManager *transactionmanager.Manager

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
}

func main() {
	var rootCmd = &cobra.Command{Use: "luminocli"}

	var stakeCmd = &cobra.Command{
		Use:   "stake [amount]",
		Short: "Stake tokens in the Lumino network",
		Args:  cobra.ExactArgs(1),
		Run:   runStake,
	}

	var unstakeCmd = &cobra.Command{
		Use:   "unstake [amount]",
		Short: "Unstake tokens from the Lumino network",
		Args:  cobra.ExactArgs(1),
		Run:   runUnstake,
	}

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize the Lumino client",
		Run:   runInit,
	}

	rootCmd.AddCommand(stakeCmd, unstakeCmd, initCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runInit(cmd *cobra.Command, args []string) {
	var err error
	client, err = blockchain.NewClient("http://localhost:8545") // Replace with actual RPC URL
	if err != nil {
		log.Fatalf("Failed to create blockchain client: %v", err)
	}

	keyManager, err = keymanager.NewManager("your_private_key_here") // Replace with actual private key
	if err != nil {
		log.Fatalf("Failed to create key manager: %v", err)
	}

	stateManager = statesync.NewManager(300, 4) // 5-minute epochs, 4 states
	txManager = transactionmanager.NewManager(client.EthClient())

	log.Info("Lumino client initialized successfully")
}

func runStake(cmd *cobra.Command, args []string) {
	amount, ok := new(big.Int).SetString(args[0], 10)
	if !ok {
		log.Fatalf("Invalid stake amount: %s", args[0])
	}

	ctx := context.Background()
	tx, err := client.Stake(ctx, amount)
	if err != nil {
		log.Fatalf("Failed to stake: %v", err)
	}

	receipt, err := txManager.WaitForReceipt(ctx, tx)
	if err != nil {
		log.Fatalf("Failed to get transaction receipt: %v", err)
	}

	log.WithFields(logrus.Fields{
		"txHash": receipt.TxHash.Hex(),
		"amount": amount.String(),
	}).Info("Staking successful")
}

func runUnstake(cmd *cobra.Command, args []string) {
	amount, ok := new(big.Int).SetString(args[0], 10)
	if !ok {
		log.Fatalf("Invalid unstake amount: %s", args[0])
	}

	ctx := context.Background()
	tx, err := client.Unstake(ctx, amount)
	if err != nil {
		log.Fatalf("Failed to unstake: %v", err)
	}

	receipt, err := txManager.WaitForReceipt(ctx, tx)
	if err != nil {
		log.Fatalf("Failed to get transaction receipt: %v", err)
	}

	log.WithFields(logrus.Fields{
		"txHash": receipt.TxHash.Hex(),
		"amount": amount.String(),
	}).Info("Unstaking successful")
}
