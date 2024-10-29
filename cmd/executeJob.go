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

// This function sets the flag appropriately and executes the UpdateJob function
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

	// jobIdStr, err := flagSet.GetString("jobId")
	// utils.CheckError("Error in getting jobId: ", err)

	// jobId, ok := new(big.Int).SetString(jobIdStr, 10)
	// if !ok {
	// 	log.Fatal("Invalid JobId format", errors.New("Failed to parse job ID string"))
	// }
	// configPath, err := flagSet.GetString("config")
	// utils.CheckError("Error in getting config path: ", err)

	pipelinePath, err := flagSet.GetString("zen-path")
	utils.CheckError("Error in getting pipeline path: ", err)

	isAdmin, err := flagSet.GetBool("isAdmin")
	utils.CheckError("Error in getting admin flag: ", err)

	if isAdmin && address != "0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771" {
		log.Fatal("Only Admin can pass the isAdmin Flag")
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
	if err := cmdUtils.ExecuteJob(ctx, client, config, account, isAdmin, pipelinePath); err != nil {
		log.WithError(err).Fatal("Job execution failed")
	}
}

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

// 	opts := protoUtils.GetOptions()
// 	stateManagerUtils.WaitForNextState(client, &opts, types.EpochStateAssign)

// 	for {
// 		currentEpoch, currentState, err := cmdUtils.GetEpochAndState(client)
// 		if err != nil {
// 			log.Error(err)
// 		}
// 		log.Infof("State: %s Epoch: %v", utils.UtilsInterface.GetStateName(currentState), currentEpoch)
// 		time.Sleep(5 * time.Second)

// 		switch currentState {
// 		case int64(types.EpochStateAssign):

// 		}

// 	}

// 	// Install dependencies with live logging
// 	log.Info("Starting dependency installation...")
// 	err = pipeline_zen.InstallDeps(pipelinePath)
// 	if err != nil {
// 		log.WithError(err).Fatal("Failed to install dependencies")
// 	}
// 	log.Info("Dependencies installed successfully")

// 	// Hardcoded, to be changed in future
// 	status := types.JobStatusQueued
// 	buffer := 0

// 	// Update job status to Queued
// 	log.Info("Updating job status to Queued...")
// 	jobUpdateTxn, err := cmdUtils.UpdateJobStatus(client, config, types.Account{
// 		Address:  address,
// 		Password: password,
// 	}, jobId, status, uint8(buffer))
// 	if err != nil {
// 		log.WithError(err).Fatal("Failed to update job status to Queued")
// 	}
// 	log.WithField("txHash", jobUpdateTxn.Hex()).Info("Job status updated to Queued")

// 	// Run the TorchTuneWrapper
// 	log.Info("Running TorchTuneWrapper...")
// 	go func() {
// 		output, err := pipeline_zen.RunTorchTuneWrapper(pipelinePath, configPath)
// 		if err != nil {
// 			log.WithError(err).Error("Error running TorchTuneWrapper")
// 			cmdUtils.UpdateJobStatus(client, config, types.Account{
// 				Address:  address,
// 				Password: password,
// 			}, jobId, types.JobStatusFailed, uint8(buffer))
// 			return
// 		}
// 		log.Info("Updating job status to Running...")
// 		runningTxnHash, err := cmdUtils.UpdateJobStatus(client, config, types.Account{
// 			Address:  address,
// 			Password: password,
// 		}, jobId, types.JobStatusRunning, uint8(buffer))
// 		log.WithField("txHash", runningTxnHash.Hex()).Info("Job status updated to Running")

// 		log.Debug("TorchTuneWrapper output: ", output)
// 		log.Info("Job execution initiated. Monitor logs for progress.")
// 	}()

// 	// Update job status to Running
// 	completedJobUpdateTxn, err := cmdUtils.UpdateJobStatus(client, config, types.Account{
// 		Address:  address,
// 		Password: password,
// 	}, jobId, types.JobStatusCompleted, uint8(buffer))
// 	log.WithField("txHash", completedJobUpdateTxn.Hex()).Info("Job status updated to Completed")

// }

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

// This function allows the admin to update an existing job
func (*UtilsStruct) ExecuteJob(ctx context.Context, client *ethclient.Client, config types.Configurations, account types.Account, isAdmin bool, pipelinePath string) error {
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

			if err := cmdUtils.HandleStateTransition(ctx, client, config, account, types.EpochState(state), epoch, isAdmin, pipelinePath); err != nil {
				log.WithError(err).Error("Error handling state transition")
			}
		}
	}
}

func init() {
	rootCmd.AddCommand(executeJobCmd)

	var (
		JobId      string
		Account    string
		Password   string
		ConfigPath string
		ZenPath    string
		IsAdmin    bool
	)

	executeJobCmd.Flags().StringVarP(&JobId, "jobId", "", string(0), "job id")
	executeJobCmd.Flags().StringVarP(&Account, "address", "a", "", "address of the compute provider")
	executeJobCmd.Flags().StringVarP(&Password, "password", "", "", "password path of compute provider to protect the keystore")
	executeJobCmd.Flags().StringVarP(&ConfigPath, "config", "c", "", "path to the job configuration file")
	executeJobCmd.Flags().StringVarP(&ZenPath, "zen-path", "z", "", "path to the pipeline-zen directory")
	executeJobCmd.Flags().BoolVarP(&IsAdmin, "isAdmin", "", false, "whether the executor is an admin")

	AddrErr := executeJobCmd.MarkFlagRequired("address")
	utils.CheckError("Address error : ", AddrErr)
	jobId := executeJobCmd.MarkFlagRequired("jobId")
	utils.CheckError("JobId error : ", jobId)
	configPath := executeJobCmd.MarkFlagRequired("config")
	utils.CheckError("Path error : ", configPath)
	zenPath := executeJobCmd.MarkFlagRequired("zen-path")
	utils.CheckError("Pipeline Path error : ", zenPath)
}
