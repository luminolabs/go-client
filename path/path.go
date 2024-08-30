package path

import (
	"os"
	pathPackage "path"
)

// GetLogFilePath returns the path for the log file
// It creates the necessary directories if they don't exist
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

// This function returns the default path
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

// This function returns the config file path
func (PathUtils) GetConfigFilePath() (string, error) {
	luminoPath, err := PathUtilsInterface.GetDefaultPath()
	if err != nil {
		return "", err
	}
	return pathPackage.Join(luminoPath, "lumino.yaml"), nil
}
