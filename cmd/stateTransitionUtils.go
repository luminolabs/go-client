package cmd

import (
	"context"
	"fmt"
	"lumino/core/types"
	pipeline_zen "lumino/pipeline-zen"
	"math/big"
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
		return cmdUtils.HandleConfirmState(ctx, client, config, account, epoch)
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
	// jobDetails, err := jobsManagerUtils.GetJobDetails(client, &opts, jobId)
	// if err != nil {
	// 	return fmt.Errorf("failed to get job details: %w", err)
	// }

	// cleanJSON, err := strconv.Unquote(`"` + jobDetails.JobDetailsInJSON + `"`)
	// if err != nil {
	// 	log.Fatalf("Error unescaping JSON string: %v", err)
	// }

	// log.WithField("JobDetails in JSON", cleanJSON).Debug(" job config JSON")

	// // Parse job configuration from jobDetailsInJSON field
	// var jobConfigMap map[string]interface{}
	// if err := json.Unmarshal([]byte(cleanJSON), &jobConfigMap); err != nil {
	// 	return fmt.Errorf("failed to parse job config: %w", err)
	// }

	// // Create job directory in .lumino
	// jobDir := filepath.Join("./.jobs", jobId.String())
	// if err := os.MkdirAll(jobDir, 0755); err != nil {
	// 	return fmt.Errorf("failed to create job directory: %w", err)
	// }

	// // Write job config to file
	// configPath := filepath.Join(jobDir, "jobConfig.json")
	// configJson, err := json.MarshalIndent(jobConfigMap, "", "  ")
	// if err != nil {
	// 	return fmt.Errorf("failed to marshal job config: %w", err)
	// }

	// if err := os.WriteFile(configPath, configJson, 0644); err != nil {
	// 	return fmt.Errorf("failed to write job config: %w", err)
	// }

	// log.WithFields(logrus.Fields{
	// 	"jobId":      jobId.String(),
	// 	"configPath": configPath,
	// }).Debug("Job configuration written to file")

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

	configPathString := "/Users/shyampatel/.lumino/config.json"

	// Execute job with the config from .lumino directory
	output, err := pipeline_zen.RunTorchTuneWrapper(pipelinePath, configPathString)
	if err != nil {
		log.WithError(err).Error("Job execution failed")
		cmdUtils.UpdateJobStatus(client, config, account, jobId, types.JobStatusFailed, 0)
		return nil
	}

	log.WithFields(logrus.Fields{
		"jobId":  jobId.String(),
		"output": output,
	}).Info("Job completed successfully")

	// log.WithField("output", output).Info("Job completed successfully")

	return nil
}

func (*UtilsStruct) HandleConfirmState(ctx context.Context, client *ethclient.Client, config types.Configurations, account types.Account, epoch uint32) error {
	stateMutex.RLock()
	currentJob := executionState.CurrentJob
	isJobRunning := executionState.IsJobRunning
	stateMutex.RUnlock()

	log.WithFields(logrus.Fields{
		"Current Task": "Executing Handle Confirm State",
	}).Info("In Confirm State")

	if currentJob == nil {
		return nil
	}

	opts := protoUtils.GetOptions()

	jobId := currentJob.JobID
	status, err := jobsManagerUtils.GetJobStatus(client, &opts, jobId)
	if err != nil {
		return fmt.Errorf("failed to get job status: %w", err)
	}

	if status == uint8(types.JobStatusRunning) && !isJobRunning {
		// Job completed, update status
		txnHash, err := cmdUtils.UpdateJobStatus(client, config, account, jobId, types.JobStatusCompleted, 0)
		if err != nil {
			return fmt.Errorf("failed to update job status to completed: %w", err)
		}
		log.WithField("txHash", txnHash.Hex()).Info("Job status updated to Completed")

		// Clear current job
		stateMutex.Lock()
		executionState.CurrentJob = nil
		executionState.IsJobRunning = false
		stateMutex.Unlock()
	}

	return nil
}
