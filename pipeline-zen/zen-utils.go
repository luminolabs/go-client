package pipeline_zen

import (
	"encoding/json"
	"fmt"
	"lumino/logger"
	"os"
	"os/exec"
	"path/filepath"
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
}

// InstallDeps runs the install-deps.sh script to install Python dependencies
func InstallDeps() (string, error) {
	cmd := exec.Command("bash", "-c", "./scripts/install-deps.sh")
	cmd.Dir = filepath.Join(".", "pipeline-zen") // Set working directory to pipeline-zen folder

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Error installing dependencies: ", err)
		return "", err
	}

	logger.Info("Dependencies installed successfully")
	return string(output), nil
}

// RunTorchTuneWrapper runs the torchtunewrapper workflow with the provided configuration
func RunTorchTuneWrapper(configFile string) (string, error) {
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

	// Build the bash command string
	command := fmt.Sprintf("./scripts/runners/celery-wf.sh torchtunewrapper --job_config_name %s --job_id %s --dataset_id %s --batch_size %s --shuffle %s --num_epochs %s --use_lora %s --use_qlora %s --lr %s --seed %s --num_gpus %s",
		config.JobConfigName, config.JobID, config.DatasetID, config.BatchSize, config.Shuffle, config.NumEpochs, config.UseLora, config.UseQlora, config.LearningRate, config.Seed, config.NumGpus)

	logger.Info("Running command: ", command)

	// Execute the command
	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = filepath.Join(".", "pipeline-zen") // Set working directory to pipeline-zen folder

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Error running torchtunewrapper workflow: ", err)
		return "", err
	}

	logger.Info("Workflow executed successfully")
	return string(output), nil
}
