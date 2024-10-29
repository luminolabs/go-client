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

func (u *UtilsStruct) HandleAssignState(ctx context.Context, client *ethclient.Client, config types.Configurations, account types.Account, epoch uint32) error {

	opts := protoUtils.GetOptions()
	numStakers := stakeManagerUtils.GetNumStakers(client, opts)
	// Get unassigned jobs and assign them
	// TODO: to be moved to jobsManagerUtils in struct Utils in future
	unassignedJobs, err := jobsManagerUtils.GetActiveJobs(client)
	if err != nil {
		return fmt.Errorf("failed to get unassigned jobs: %w", err)
	}

	for _, jobId := range unassignedJobs {
		// Pick a random staker
		selectedStaker := activeStakers[rand.Intn(len(activeStakers))]

		log.WithFields(logrus.Fields{
			"jobId":  jobId.String(),
			"staker": selectedStaker.Hex(),
		}).Info("Assigning job")

		txnHash, err := u.AssignJob(client, config, account, selectedStaker.Hex(), jobId)
		if err != nil {
			log.WithError(err).Error("Failed to assign job")
			continue
		}

		log.WithField("txHash", txnHash.Hex()).Info("Job assigned successfully")
	}

	return nil
}

func (u *UtilsStruct) HandleUpdateState(ctx context.Context, client *ethclient.Client, config types.Configurations, account types.Account, epoch uint32, pipelinePath string) error {
	// Check if already running a job
	stateMutex.RLock()
	isJobRunning := executionState.IsJobRunning
	stateMutex.RUnlock()

	if isJobRunning {
		log.Info("Already running a job")
		return nil
	}

	// Get job assigned to this staker
	jobId, err := jobsManagerUtils.GetJobForStaker(client, common.HexToAddress(account.Address))
	if err != nil {
		return fmt.Errorf("failed to get job for staker: %w", err)
	}

	if jobId.Cmp(big.NewInt(0)) == 0 {
		log.Debug("No job assigned")
		return nil
	}

	// Get job status
	status, err := jobsManagerUtils.GetJobStatus(client, jobId)
	if err != nil {
		return fmt.Errorf("failed to get job status: %w", err)
	}

	if status != types.JobStatusAssigned {
		log.WithFields(logrus.Fields{
			"jobId":  jobId.String(),
			"status": status.String(),
		}).Debug("Job not in assigned state")
		return nil
	}

	// Get job details
	jobDetails, err := jobsManagerUtils.GetJobDetails(client, jobId)
	if err != nil {
		return fmt.Errorf("failed to get job details: %w", err)
	}

	// Parse job configuration
	var jobConfig types.JobConfig
	if err := json.Unmarshal([]byte(jobDetails), &jobConfig); err != nil {
		return fmt.Errorf("failed to parse job config: %w", err)
	}

	// Update status to queued
	txnHash, err := u.UpdateJobStatus(client, config, account, jobId, types.JobStatusQueued)
	if err != nil {
		return fmt.Errorf("failed to update job status to queued: %w", err)
	}
	log.WithField("txHash", txnHash.Hex()).Info("Job status updated to Queued")

	// Start job execution in goroutine
	go u.executeJobProcess(ctx, client, config, account, jobId, jobConfig, pipelinePath)

	return nil
}

func (u *UtilsStruct) executeJobProcess(ctx context.Context, client *ethclient.Client, config types.Configurations, account types.Account, jobId *big.Int, jobConfig types.JobConfig, pipelinePath string) {
	// Update state
	stateMutex.Lock()
	executionState.CurrentJob = &types.JobExecution{
		JobID:      jobId,
		Status:     types.JobStatusRunning,
		StartTime:  time.Now(),
		Executor:   account.Address,
		PipelineID: jobConfig.PipelineID,
	}
	executionState.IsJobRunning = true
	stateMutex.Unlock()

	// Update status to running
	txnHash, err := u.UpdateJobStatus(client, config, account, jobId, types.JobStatusRunning)
	if err != nil {
		log.WithError(err).Error("Failed to update job status to running")
		return
	}
	log.WithField("txHash", txnHash.Hex()).Info("Job status updated to Running")

	// Write job config to file
	configPath := filepath.Join(pipelinePath, fmt.Sprintf("job-%s.json", jobId.String()))
	configJson, err := json.MarshalIndent(jobConfig, "", "  ")
	if err != nil {
		log.WithError(err).Error("Failed to marshal job config")
		return
	}

	if err := os.WriteFile(configPath, configJson, 0644); err != nil {
		log.WithError(err).Error("Failed to write job config")
		return
	}

	// Install dependencies
	if err := pipeline_zen.InstallDeps(pipelinePath); err != nil {
		log.WithError(err).Error("Failed to install dependencies")
		u.updateJobStatusToFailed(client, config, account, jobId)
		return
	}

	// Execute job
	output, err := pipeline_zen.RunTorchTuneWrapper(configPath)
	if err != nil {
		log.WithError(err).Error("Job execution failed")
		u.updateJobStatusToFailed(client, config, account, jobId)
		return
	}

	log.WithField("output", output).Info("Job completed successfully")
}

func (u *UtilsStruct) HandleConfirmState(ctx context.Context, client *ethclient.Client, config types.Configurations, account types.Account, epoch uint32) error {
	stateMutex.RLock()
	currentJob := executionState.CurrentJob
	isJobRunning := executionState.IsJobRunning
	stateMutex.RUnlock()

	if currentJob == nil {
		return nil
	}

	jobId := currentJob.JobID
	status, err := jobsManagerUtils.GetJobStatus(client, jobId)
	if err != nil {
		return fmt.Errorf("failed to get job status: %w", err)
	}

	if status == types.JobStatusRunning && !isJobRunning {
		// Job completed, update status
		txnHash, err := u.UpdateJobStatus(client, config, account, jobId, types.JobStatusCompleted)
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

func (u *UtilsStruct) updateJobStatusToFailed(client *ethclient.Client, config types.Configurations, account types.Account, jobId *big.Int) {
	txnHash, err := u.UpdateJobStatus(client, config, account, jobId, types.JobStatusFailed)
	if err != nil {
		log.WithError(err).Error("Failed to update job status to failed")
		return
	}
	log.WithField("txHash", txnHash.Hex()).Info("Job status updated to Failed")

	stateMutex.Lock()
	executionState.CurrentJob = nil
	executionState.IsJobRunning = false
	stateMutex.Unlock()
}
