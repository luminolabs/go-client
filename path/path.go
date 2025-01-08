// Package path provides filesystem operations and path management functionality
// for the Lumino client, handling file access and directory operations securely.
package path

import (
	"os"
	pathPackage "path"
)

// GetLogFilePath generates and verifies the path for log files.
// Creates necessary directories if they don't exist and ensures
// proper file permissions are set. Returns the complete path
// to the log file or an error if setup fails.
func (PathUtils) GetLogFilePath(fileName string) (string, error) {
	luminoPath, err := PathUtilsInterface.GetDefaultPath()
	if err != nil {
		return "", err
	}
	defaultPath := pathPackage.Join(luminoPath, "logs")
	if _, err := OSUtilsInterface.Stat(defaultPath); OSUtilsInterface.IsNotExist(err) {
		mkdirErr := OSUtilsInterface.Mkdir(defaultPath, 0700)
		if mkdirErr != nil {
			return "", mkdirErr
		}
	}

	logFilepath := pathPackage.Join(defaultPath, fileName+".log")
	f, err := OSUtilsInterface.OpenFile(logFilepath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return logFilepath, nil
}

// GetDefaultPath returns the default Lumino working directory path.
// Creates the directory if it doesn't exist and ensures proper
// permissions are set. Usually points to ~/.lumino/.
func (PathUtils) GetDefaultPath() (string, error) {
	home, err := OSUtilsInterface.UserHomeDir()
	if err != nil {
		return "", err
	}
	defaultPath := pathPackage.Join(home, ".lumino")
	if _, err := OSUtilsInterface.Stat(defaultPath); OSUtilsInterface.IsNotExist(err) {
		mkdirErr := OSUtilsInterface.Mkdir(defaultPath, 0700)
		if mkdirErr != nil {
			return "", mkdirErr
		}
	}
	return defaultPath, nil
}

// GetConfigFilePath returns the path to the Lumino configuration file.
// Builds upon the default path to locate the config file.
// Returns an error if the default path cannot be determined.
func (PathUtils) GetConfigFilePath() (string, error) {
	luminoPath, err := PathUtilsInterface.GetDefaultPath()
	if err != nil {
		return "", err
	}
	return pathPackage.Join(luminoPath, "lumino.yaml"), nil
}
