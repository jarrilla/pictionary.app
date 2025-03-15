package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	infoLogger    *log.Logger
	errorLogger   *log.Logger
	warningLogger *log.Logger
	debugLogger   *log.Logger
)

// Init initializes the loggers with file and console output
func Init(logDir string) error {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	// Open the log file with the current date
	currentTime := time.Now()
	logFileName := filepath.Join(logDir, fmt.Sprintf("%s.log", currentTime.Format("2006-01-02")))
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	// Create multi-writer for both file and console
	infoWriter := io.MultiWriter(os.Stdout, logFile)
	errorWriter := io.MultiWriter(os.Stderr, logFile)

	// Initialize loggers with different prefixes and writers
	infoLogger = log.New(infoWriter, "INFO: ", log.Ldate|log.Ltime)
	errorLogger = log.New(errorWriter, "ERROR: ", log.Ldate|log.Ltime)
	warningLogger = log.New(infoWriter, "WARNING: ", log.Ldate|log.Ltime)
	debugLogger = log.New(infoWriter, "DEBUG: ", log.Ldate|log.Ltime)

	return nil
}

// getCallerInfo returns the file name and line number of the caller
func getCallerInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown:0"
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

// Info logs an info message with caller information
func Info(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	infoLogger.Printf("[%s] %s", getCallerInfo(), msg)
}

// Error logs an error message with caller information
func Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	errorLogger.Printf("[%s] %s", getCallerInfo(), msg)
}

// Warning logs a warning message with caller information
func Warning(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	warningLogger.Printf("[%s] %s", getCallerInfo(), msg)
}

// Debug logs a debug message with caller information
func Debug(format string, v ...interface{}) {
	if os.Getenv("NODE_ENV") == "development" {
		msg := fmt.Sprintf(format, v...)
		debugLogger.Printf("[%s] %s", getCallerInfo(), msg)
	}
}
