package pipeline_zen

import (
	"fmt"
	"os/exec"
)

// InstallDeps runs the install-deps.sh script to install Python dependencies
func InstallDeps() (string, error) {
	cmd := exec.Command("bash", "-c", "./scripts/install-deps.sh")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Error installing dependencies: %v", err)
	}

	return string(output), nil
}
