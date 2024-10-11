package pipeline_zen

import (
	"fmt"
	"os/exec"
)

// RunTorchTuneWrapper runs the torchtunewrapper workflow with the provided arguments
func RunTorchTuneWrapper(jobConfigName, jobID, datasetID, batchSize, shuffle, numEpochs, useLora, useQlora, learningRate, seed, numGpus string) (string, error) {
	// Build the bash command string
	command := fmt.Sprintf("./scripts/runners/single-wf.sh torchtunewrapper --job_config_name %s --job_id %s --dataset_id %s --batch_size %s --shuffle %s --num_epochs %s --use_lora %s --use_qlora %s --lr %s --seed %s --num_gpus %s",
		jobConfigName, jobID, datasetID, batchSize, shuffle, numEpochs, useLora, useQlora, learningRate, seed, numGpus)

	// Execute the command
	cmd := exec.Command("bash", "-c", command)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Error running torchtunewrapper workflow: %v", err)
	}

	return string(output), nil
}
