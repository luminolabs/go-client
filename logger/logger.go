// Package logger implements structured logging functionality for the Lumino client,
// providing consistent log formatting, rotation, and level-based filtering.
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

// StandardLogger wraps logrus.Logger to provide additional functionality
// specific to Lumino's logging needs. Implements log rotation and
// structured logging with consistent formatting.
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

// init initializes the default logger configuration.
// Sets up output streams, log formatting, and captures basic system
// information for logging context. Must be called before any logging
// operations.
func init() {
	standardLogger.SetOutput(os.Stdout)
	standardLogger.SetLevel(logrus.InfoLevel)

	InitializeLogger(FileName)

	osInfo, err := goInfo.GetInfo()
	if err != nil {
		standardLogger.Fatal("Error in fetching OS Info: ", err)
	}
	standardLogger.WithFields(logrus.Fields{
		"Operating System":  osInfo.OS,
		"Core":              osInfo.Core,
		"Platform":          osInfo.Platform,
		"CPUs":              osInfo.CPUs,
		"lumino-go version": core.VersionWithMeta,
		"go version":        runtime.Version(),
	}).Info()
}

// InitializeLogger configures the logging system with file rotation support.
// Sets up log file paths, rotation policies, and output formatting.
// Supports both file and console logging with different formatters.
// Parameters:
// - fileName: Optional log file name. If empty, logs to stderr only.
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

// NewLogger returns a new instance of StandardLogger initialized
// with the default configuration.
func NewLogger() *StandardLogger {
	return standardLogger
}

// SetLoggerParameters configures global logging parameters.
// Sets up the blockchain client and account address for logging context.
// Starts a goroutine to track the latest block number for logging.
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

func SetLoggerParameters(client *ethclient.Client, address string) {
	Address = address
	Client = client
	go core.CalculateLatestBlock(client)
}
