package utils

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"wyvern-api/config"

	"github.com/sirupsen/logrus"
)

/**
* Logger Helper Using Logrus

Example Usage :

	import "wyvern-api/logger"

    log := logger.NewLogger("Main", 1).AddGWContext("ABC").Controller().Start()
		or
	log := logger.NewLogger("Main", 1)
		or
	log := logger.DefaultCBLogger()

	log.Info("Hello World")

*/

// DefaultLogDateFormat default log date format
const (
	DefaultLogDateFormat = "2006/01/02 15:04:05.000"
)

var (
	log     *logrus.Logger
	once    sync.Once
	logFile *os.File
)

// DefaultFormatter struct
type DefaultFormatter struct {
	*logrus.TextFormatter
}

// Format func to override format
func (f *DefaultFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format(DefaultLogDateFormat)
	level := strings.ToUpper(entry.Level.String())[0:1]
	//file := fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
	//funcName := fmt.Sprintf("%s()", entry.Caller.Function)
	msg := entry.Message

	return []byte(fmt.Sprintf("%s [%s] %s\n", timestamp, level, msg)), nil
}

// InitLogger func
func InitLogger() {
	once.Do(func() {
		log = logrus.New()
		logFileName := config.GetString("LOGGING_FILE_NAME", "logs/app.log")

		switch config.GetString("LOGGING_LEVEL", "debug") {
		case "debug":
			logrus.SetLevel(logrus.DebugLevel)
		case "info":
			logrus.SetLevel(logrus.InfoLevel)
		case "warn":
			logrus.SetLevel(logrus.WarnLevel)
		case "error":
			logrus.SetLevel(logrus.ErrorLevel)
		case "fatal":
			logrus.SetLevel(logrus.FatalLevel)
		case "panic":
			logrus.SetLevel(logrus.PanicLevel)
		case "trace":
			logrus.SetLevel(logrus.TraceLevel)
		default: //
			logrus.SetLevel(logrus.InfoLevel)
		}

		// Extract directory path from logFileName
		dirPath := filepath.Dir(logFileName)

		// Ensure the logs directory exists
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil && !os.IsExist(err) {
			log.Fatalf("Failed to create logs directory: %v", err)
		}

		// Create or open the log file
		var err error
		logFile, err = os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to create log file: %v", err)
		}

		// Set logrus to write to the log file and stdout
		logrus.SetOutput(io.MultiWriter(logFile, os.Stdout))

		logrus.SetReportCaller(true)

		logrus.SetFormatter(new(DefaultFormatter))
	})
}

type delegate func(format string, v ...interface{})

// Logger struct holds logging information
type Logger struct {
	*logrus.Entry
	Context        string
	CallStackLevel int
	StartTime      time.Time
	EnableStack    bool
	Identifier     int64
	GwIdentifier   int64
	GwContext      string
	LogLevel       string
	VersionContext string
	IsSetGWContext bool
}

// DefaultCBLogger initialize default Logger
func DefaultCBLogger() Logger {
	logEntry := logrus.NewEntry(log)

	cbLogger := Logger{}
	cbLogger.Entry = logEntry

	return cbLogger
}

// NewLogger creates a new Logger instance
func NewLogger(context string, callStackLevel int) *Logger {
	identifier := time.Now().UnixNano()
	logger := DefaultCBLogger()

	logger.Context = context
	logger.CallStackLevel = callStackLevel
	logger.StartTime = time.Now()
	logger.EnableStack = config.GetBool("LOGGING_WRAPPER_CALLSTACK", true)
	logger.Identifier = identifier
	logger.IsSetGWContext = false

	return &logger
}

// NewLoggerIdentifier initialize Logger with specific identifier
func NewLoggerIdentifier(context string, callStackLevel int, identifier int64) *Logger {
	logger := DefaultCBLogger()

	logger.Context = context
	logger.CallStackLevel = callStackLevel
	logger.StartTime = time.Now()
	logger.EnableStack = config.GetBool("LOGGING_WRAPPER_CALLSTACK", true)
	logger.Identifier = identifier
	logger.IsSetGWContext = false

	return &logger
}

// formatContext formats the log message with additional context
func (logger *Logger) formatContext(format string) string {
	if !logger.EnableStack {
		return format
	}

	var callstack string
	for i := logger.CallStackLevel + 3; i >= 3; i-- {
		_, file, line, _ := runtime.Caller(i)
		_, filename := path.Split(file)
		callstack = fmt.Sprintf("%s/%s:%d", callstack, filename, line)
	}

	if logger.GwIdentifier != 0 {
		logger.Identifier = logger.GwIdentifier
	}

	if logger.IsSetGWContext == false {
		logger.GwContext = ""
	}

	return fmt.Sprintf("[%s] [%X] %s [%s] %s %s%s", callstack, logger.Identifier, logger.GwContext, logger.Context, logger.LogLevel, logger.VersionContext, format)
}

// log logs a message with the given level and format
func (logger *Logger) log(fn delegate, format string, v ...interface{}) *Logger {
	message := logger.formatContext(format)
	fn(message, v...)

	return logger
}

// Error logs a message at Error level
func (logger *Logger) Error(format string, v ...interface{}) {
	logger.log(logrus.Errorf, format, v...)
}

// Warn logs a message at Warn level
func (logger *Logger) Warn(format string, v ...interface{}) {
	logger.log(logrus.Warnf, format, v...)
}

// Info logs a message at Info level
func (logger *Logger) Info(format string, v ...interface{}) {
	logger.log(logrus.Infof, format, v...)
}

// Debug logs a message at Debug level
func (logger *Logger) Debug(format string, v ...interface{}) {
	logger.log(logrus.Debugf, format, v...)
}

// Performance logs performance time
func (logger *Logger) Performance() *Logger {
	elapsed := time.Since(logger.StartTime)
	logger.WithField("performance", fmt.Sprintf("%fsec", elapsed.Seconds())).Info("Performance")
	return logger
}

// Start logs a message for service start
func (logger *Logger) Start() *Logger {
	logger.Info("[SERVICESTART]")
	return logger
}

// End logs a message for service end
func (logger *Logger) End() {
	logger.Info("[SERVICEEND]")
}

// Controller sets the log level to Controller context
func (logger *Logger) Controller() *Logger {
	logger.LogLevel = "[Controller]"
	return logger
}

// Service sets the log level to Service context
func (logger *Logger) Service() *Logger {
	logger.LogLevel = "[Service]"
	return logger
}

// Version sets the log level to Version context
func (logger *Logger) Version(version ...string) *Logger {
	ver := "1.0"

	if len(version) > 0 && version[0] != "" {
		ver = version[0]
	}

	logger.VersionContext = fmt.Sprintf("[%s] ", ver)
	return logger
}

// Host sets the log level to Host context
func (logger *Logger) Host() *Logger {
	logger.LogLevel = "[Host]"
	return logger
}

// Repository sets the log level to Repository context
func (logger *Logger) Repository() *Logger {
	logger.LogLevel = "[Repository]"
	return logger
}

// AddGWContext adds gateway context to the log
func (logger *Logger) AddGWContext(context string) *Logger {
	logger.GwContext = fmt.Sprintf("[%s]", context)
	return logger
}

// AddGWIdentifier adds gateway identifier to the log
func (logger *Logger) AddGWIdentifier(identifier string) *Logger {
	output, _ := strconv.ParseInt(identifier, 16, 64)
	logger.GwIdentifier = output
	return logger
}

// GWSIdentifier gets the gateway identifier as a string
func (logger *Logger) GWSIdentifier() string {
	return fmt.Sprintf("%X", logger.GwIdentifier)
}
