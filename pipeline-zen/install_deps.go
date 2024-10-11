package pipeline_zen

import (
	"lumino/logger"
	"os/exec"
)

// InstallDeps runs the install-deps.sh script to install Python dependencies
func InstallDeps() (string, error) {
	cmd := exec.Command("bash", "-c", "./scripts/install-deps.sh")

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Error installing dependencies: ", err)
		return "", err
	}

	logger.Info("Dependencies installed successfully")
	return string(output), nil
}
