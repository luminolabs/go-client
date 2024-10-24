package pipeline_zen

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"lumino/logger"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

// TorchTuneWrapperConfig defines the structure of the config file for the workflow
type TorchTuneWrapperConfig struct {
	JobConfigName string `json:"job_config_name"`
	JobID         string `json:"job_id"`
	DatasetID     string `json:"dataset_id"`
	BatchSize     string `json:"batch_size"`
	Shuffle       string `json:"shuffle"`
	NumEpochs     string `json:"num_epochs"`
	UseLora       string `json:"use_lora"`
	UseQlora      string `json:"use_qlora"`
	LearningRate  string `json:"lr"`
	Seed          string `json:"seed"`
	NumGpus       string `json:"num_gpus"`
	UserId        string `json:"user_id"`
}

// log is the package-level logger instance
var log = logger.NewLogger()

// streamOutput reads from a pipe and logs each line
func streamOutput(reader io.ReadCloser, prefix string, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		// Use structured logging but with the raw output from the script
		log.WithField("source", prefix).Info(line)
	}
}

// InstallDeps runs the install-deps.sh script to install Python dependencies
func InstallDeps(pipelineZenPath string) error {
	// Verify pipeline-zen directory exists
	if _, err := os.Stat(pipelineZenPath); os.IsNotExist(err) {
		log.WithError(err).Error("Pipeline-zen directory not found")
		return fmt.Errorf("pipeline-zen directory not found at %s: %w", pipelineZenPath, err)
	}

	// Construct path to install script
	scriptPath := filepath.Join(pipelineZenPath, "scripts", "install-deps.sh")

	// Verify script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		log.WithError(err).Error("Install script not found")
		return fmt.Errorf("install script not found at %s: %w", scriptPath, err)
	}

	// Create the command
	cmd := exec.Command("bash", scriptPath)
	cmd.Dir = pipelineZenPath
	cmd.Env = append(os.Environ())

	// Get pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.WithError(err).Error("Failed to create stdout pipe")
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.WithError(err).Error("Failed to create stderr pipe")
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	log.WithField("scriptPath", scriptPath).Debug("Starting installation script")
	if err := cmd.Start(); err != nil {
		log.WithError(err).Error("Failed to start installation script")
		return fmt.Errorf("failed to start installation script: %w", err)
	}

	// Use WaitGroup to ensure we capture all output
	var wg sync.WaitGroup
	wg.Add(2)

	// Stream stdout and stderr
	go streamOutput(stdout, "install-stdout", &wg)
	go streamOutput(stderr, "install-stderr", &wg)

	// Wait for command to complete and output to be processed
	err = cmd.Wait()
	wg.Wait()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			log.WithError(err).WithField("exitCode", exitErr.ExitCode()).Error("Installation script failed")
			return fmt.Errorf("installation script failed with exit code %d: %w", exitErr.ExitCode(), err)
		}
		log.WithError(err).Error("Installation failed")
		return fmt.Errorf("installation failed: %w", err)
	}

	log.Info("Installation completed successfully")
	return nil
}

// RunTorchTuneWrapper runs the torchtunewrapper workflow with the provided configuration
func RunTorchTuneWrapper(pipelineZenPath string, configFile string) (string, error) {
	// Read the config file
	configData, err := os.ReadFile(configFile)
	if err != nil {
		logger.Error("Error reading config file: ", err)
		return "", err
	}

	// Unmarshal the JSON data into the TorchTuneWrapperConfig struct
	var config TorchTuneWrapperConfig
	err = json.Unmarshal(configData, &config)
	if err != nil {
		logger.Error("Error parsing JSON config: ", err)
		return "", err
	}

	// Validate required fields
	if err := config.Validate(); err != nil {
		logger.Error("Invalid config: ", err)
		return "", err
	}

	scriptPath := filepath.Join(pipelineZenPath, "scripts", "runners", "celery-wf.sh")
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		logger.Error("Script not found at path: ", scriptPath)
		return "", fmt.Errorf("script not found")
	}

	// Build the bash command string without the 'cd' and pipe
	// command := fmt.Sprintf("./scripts/runners/celery-wf.sh torchtunewrapper --job_config_name %s --job_id %s --dataset_id %s --batch_size %s --shuffle %s --num_epochs %s --use_lora %s --use_qlora %s --lr %s --seed %s --num_gpus %s",
	// 	config.JobConfigName, config.JobID, config.DatasetID, config.BatchSize, config.Shuffle, config.NumEpochs, config.UseLora, config.UseQlora, config.LearningRate, config.Seed, config.NumGpus)
	command := fmt.Sprintf("%s torchtunewrapper --job_config_name %s --job_id %s --dataset_id %s --batch_size %s --shuffle %s --num_epochs %s --use_lora %s --use_qlora %s --lr %s --seed %s --num_gpus %s --user_id %s",
		scriptPath, config.JobConfigName, config.JobID, config.DatasetID, config.BatchSize, config.Shuffle, config.NumEpochs, config.UseLora, config.UseQlora, config.LearningRate, config.Seed, config.NumGpus, config.UserId)

	logger.Info("Running command: ", command)

	// Set working directory to 'pipeline-zen' folder
	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = pipelineZenPath

	// Execute the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Error running torchtunewrapper workflow: ", err)
		return "", err
	}

	logger.Info("Command output: ", string(output))
	return string(output), nil
}
