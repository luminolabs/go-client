package cmd

import (
	"errors"
	"lumino/core"
	"lumino/core/types"
	"lumino/logger"
	"lumino/pkg/bindings"
	"lumino/utils"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"

	"github.com/spf13/cobra"
)

var assignJobCmd = &cobra.Command{
	Use:   "assignJob",
	Short: "Assign a job to a compute provider",
	Long: `[ADMIN ONLY] Assign a new job by providing a job Id and an assignee.

Example:
  ./lumino assignJob -a 0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771 --assignee 0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771 --jobId 1`,
	Run: initialiseAssignJob,
}

// This function initialises the ExecuteCreateJob function
func initialiseAssignJob(cmd *cobra.Command, args []string) {
	cmdUtils.ExecuteAssignJob(cmd.Flags())
}

// ExecuteAssignJob is the main entry point for the job assignment process.
// It manages the workflow of job assignment:
// 1. Loads and validates configuration
// 2. Sets up blockchain connection and logging
// 3. Validates assignee address and job ID
// 4. Executes the job assignment
// Returns early if any validation step fails.
func (*UtilsStruct) ExecuteAssignJob(flagSet *pflag.FlagSet) {
	config, err := cmdUtils.GetConfigData()
	utils.CheckError("Error in getting config: ", err)
	log.Debugf("RunAssignJob: Config: %+v", config)

	client := protoUtils.ConnectToEthClient(config.Provider)

	address, err := flagSetUtils.GetStringAddress(flagSet)
	utils.CheckError("Error in getting address: ", err)

	logger.SetLoggerParameters(client, address)
	log.Debug("Checking to assign log file...")
	protoUtils.AssignLogFile(flagSet)

	log.Debug("Getting password...")
	password := protoUtils.AssignPassword(flagSet)

	assigneeAddress, err := flagSet.GetString("assignee")
	utils.CheckError("Error in getting assignee address: ", err)

	if !common.IsHexAddress(assigneeAddress) {
		log.WithField("assignee", assigneeAddress).Fatal("Invalid assignee address format")
	}

	jobIdStr, err := flagSet.GetString("jobId")
	utils.CheckError("Error in getting jobId: ", err)

	jobId, ok := new(big.Int).SetString(jobIdStr, 10)
	if !ok {
		log.Fatal("Invalid JobId format", errors.New("Failed to parse job ID string"))
	}

	account := types.Account{
		Address:  address,
		Password: password,
	}

	buffer := uint8(0)
	// Assign the job
	log.Info("Assigning job...")
	txnHash, err := cmdUtils.AssignJob(client, config, account, assigneeAddress, jobId, buffer)
	utils.CheckError("Error assigning job: ", err)

	log.WithFields(logrus.Fields{
		"txHash":   txnHash.Hex(),
		"jobId":    jobId.String(),
		"assignee": assigneeAddress,
	}).Info("Job assigned successfully")
}

// AssignJob assigns an existing job to a compute provider in the network. This function:
// 1. Validates the job and assignee addresses
// 2. Constructs and submits the assignment transaction
// 3. Monitors transaction confirmation
// Returns the transaction hash once the assignment is confirmed.
func (u *UtilsStruct) AssignJob(client *ethclient.Client, config types.Configurations, account types.Account, assigneeAddress string, jobId *big.Int, buffer uint8) (common.Hash, error) {
	if client == nil {
		log.Error("Client is nil")
		return common.Hash{}, errors.New("client is nil")
	}

	if !common.IsHexAddress(assigneeAddress) {
		log.WithField("assignee", assigneeAddress).Error("Invalid assignee address")
		return common.Hash{}, errors.New("invalid assignee address")
	}

	if jobId == nil {
		log.Error("JobId is nil")
		return common.Hash{}, errors.New("jobId is nil")
	}

	log.WithFields(logrus.Fields{
		"owner":    account.Address,
		"assignee": assigneeAddress,
		"jobId":    jobId.String(),
	}).Debug("Assigning job")

	// TODO: Implement GetJobStatus and wait for the assignState
	// status, err := jobsManagerUtils.GetJobStatus(client, jobId)
	// if err != nil {
	//     log.WithError(err).WithField("jobId", jobId.String()).Error("Failed to get job status")
	//     return common.Hash{}, err
	// }

	// Check if job status allows assignment
	// if status != types.JobStatusNew {
	//     log.WithFields(logrus.Fields{
	//         "jobId": jobId.String(),
	//         "status": status,
	//     }).Error("Job is not in assignable state")
	//     return common.Hash{}, errors.New("job must be in NEW status to be assigned")
	// }

	txnArgs := types.TransactionOptions{
		Client:          client,
		AccountAddress:  account.Address,
		Password:        account.Password,
		ChainId:         core.ChainID,
		Config:          config,
		ContractAddress: core.JobManagerAddress,
		MethodName:      "assignJob",
		Parameters:      []interface{}{jobId, common.HexToAddress(assigneeAddress), buffer},
		ABI:             bindings.JobManagerABI,
	}

	if jobsManagerUtils == nil {
		log.Error("JobManagerUtils is nil")
		return common.Hash{}, errors.New("jobManagerUtils is nil")
	}

	log.WithFields(logrus.Fields{
		"jobId":    jobId.String(),
		"assignee": assigneeAddress,
	}).Debug("Executing assignJob transaction")

	txnOpts := protoUtils.GetTransactionOpts(txnArgs)

	txn, err := jobsManagerUtils.AssignJob(client, txnOpts, jobId, common.HexToAddress(assigneeAddress), buffer)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"jobId":    jobId.String(),
			"assignee": assigneeAddress,
		}).Error("Failed to assign job")
		return common.Hash{}, err
	}

	if txn == nil {
		log.Error("Transaction is nil")
		return common.Hash{}, errors.New("transaction is nil")
	}

	txnHash := transactionUtils.Hash(txn)
	log.WithFields(logrus.Fields{
		"txHash":   txnHash.Hex(),
		"jobId":    jobId.String(),
		"assignee": assigneeAddress,
	}).Info("Job assignment transaction submitted")

	err = protoUtils.WaitForBlockCompletion(txnArgs.Client, txnHash.Hex())
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{
			"jobId":    jobId.String(),
			"assignee": assigneeAddress,
		}).Error("Failed to wait for block completion")
		return common.Hash{}, err
	}

	return txnHash, nil
}

// Configures the job assignment command including required flags for account address,
// assignee address and job ID. Sets up help text and usage information.
func init() {
	rootCmd.AddCommand(assignJobCmd)

	var (
		Account  string
		Password string
		Assignee string
		JobId    string
	)

	assignJobCmd.Flags().StringVarP(&Account, "address", "a", "", "address of the job owner")
	assignJobCmd.Flags().StringVarP(&Password, "password", "", "", "password path of job owner")
	assignJobCmd.Flags().StringVarP(&Assignee, "assignee", "", "", "address of the compute provider to assign")
	assignJobCmd.Flags().StringVarP(&JobId, "jobId", "", "", "ID of the job to assign")

	// Check errors when marking flags as required
	if err := assignJobCmd.MarkFlagRequired("address"); err != nil {
		log.WithError(err).Fatal("Error marking 'address' flag as required")
	}
	if err := assignJobCmd.MarkFlagRequired("assignee"); err != nil {
		log.WithError(err).Fatal("Error marking 'assignee' flag as required")
	}
	if err := assignJobCmd.MarkFlagRequired("jobId"); err != nil {
		log.WithError(err).Fatal("Error marking 'jobId' flag as required")
	}
}
