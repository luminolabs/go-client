// Package cmd provides all functions related to command line
package cmd

import (
	"context"
	"errors"
	"lumino/core"
	"lumino/core/types"
	"lumino/logger"
	"lumino/pkg/bindings"
	"lumino/utils"
	"math/big"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	executionState types.JobExecutionState
	stateMutex     sync.RWMutex
)

var executeJobCmd = &cobra.Command{
	Use:   "executeJob",
	Short: "[COMPUTE PROVIDER ONLY]executeJob can be used to execute an existing job",
	Long: `A job consists of parameters and config for training and fine-tuning the model. The executeJob command can be used to execute an existing ML Job, using the ML-pipeline package.

Example:
  ./lumino executeJob -a 0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771 --jobId 1 --config /path/to/config --pipeline-path /path/to/pipeline-zen  
  [FOR ADMIN]
  ./lumino executeJob -a 0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771 --jobId 1 --config /path/to/config --pipeline-path /path/to/pipeline-zen  --isAdmin

Note: 
  This command only works for the compute provider.
`,
	Run: initialiseExecuteJob,
}

// This function initialises the ExecuteUpdateJob function
func initialiseExecuteJob(cmd *cobra.Command, args []string) {
	cmdUtils.RunExecuteJob(cmd.Flags())
}

// RunExecuteJob is the entry point for job execution that sets up the execution environment
// and initiates job processing. This function:
// 1. Validates all input parameters and configuration
// 2. Sets up graceful shutdown handlers
// 3. Initializes execution state tracking
// 4. Launches the main execution loop
// Returns early if validation fails or if admin checks fail.
func (*UtilsStruct) RunExecuteJob(flagSet *pflag.FlagSet) {
	config, err := cmdUtils.GetConfigData()
	utils.CheckError("Error in getting config: ", err)
	log.Debugf("RunExecuteJob: Config: %+v", config)

	client := protoUtils.ConnectToEthClient(config.Provider)

	address, err := flagSetUtils.GetStringAddress(flagSet)
	utils.CheckError("Error in getting address: ", err)

	logger.SetLoggerParameters(client, address)
	log.Debug("Checking to assign log file...")
	protoUtils.AssignLogFile(flagSet)

	log.Debug("Getting password...")
	password := protoUtils.AssignPassword(flagSet)

	pipelinePath, err := flagSet.GetString("zen-path")
	utils.CheckError("Error in getting pipeline path: ", err)

	isAdmin, err := flagSet.GetBool("isAdmin")
	utils.CheckError("Error in getting admin flag: ", err)

	isRandom, err := flagSet.GetBool("isRandom")
	utils.CheckError("Error in getting random flag: ", err)

	if isAdmin && address != "0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771" {
		log.Fatal("Only Admin can pass the isAdmin Flag")
	}
	if isRandom && address != "0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771" {
		log.Fatal("Only Admin can pass the isRandom Flag")
	}

	account := types.Account{
		Address:  address,
		Password: password,
	}

	// Initialize execution state
	executionState = types.JobExecutionState{
		IsJobRunning: false,
		CurrentJob:   nil,
	}

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handleGracefulShutdown(ctx, cancel)

	// Start the main execution loop
	if err := cmdUtils.ExecuteJob(ctx, client, config, account, isAdmin, isRandom, pipelinePath); err != nil {
		log.WithError(err).Fatal("Job execution failed")
	}
}

// Sets up signal handling for graceful shutdown of job execution. Ensures proper cleanup
// of resources and updates job status appropriately when shutdown is triggered.
func handleGracefulShutdown(ctx context.Context, cancel context.CancelFunc) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		select {
		case <-signalChan:
			log.Warn("Received interrupt signal. Starting graceful shutdown...")
			stateMutex.RLock()
			if executionState.CurrentJob != nil {
				log.Info("Currently executing job will be marked as failed")
				// Handle cleanup for current job
			}
			stateMutex.RUnlock()
			log.Info("Press CTRL+C again to force terminate.")
			cancel()
		case <-ctx.Done():
		}
		<-signalChan
		os.Exit(2)
	}()
}

// UpdateJobStatus updates the on-chain status of a job by submitting a transaction with the new status.
// Handles transaction construction, submission and monitoring for all job state transitions.
// Takes care of gas estimation and transaction confirmation.
func (*UtilsStruct) UpdateJobStatus(client *ethclient.Client, config types.Configurations, account types.Account, jobId *big.Int, status types.JobStatus, buffer uint8) (common.Hash, error) {
	if client == nil {
		log.Error("Client is nil")
		return common.Hash{}, errors.New("client is nil")
	}

	if jobId == nil {
		log.Error("JobId is nil")
		return common.Hash{}, errors.New("jobId is nil")
	}

	log.WithFields(logrus.Fields{
		"address": account.Address,
		"jobId":   jobId.String(),
		"status":  status,
	}).Debug("Updating job status")

	if jobsManagerUtils == nil {
		log.Error("JobManagerUtils is nil")
		return common.Hash{}, errors.New("jobManagerUtils is nil")
	}

	log.WithFields(logrus.Fields{
		"jobId":  jobId.String(),
		"status": status,
	}).Debug("Executing updateJobStatus transaction")

	txnArgs := types.TransactionOptions{
		Client:          client,
		AccountAddress:  account.Address,
		Password:        account.Password,
		ChainId:         core.ChainID,
		Config:          config,
		ContractAddress: core.JobManagerAddress,
		MethodName:      "updateJobStatus",
		Parameters:      []interface{}{jobId, uint8(status), buffer},
		ABI:             bindings.JobManagerABI,
	}

	txnOpts := protoUtils.GetTransactionOpts(txnArgs)

	txn, err := jobsManagerUtils.UpdateJobStatus(txnArgs.Client, txnOpts, jobId, uint8(status), buffer)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"jobId":  jobId.String(),
			"status": status,
		}).Error("Failed to update job status")
		return common.Hash{}, err
	}

	if txn == nil {
		log.Error("Transaction is nil")
		return common.Hash{}, errors.New("transaction is nil")
	}

	txnHash := transactionUtils.Hash(txn)
	log.WithFields(logrus.Fields{
		"txHash": txnHash.Hex(),
		"jobId":  jobId.String(),
		"status": status,
	}).Info("Job status update transaction submitted")

	err = protoUtils.WaitForBlockCompletion(txnArgs.Client, txnHash.Hex())
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"jobId":  jobId.String(),
			"status": status,
		}).Error("Failed to wait for block completion")
		return common.Hash{}, err
	}

	return txnHash, nil
}

// ExecuteJob is the core job execution function that manages the job lifecycle through various network states.
// Continuously monitors the network state and responds to changes by:
// 1. Managing state transitions
// 2. Handling job execution updates
// 3. Coordinating with the blockchain for job progression
// Uses a ticker to periodically check and update job status.
func (*UtilsStruct) ExecuteJob(ctx context.Context, client *ethclient.Client, config types.Configurations, account types.Account, isAdmin bool, isRandom bool, pipelinePath string) error {
	ticker := time.NewTicker(time.Duration(core.StateCheckInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			epoch, state, err := cmdUtils.GetEpochAndState(client)
			if err != nil {
				log.WithError(err).Error("Failed to get current state and epoch")
				continue
			}

			stateMutex.Lock()
			executionState.CurrentState = types.EpochState(state)
			executionState.CurrentEpoch = epoch
			stateMutex.Unlock()

			log.WithFields(logrus.Fields{
				"state": utils.UtilsInterface.GetStateName(state),
				"epoch": epoch,
			}).Debug("Current network state")

			if err := cmdUtils.HandleStateTransition(ctx, client, config, account, types.EpochState(state), epoch, isAdmin, isRandom, pipelinePath); err != nil {
				log.WithError(err).Error("Error handling state transition")
			}
		}
	}
}

func init() {
	rootCmd.AddCommand(executeJobCmd)

	var (
		Account  string
		Password string
		ZenPath  string
		IsAdmin  bool
		IsRandom bool
	)

	executeJobCmd.Flags().StringVarP(&Account, "address", "a", "", "address of the compute provider")
	executeJobCmd.Flags().StringVarP(&Password, "password", "", "", "password path of compute provider to protect the keystore")
	executeJobCmd.Flags().StringVarP(&ZenPath, "zen-path", "z", "", "path to the pipeline-zen directory")
	executeJobCmd.Flags().BoolVarP(&IsAdmin, "isAdmin", "", false, "whether the executor is an admin")
	executeJobCmd.Flags().BoolVarP(&IsRandom, "isRandom", "", false, "whether the job to be assigned in random manner or just to admin")

	AddrErr := executeJobCmd.MarkFlagRequired("address")
	utils.CheckError("Address error : ", AddrErr)
	zenPath := executeJobCmd.MarkFlagRequired("zen-path")
	utils.CheckError("Pipeline Path error : ", zenPath)
}
