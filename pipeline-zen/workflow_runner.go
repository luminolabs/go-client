package pipeline_zen

import (
	"fmt"
	"lumino/logger"
)

// Validate ensures all required fields are present and non-empty, with logging for each missing field
func (c *TorchTuneWrapperConfig) Validate() error {
	var missingFields []string

	if c.JobConfigName == "" {
		missingFields = append(missingFields, "job_config_name")
		logger.Error("Missing required field: job_config_name")
	}
	if c.JobID == "" {
		missingFields = append(missingFields, "job_id")
		logger.Error("Missing required field: job_id")
	}
	if c.DatasetID == "" {
		missingFields = append(missingFields, "dataset_id")
		logger.Error("Missing required field: dataset_id")
	}
	if c.BatchSize == "" {
		missingFields = append(missingFields, "batch_size")
		logger.Error("Missing required field: batch_size")
	}
	if c.NumEpochs == "" {
		missingFields = append(missingFields, "num_epochs")
		logger.Error("Missing required field: num_epochs")
	}
	if c.LearningRate == "" {
		missingFields = append(missingFields, "lr (learning rate)")
		logger.Error("Missing required field: lr")
	}
	if c.Seed == "" {
		missingFields = append(missingFields, "seed")
		logger.Error("Missing required field: seed")
	}
	if c.NumGpus == "" {
		missingFields = append(missingFields, "num_gpus")
		logger.Error("Missing required field: num_gpus")
	}
	if c.UserId == "" {
		missingFields = append(missingFields, "user_id")
		logger.Error("Missing required field: user_id")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required fields: %v", missingFields)
	}

	return nil
}
