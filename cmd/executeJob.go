// Package cmd provides all functions related to command line
package cmd

import (
	"lumino/core"
	"lumino/core/types"
	"lumino/logger"
	pipeline_zen "lumino/pipeline-zen"
	"lumino/pkg/bindings"
	"lumino/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var executeJobCmd = &cobra.Command{
	Use:   "executeJob",
	Short: "[COMPUTE PROVIDER ONLY]executeJob can be used to execute an existing job",
	Long: `A job consists of parameters and config for training and fine-tuning the model. The executeJob command can be used to execute an existing ML Job, using the ML-pipeline package.

Example:
  ./lumino executeJob -a 0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771 --jobId 1 --config /path/to/config 

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
func (u *UtilsStruct) RunExecuteJob(flagSet *pflag.FlagSet) {
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
	// password := protoUtils.AssignPassword(flagSet)

	jobId, err := flagSetUtils.GetUint16JobId(flagSet)
	utils.CheckError("Error in getting jobId: ", err)

	// configPath, err := flagSet.GetString("config")
	// utils.CheckError("Error in getting config path: ", err)

	pipelinePath, err := flagSet.GetString("zen-path")
	utils.CheckError("Error in getting pipeline path: ", err)

	// Install dependencies with live logging
	log.Info("Starting dependency installation...")
	err = pipeline_zen.InstallDeps(pipelinePath)
	if err != nil {
		log.WithError(err).Fatal("Failed to install dependencies")
	}
	log.Info("Dependencies installed successfully")

	txnArgs := types.TransactionOptions{
		Client:          client,
		AccountAddress:  account.Address,
		Password:        account.Password,
		ChainId:         core.ChainID,
		Config:          config,
		ContractAddress: core.JobManagerAddress,
		MethodName:      "updateJobStatus",
		Parameters:      []interface{}{jobId, status, buffer: 0},
		ABI:             bindings.JobManagerABI,
	}

	txnOpts := protoUtils.GetTransactionOpts(txnArgs)

	// Update job status to Queued
	log.Info("Updating job status to Queued...")
	jobUpdateTxn, err = jobsManagerUtils.UpdateJobStatus(client, txnOpts, jobId, types.JobStatusQueued, 0)
	utils.CheckError("Error updating job status to Queued: ", err)

	// TODO:
	// add UpdateJob txns function
	// add types for jobsStatus
	// complete this and push it to a go-routine and stream logs from there
	// to this thread

	// // Run the TorchTuneWrapper
	// log.Info("Running TorchTuneWrapper...")
	// go func() {
	// 	output, err := pipeline_zen.RunTorchTuneWrapper(configPath)
	// 	if err != nil {
	// 		log.WithError(err).Error("Error running TorchTuneWrapper")
	// 		u.updateJobStatus(client, config, jobId, types.JobStatusFailed)
	// 		return
	// 	}
	// 	log.Debug("TorchTuneWrapper output: ", output)
	// 	u.updateJobStatus(client, config, jobId, types.JobStatusCompleted)
	// }()

	// // Update job status to Running
	// log.Info("Updating job status to Running...")
	// err = u.updateJobStatus(client, config, jobId, types.JobStatusRunning)
	// utils.CheckError("Error updating job status to Running: ", err)

	log.Info("Job execution initiated. Monitor logs for progress.")
}

// func (u *UtilsStruct) updateJobStatus(client *ethclient.Client, config types.Configurations, jobId uint16, status types.JobStatus) error {
// 	txnArgs := types.TransactionOptions{
// 		Client:          client,
// 		AccountAddress:  config.Address,
// 		Password:        config.Password,
// 		ChainId:         core.ChainID,
// 		Config:          config,
// 		ContractAddress: core.JobManagerAddress,
// 		MethodName:      "updateJobStatus",
// 		Parameters:      []interface{}{jobId, uint8(status)},
// 		ABI:             bindings.JobManagerABI,
// 	}

// 	txnOpts := protoUtils.GetTransactionOpts(txnArgs)
// 	log.Debugf("Executing updateJobStatus transaction with jobId = %d, status = %d", jobId, status)
// 	txn, err := jobManagerUtils.UpdateJobStatus(txnArgs.Client, txnOpts, jobId, uint8(status))
// 	if err != nil {
// 		return err
// 	}
// 	log.Info("Transaction Hash: ", transactionUtils.Hash(txn))
// 	return protoUtils.WaitForBlockCompletion(txnArgs.Client, transactionUtils.Hash(txn).String())
// }

// This function allows the admin to update an existing job
func (*UtilsStruct) ExecuteJob(client *ethclient.Client, config types.Configurations, jobId uint16) (common.Hash, error) {
	// _, err := cmdUtils.WaitIfCommitState(client, "update job")
	// if err != nil {
	// 	log.Error("Error in fetching state")
	// 	return core.NilHash, err
	// }
	// txnArgs := protoUtils.GetTxnOpts(types.TransactionOptions{
	// 	Client:          client,
	// 	Password:        jobInput.Password,
	// 	AccountAddress:  jobInput.Address,
	// 	ChainId:         core.ChainId,
	// 	Config:          config,
	// 	ContractAddress: core.CollectionManagerAddress,
	// 	MethodName:      "updateJob",
	// 	Parameters:      []interface{}{jobId, jobInput.Weight, jobInput.Power, jobInput.SelectorType, jobInput.Selector, jobInput.Url},
	// 	ABI:             bindings.CollectionManagerABI,
	// })
	// log.Debugf("Executing UpdateJob transaction with arguments jobId = %d, weight = %d, power = %d, selector type = %d, selector = %s, URL = %s", jobId, jobInput.Weight, jobInput.Power, jobInput.SelectorType, jobInput.Selector, jobInput.Url)
	// txn, err := assetManagerUtils.UpdateJob(client, txnArgs, jobId, jobInput.Weight, jobInput.Power, jobInput.SelectorType, jobInput.Selector, jobInput.Url)
	// if err != nil {
	// 	return core.NilHash, err
	// }
	// return transactionUtils.Hash(txn), nil
	return core.NilHash, nil
}

func init() {
	rootCmd.AddCommand(executeJobCmd)

	var (
		JobId      uint16
		Account    string
		Password   string
		ConfigPath string
		ZenPath    string
	)

	executeJobCmd.Flags().Uint16VarP(&JobId, "jobId", "", 0, "job id")
	executeJobCmd.Flags().StringVarP(&Account, "address", "a", "", "address of the compute provider")
	executeJobCmd.Flags().StringVarP(&Password, "password", "", "", "password path of compute provider to protect the keystore")
	executeJobCmd.Flags().StringVarP(&ConfigPath, "config", "c", "", "path to the job configuration file")
	executeJobCmd.Flags().StringVarP(&ZenPath, "zen-path", "z", "", "path to the pipeline-zen directory")

	AddrErr := executeJobCmd.MarkFlagRequired("address")
	utils.CheckError("Address error : ", AddrErr)
	jobId := executeJobCmd.MarkFlagRequired("jobId")
	utils.CheckError("JobId error : ", jobId)
	configPath := executeJobCmd.MarkFlagRequired("config")
	utils.CheckError("Path error : ", configPath)
	zenPath := executeJobCmd.MarkFlagRequired("zen-path")
	utils.CheckError("Pipeline Path error : ", zenPath)
}
