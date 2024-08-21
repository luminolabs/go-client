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

var log = logrus.New()

type StandardLogger struct {
	*logrus.Logger
}

var standardLogger = &StandardLogger{logrus.New()}

var Address string
var Epoch uint32
var BlockNumber *big.Int
var FileName string
var Client *ethclient.Client

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)

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

func NewLogger() *StandardLogger {
	return standardLogger
}

func SetLogLevel(level string) {
	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
}
