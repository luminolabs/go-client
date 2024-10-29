// Package cmd provides all functions related to command line
package cmd

import (
	"errors"
	"lumino/core"
	"lumino/core/types"
	"lumino/logger"
	pipeline_zen "lumino/pipeline-zen"
	"lumino/pkg/bindings"
	"lumino/utils"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
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
	password := protoUtils.AssignPassword(flagSet)

	jobIdStr, err := flagSet.GetString("jobId")
	utils.CheckError("Error in getting jobId: ", err)

	jobId, ok := new(big.Int).SetString(jobIdStr, 10)
	if !ok {
		log.Fatal("Invalid JobId format", errors.New("Failed to parse job ID string"))
	}
	configPath, err := flagSet.GetString("config")
	utils.CheckError("Error in getting config path: ", err)

	pipelinePath, err := flagSet.GetString("zen-path")
	utils.CheckError("Error in getting pipeline path: ", err)

	// Install dependencies with live logging
	//log.Info("Starting dependency installation...")
	//err = pipeline_zen.InstallDeps(pipelinePath)
	//if err != nil {
	//	log.WithError(err).Fatal("Failed to install dependencies")
	//}
	//log.Info("Dependencies installed successfully")

	// Hardcoded, to be changed in future
	status := types.JobStatusQueued
	buffer := 0

	// Update job status to Queued
	log.Info("Updating job status to Queued...")
	jobUpdateTxn, err := cmdUtils.UpdateJobStatus(client, config, types.Account{
		Address:  address,
		Password: password,
	}, jobId, status, uint8(buffer))
	if err != nil {
		log.WithError(err).Fatal("Failed to update job status to Queued")
	}
	log.WithField("txHash", jobUpdateTxn.Hex()).Info("Job status updated to Queued")

	// Run the TorchTuneWrapper
	log.Info("Running TorchTuneWrapper...")
	go func() {
		output, err := pipeline_zen.RunTorchTuneWrapper(pipelinePath, configPath)
		if err != nil {
			log.WithError(err).Error("Error running TorchTuneWrapper")
			cmdUtils.UpdateJobStatus(client, config, types.Account{
				Address:  address,
				Password: password,
			}, jobId, types.JobStatusFailed, uint8(buffer))
			return
		}
		log.Info("Updating job status to Running...")
		runningTxnHash, err := cmdUtils.UpdateJobStatus(client, config, types.Account{
			Address:  address,
			Password: password,
		}, jobId, types.JobStatusRunning, uint8(buffer))
		log.WithField("txHash", runningTxnHash.Hex()).Info("Job status updated to Running")

		log.Debug("TorchTuneWrapper output: ", output)
		log.Info("Job execution initiated. Monitor logs for progress.")
	}()

	// Update job status to Running
	completedJobUpdateTxn, err := cmdUtils.UpdateJobStatus(client, config, types.Account{
		Address:  address,
		Password: password,
	}, jobId, types.JobStatusCompleted, uint8(buffer))
	log.WithField("txHash", completedJobUpdateTxn.Hex()).Info("Job status updated to Completed")

}

func (u *UtilsStruct) UpdateJobStatus(client *ethclient.Client, config types.Configurations, account types.Account, jobId *big.Int, status types.JobStatus, buffer uint8) (common.Hash, error) {
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
		JobId      string
		Account    string
		Password   string
		ConfigPath string
		ZenPath    string
	)

	executeJobCmd.Flags().StringVarP(&JobId, "jobId", "", string(0), "job id")
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
