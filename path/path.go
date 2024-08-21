package path

import (
	"os"
	pathPackage "path"
)

// This function returns the log file path
func (PathUtils) GetLogFilePath(fileName string) (string, error) {
	razorPath, err := PathUtilsInterface.GetDefaultPath()
	if err != nil {
		return "", err
	}
	defaultPath := pathPackage.Join(razorPath, "logs")
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
