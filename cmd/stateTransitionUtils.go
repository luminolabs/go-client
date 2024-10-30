package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"lumino/core/types"
	pipeline_zen "lumino/pipeline-zen"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/rand"
)

func (*UtilsStruct) HandleStateTransition(ctx context.Context, client *ethclient.Client, config types.Configurations, account types.Account, state types.EpochState, epoch uint32, isAdmin bool, pipelinePath string) error {
	switch state {
	case types.EpochStateAssign:
		if !isAdmin {
			log.Debug("Not an admin node, skipping assignment state")
			return nil
		}
		return cmdUtils.HandleAssignState(ctx, client, config, account, epoch)
	case types.EpochStateUpdate:
		return cmdUtils.HandleUpdateState(ctx, client, config, account, epoch, pipelinePath)
	case types.EpochStateConfirm:
		return cmdUtils.HandleConfirmState(ctx, client, config, account, epoch, pipelinePath)
	default:
		log.WithField("state", state).Info("Waiting for next assign state")
		opts := protoUtils.GetOptions()
		return stateManagerUtils.WaitForNextState(client, &opts, types.EpochStateAssign)
	}
}

func (*UtilsStruct) HandleAssignState(ctx context.Context, client *ethclient.Client, config types.Configurations, account types.Account, epoch uint32) error {

	log.WithFields(logrus.Fields{
		"Current Task": "Executing Handle Assign State",
	}).Info("In Assign State")
	opts := protoUtils.GetOptions()
	config, err := cmdUtils.GetConfigData()
	// TODO: to be replaced by activeStaker or a better mechanism
	// numStakers, err := stakeManagerUtils.GetNumStakers(client, &opts)
	// if err != nil {
	// 	log.Error("Error in getting Num stakers: ", err)
	// 	return fmt.Errorf("failed to get Number of stakers jobs: %w", err)
	// }
	log.Debug("Num stakers : ", 3)

	activeStakers := [1]string{"0xC4481aa21AeAcAD3cCFe6252c6fe2f161A47A771"}

	// Get unassigned jobs and assign them
	// TODO: to be moved to jobsManagerUtils in struct Utils in future
	unassignedJobs, err := jobsManagerUtils.GetActiveJobs(client, &opts)
	if err != nil {
		return fmt.Errorf("failed to get unassigned jobs: %w", err)
	}
	if len(unassignedJobs) == 0 {
		log.Info("No active jobs to be assigned")
	} else {
		log.Debug("ActiveUnassignedJobs : ", unassignedJobs)
	}

	// TODO: make assignJob accept array input
	for _, jobId := range unassignedJobs {
		// Pick a random staker
		selectedStaker := activeStakers[rand.Intn(len(activeStakers))]

		if jobId.Cmp(big.NewInt(0)) == 1 {

			log.WithFields(logrus.Fields{
				"jobId":  jobId.String(),
				"staker": selectedStaker,
			}).Info("Assigning job")

			txnHash, err := cmdUtils.AssignJob(client, config, account, selectedStaker, jobId, 0)
			if err != nil {
				log.WithError(err).Error("Failed to assign job")
				continue
			}

			log.WithField("txHash", txnHash.Hex()).Info("Job assigned successfully")
		}
	}
	return nil
}

func (*UtilsStruct) HandleUpdateState(ctx context.Context, client *ethclient.Client, config types.Configurations, account types.Account, epoch uint32, pipelinePath string) error {
	// Check if already running a job
	stateMutex.RLock()
	isJobRunning := executionState.IsJobRunning
	stateMutex.RUnlock()

	log.WithFields(logrus.Fields{
		"Current Task": "Executing Handle Update State",
	}).Info("In Update State")

	opts := protoUtils.GetOptions()

	// Get job assigned to this staker
	jobId, err := jobsManagerUtils.GetJobForStaker(client, &opts, common.HexToAddress(account.Address))
	if err != nil {
		return fmt.Errorf("failed to get job for staker: %w", err)
	}

	if jobId.Cmp(big.NewInt(0)) == 0 {
		log.Debug("No job assigned")
		return nil
	}

	// Get job status
	status, err := jobsManagerUtils.GetJobStatus(client, &opts, jobId)
	if err != nil {
		return fmt.Errorf("failed to get job status: %w", err)
	}

	if isJobRunning && status == uint8(types.JobStatusRunning) {
		log.Info("Already running a job")
		return nil
	}

	if status != uint8(types.JobStatusQueued) {
		log.WithFields(logrus.Fields{
			"jobId":  jobId.String(),
			"status": status,
		}).Debug("Job not Queued yet")
		return nil
	}

	// Get job details
	jobDetails, err := jobsManagerUtils.GetJobDetails(client, &opts, jobId)
	if err != nil {
		return fmt.Errorf("failed to get job details: %w", err)
	}

	// Log the input for debugging
	log.WithField("rawJSON", jobDetails.JobDetailsInJSON).Debug("Received job details JSON")

	// Clean up the JSON string by removing escaped characters
	cleanJSON := cleanJSONString(jobDetails.JobDetailsInJSON)
	log.WithField("cleanJSON", cleanJSON).Debug("Cleaned JSON string")
	// First unmarshal the incoming JSON into a map to handle unknown structure
	var rawConfig map[string]interface{}
	if err := json.Unmarshal([]byte(cleanJSON), &rawConfig); err != nil {
		return fmt.Errorf("failed to parse job details: %w", err)
	}

	// Create the new config structure with the desired format
	jobConfig := types.JobConfig{
		JobConfigName: getString(rawConfig, "job_config_name"),
		JobID:         jobId.String(),
		DatasetID:     getString(rawConfig, "dataset_id"),
		BatchSize:     getString(rawConfig, "batch_size"),
		Shuffle:       getString(rawConfig, "shuffle"),
		NumEpochs:     getString(rawConfig, "num_epochs"),
		UseLora:       getString(rawConfig, "use_lora"),
		UseQlora:      getString(rawConfig, "use_qlora"),
		LearningRate:  getString(rawConfig, "lr", "1e-2"),
		OverrideEnv:   getString(rawConfig, "override_env", "prod"), // default to "prod"
		Seed:          getString(rawConfig, "seed", "42"),           // default to "42"
		NumGPUs:       getString(rawConfig, "num_gpus"),
		UserID:        jobDetails.Creator.String(),
	}

	// Create job directory in .lumino
	jobDir := filepath.Join("./.jobs", jobId.String())
	if err := os.MkdirAll(jobDir, 0755); err != nil {
		return fmt.Errorf("failed to create job directory: %w", err)
	}

	// Create job directory if it doesn't exist
	if err := os.MkdirAll(jobDir, 0755); err != nil {
		return fmt.Errorf("failed to create job directory: %w", err)
	}

	// Marshal the config with proper indentation
	configJson, err := json.MarshalIndent(jobConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal job config: %w", err)
	}

	configJson = append(configJson, '\n')

	// Write to file
	configPath := filepath.Join(jobDir, "config.json")
	if err := os.WriteFile(configPath, configJson, 0644); err != nil {
		return fmt.Errorf("failed to write job config: %w", err)
	}

	log.WithFields(logrus.Fields{
		"jobId":      jobId.String(),
		"configPath": configPath,
	}).Debug("Job configuration written to file")

	// Start job execution in goroutine
	// Update state
	stateMutex.Lock()
	executionState.CurrentJob = &types.JobExecution{
		JobID:     jobId,
		Status:    types.JobStatusRunning,
		StartTime: time.Now(),
		Executor:  account.Address,
	}
	executionState.IsJobRunning = true
	stateMutex.Unlock()

	// Update status to running
	txnHash, err := cmdUtils.UpdateJobStatus(client, config, account, jobId, types.JobStatusRunning, 0)
	if err != nil {
		log.WithError(err).Error("Failed to update job status to running")
		return fmt.Errorf("failed to update jobStatus: %w", err)
	}
	log.WithField("txHash", txnHash.Hex()).Info("Job status updated to Running")

	// Install dependencies
	// if err := pipeline_zen.InstallDeps(pipelinePath); err != nil {
	// 	log.WithError(err).Error("Failed to install dependencies")
	// 	cmdUtils.updateJobStatus(client, config, account, jobId)
	// 	return nil
	// }

	// Execute job with the config from .lumino directory
	output, err := pipeline_zen.RunTorchTuneWrapper(pipelinePath, configPath)
	if err != nil {
		log.WithError(err).Error("Job execution failed")
		stateMutex.Lock()
		executionState.CurrentJob = &types.JobExecution{
			JobID:    jobId,
			Status:   types.JobStatusFailed,
			Executor: account.Address,
		}
		stateMutex.Unlock()
		// cmdUtils.UpdateJobStatus(client, config, account, jobId, types.JobStatusFailed, 0)
		return nil
	}

	log.WithFields(logrus.Fields{
		"jobId":  jobId.String(),
		"output": output,
	}).Info("Job is running successfully")

	return nil
}

// getString safely extracts a string value from the map, with optional default value
func getString(m map[string]interface{}, key string, defaultValue ...string) string {
	if val, ok := m[key]; ok {
		// Convert the value to string regardless of its original type
		switch v := val.(type) {
		case string:
			return v
		case bool:
			return fmt.Sprintf("%v", v)
		case float64:
			// Handle integers vs decimals differently
			if v == float64(int(v)) {
				return fmt.Sprintf("%.0f", v)
			}
			return fmt.Sprintf("%g", v)
		default:
			return fmt.Sprintf("%v", v)
		}
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

// cleanJSONString properly formats the JSON string
func cleanJSONString(input string) string {
	// Handle the specific case of the unquoted job_config_name
	input = strings.Replace(input, "job_config_name:", `"job_config_name":`, 1)

	// Remove any whitespace between the opening brace and first key
	input = strings.TrimSpace(input)
	if strings.HasPrefix(input, "{") {
		input = "{" + strings.TrimSpace(input[1:])
	}

	return input
}

func (*UtilsStruct) HandleConfirmState(ctx context.Context, client *ethclient.Client, config types.Configurations, account types.Account, epoch uint32, pipelinePath string) error {

	log.WithFields(logrus.Fields{
		"Current Task": "Executing Handle Confirm State",
	}).Info("In Confirm State")

	stateMutex.RLock()
	currentJob := executionState.CurrentJob
	// isJobRunning := executionState.IsJobRunning
	stateMutex.RUnlock()

	if currentJob == nil {
		log.Debug("No current job found")
		return nil
	}

	jobId := currentJob.JobID
	currentStatus := currentJob.Status

	log.WithFields(logrus.Fields{
		"jobId":  jobId.String(),
		"status": currentStatus,
		"epoch":  epoch,
	}).Debug("Current job state")

	opts := protoUtils.GetOptions()
	// Get job details
	jobDetails, err := jobsManagerUtils.GetJobDetails(client, &opts, jobId)
	if err != nil {
		return fmt.Errorf("failed to get job details: %w", err)
	}
	resultsPath := ".results/" + jobDetails.Creator.String() + currentJob.JobID.String()
	zenPath := filepath.Join(pipelinePath, resultsPath)
	startedFile := filepath.Join(zenPath, ".started")
	finishedFile := filepath.Join(zenPath, ".finished")

	log.WithFields(logrus.Fields{
		"zenPath":      zenPath,
		"startedFile":  startedFile,
		"finishedFile": finishedFile,
	}).Debug("Checking job status files")

	// Handle failed status first
	if currentStatus == types.JobStatusFailed {
		log.WithField("jobId", jobId.String()).Info("Updating failed job status")
		txnHash, err := cmdUtils.UpdateJobStatus(client, config, account, jobId, types.JobStatusFailed, 0)
		if err != nil {
			return fmt.Errorf("failed to update job status to failed: %w", err)
		}
		log.WithField("txHash", txnHash.Hex()).Info("Job status updated to Failed")

		// Clear job state
		stateMutex.Lock()
		executionState.CurrentJob = nil
		executionState.IsJobRunning = false
		stateMutex.Unlock()

		return nil
	}

	// Check if .started file exists
	startedExists := false
	if _, err := os.Stat(startedFile); err == nil {
		startedExists = true
	}

	// Check if .finished file exists
	finishedExists := false
	if _, err := os.Stat(finishedFile); err == nil {
		finishedExists = true
	}

	if !startedExists {
		log.WithField("jobId", jobId.String()).Info("No job is running")
		return nil
	}

	if startedExists && !finishedExists {
		log.WithField("jobId", jobId.String()).Info("Job is still running")
		return nil
	}

	if startedExists && finishedExists {

		log.WithFields(logrus.Fields{
			"jobId": jobId.String(),
		}).Info("Job has concluded successfully")

		log.WithField("jobId", jobId.String()).Info("Updating status")

		// Job completed, update status
		txnHash, err := cmdUtils.UpdateJobStatus(client, config, account, jobId, types.JobStatusCompleted, 0)
		if err != nil {
			return fmt.Errorf("failed to update job status to completed: %w", err)
		}
		log.WithField("txHash", txnHash.Hex()).Info("Job status updated to Completed")

		// Clear job state
		stateMutex.Lock()
		executionState.CurrentJob = nil
		executionState.IsJobRunning = false
		stateMutex.Unlock()
	}

	return nil
}
