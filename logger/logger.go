package logger

import (
	"io"
	"math/big"
	"os"
	"runtime"

	"lumino/core"
	"lumino/path"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/matishsiao/goInfo"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// StandardLogger wraps the logrus Logger
type StandardLogger struct {
	*logrus.Logger
}

var standardLogger = &StandardLogger{logrus.New()}

// Global variables for logging context
var (
	Address     string
	Epoch       uint32
	BlockNumber *big.Int
	FileName    string
	Client      *ethclient.Client
)

// init initializes the logger
func init() {
	standardLogger.SetOutput(os.Stdout)
	standardLogger.SetLevel(logrus.InfoLevel)

	InitializeLogger(FileName)

	osInfo, err := goInfo.GetInfo()
	if err != nil {
		standardLogger.Fatal("Error in fetching OS Info: ", err)
	}
	standardLogger.WithFields(logrus.Fields{
		"Operating System": osInfo.OS,
		"Core":             osInfo.Core,
		"Platform":         osInfo.Platform,
		"CPUs":             osInfo.CPUs,
		"razor-go version": core.VersionWithMeta,
		"go version":       runtime.Version(),
	}).Info()
}

// InitializeLogger sets up the logger with file rotation if a filename is provided
func InitializeLogger(fileName string) {
	if fileName != "" {
		logFilePath, err := path.PathUtilsInterface.GetLogFilePath(fileName)
		if err != nil {
			standardLogger.Fatal("Error in fetching log file path: ", err)
		}

		lumberJackLogger := &lumberjack.Logger{
			Filename:   logFilePath,
			MaxSize:    200, // Maximum Size of a log file
			MaxBackups: 52,  // Maximum number of log files
			MaxAge:     365, // Maximum number of days to retain old files
		}

		out := os.Stderr
		mw := io.MultiWriter(out, lumberJackLogger)
		standardLogger.Formatter = &logrus.JSONFormatter{}
		standardLogger.SetOutput(mw)

	} else {
		standardLogger.Formatter = &logrus.JSONFormatter{}
	}
}

// NewLogger returns a new instance of StandardLogger
func NewLogger() *StandardLogger {
	return standardLogger
}

// SetLogLevel sets the log level based on the provided string
func SetLogLevel(level string) {
	switch level {
	case "debug":
		standardLogger.SetLevel(logrus.DebugLevel)
	case "info":
		standardLogger.SetLevel(logrus.InfoLevel)
	case "warn":
		standardLogger.SetLevel(logrus.WarnLevel)
	case "error":
		standardLogger.SetLevel(logrus.ErrorLevel)
	default:
		standardLogger.SetLevel(logrus.InfoLevel)
	}
}

// Debug logs a debug message
func Debug(args ...interface{}) {
	standardLogger.Debug(args...)
}

// Info logs an info message
func Info(args ...interface{}) {
	standardLogger.Info(args...)
}

// Warn logs a warning message
func Warn(args ...interface{}) {
	standardLogger.Warn(args...)
}

// Error logs an error message
func Error(args ...interface{}) {
	standardLogger.Error(args...)
}

// Fatal logs a fatal message and exits
func Fatal(args ...interface{}) {
	standardLogger.Fatal(args...)
}
