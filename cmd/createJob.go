// Package cmd provides all functions related to command line
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

var createJobCmd = &cobra.Command{
	Use:   "createJob",
	Short: "Create a new job",
	Long: `Create a new job by providing a job configuration file and job fee.

Example:
  ./lumino createJob -a 0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771 --config /path/to/config.json --jobFee 1000000000000000000`,
	Run: initialiseCreateJob,
}

// This function initialises the ExecuteCreateJob function
func initialiseCreateJob(cmd *cobra.Command, args []string) {
	cmdUtils.ExecuteCreateJob(cmd.Flags())
}

// This function sets the flags appropriately and executes the CreateJob function
func (*UtilsStruct) ExecuteCreateJob(flagSet *pflag.FlagSet) {
	config, err := cmdUtils.GetConfigData()
	utils.CheckError("Error in getting config: ", err)
	log.Debugf("RunCreateJob: Config: %+v", config)

	client := protoUtils.ConnectToEthClient(config.Provider)

	address, err := flagSetUtils.GetStringAddress(flagSet)
	utils.CheckError("Error in getting address: ", err)

	logger.SetLoggerParameters(client, address)
	log.Debug("Checking to assign log file...")
	protoUtils.AssignLogFile(flagSet)

	log.Debug("Getting password...")
	password := protoUtils.AssignPassword(flagSet)

	configPath, err := flagSet.GetString("config")
	utils.CheckError("Error in getting config path: ", err)

	jobFeeStr, err := flagSet.GetString("jobFee")
	utils.CheckError("Error in getting job fee: ", err)

	jobFee, ok := new(big.Int).SetString(jobFeeStr, 10)
	if !ok {
		utils.CheckError("Error converting job fee to big.Int: ", errors.New("invalid job fee"))
	}

	// Read and parse the job configuration file
	jobDetailsJSON, err := osUtils.ReadFile(configPath)
	utils.CheckError("Error reading job configuration file: ", err)

	// Create the job
	log.Info("Creating job...")
	txnHash, err := cmdUtils.CreateJob(client, config, types.Account{
		Address:  address,
		Password: password,
	}, string(jobDetailsJSON), jobFee)
	utils.CheckError("Error creating job: ", err)

	log.Info("Job created successfully. Transaction Hash: ", txnHash.Hex())
}

func (u *UtilsStruct) CreateJob(client *ethclient.Client, config types.Configurations, account types.Account, jobDetailsJSON string, jobFee *big.Int) (common.Hash, error) {
	if client == nil {
		log.Error("Client is nil")
		return common.Hash{}, errors.New("client is nil")
	}
	if jobFee == nil {
		log.Error("JobFee is nil")
		return common.Hash{}, errors.New("jobFee is nil")
	}

	log.WithFields(logrus.Fields{
		"address": account.Address,
		"jobFee":  jobFee.String(),
	}).Debug("Creating job")

	txnArgs := types.TransactionOptions{
		Client:          client,
		AccountAddress:  account.Address,
		Password:        account.Password,
		ChainId:         core.ChainID,
		Config:          config,
		ContractAddress: core.JobManagerAddress,
		MethodName:      "createJob",
		Parameters:      []interface{}{jobDetailsJSON},
		ABI:             bindings.JobManagerABI,
		EtherValue:      jobFee,
	}

	if jobsManagerUtils == nil {
		log.Error("JobManagerUtils is nil")
		return common.Hash{}, errors.New("jobManagerUtils is nil")
	}

	log.WithFields(logrus.Fields{
		"jobFee": jobFee.String(),
	}).Debug("Executing createJob transaction")

	txnOpts := protoUtils.GetTransactionOpts(txnArgs)

	txn, err := jobsManagerUtils.CreateJob(txnArgs.Client, txnOpts, jobDetailsJSON)
	if err != nil {
		log.WithError(err).Error("Failed to create job")
		return common.Hash{}, err
	}

	if txn == nil {
		log.Error("Transaction is nil")
		return common.Hash{}, errors.New("transaction is nil")
	}

	txnHash := transactionUtils.Hash(txn)
	log.WithField("txHash", txnHash.Hex()).Info("Job creation transaction submitted")

	err = protoUtils.WaitForBlockCompletion(txnArgs.Client, txnHash.Hex())
	if err != nil {
		log.WithError(err).Error("Failed to wait for block completion")
		return common.Hash{}, err
	}

	return txnHash, nil
}

func init() {
	rootCmd.AddCommand(createJobCmd)

	var (
		Account    string
		Password   string
		ConfigPath string
		JobFee     string
	)

	createJobCmd.Flags().StringVarP(&Account, "address", "a", "", "address of the job creator")
	createJobCmd.Flags().StringVarP(&Password, "password", "", "", "password path of job creator to protect the keystore")
	createJobCmd.Flags().StringVarP(&ConfigPath, "config", "c", "", "path to the job configuration file")
	createJobCmd.Flags().StringVarP(&JobFee, "jobFee", "f", "", "job fee in wei")

	AddrErr := createJobCmd.MarkFlagRequired("address")
	utils.CheckError("Address error: ", AddrErr)
	configPath := createJobCmd.MarkFlagRequired("config")
	utils.CheckError("Path error : ", configPath)
	jobFee := createJobCmd.MarkFlagRequired("jobFee")
	utils.CheckError("JobFee error : ", jobFee)
}
