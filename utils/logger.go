package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const logDirectory = "logs"

// Logger struct to hold the logger instance
type Logger struct {
	file     *os.File
	filePath string
	mu       sync.Mutex
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {

	// Prepare the log file path with the current date and time
	currentTime := time.Now().Format("2006-01-02")
	fileName := fmt.Sprintf("log_%s.log", currentTime)
	filePath := filepath.Join(logDirectory, fileName)

	return &Logger{filePath: filePath}
}

// ensureFile ensures the log file is created and opened
func (l *Logger) ensureFile() error {

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file == nil {
		// Ensure the logs directory exists
		if err := os.MkdirAll(filepath.Dir(l.filePath), 0755); err != nil {
			return fmt.Errorf("failed to create logs directory: %w", err)
		}

		// Open the log file
		file, err := os.OpenFile(l.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		l.file = file
		log.SetOutput(l.file)
	}

	return nil
}

// Info writes an info message to the log file
func (l *Logger) Info(message string) {

	if err := l.ensureFile(); err != nil {
		fmt.Printf("Error ensuring log file: %v\n", err)
		return
	}
	log.Println("[INFO]", message)
}

// Infof writes a formatted info message to the log file
func (l *Logger) Infof(format string, args ...interface{}) {

	if err := l.ensureFile(); err != nil {
		fmt.Printf("Error ensuring log file: %v\n", err)
		return
	}
	log.Printf("[INFO] "+format, args...)
}

// Warn writes a warning message to the log file
func (l *Logger) Warn(message string) {

	if err := l.ensureFile(); err != nil {
		fmt.Printf("Error ensuring log file: %v\n", err)
		return
	}
	log.Println("[WARN]", message)
}

// Warnf writes a formatted warning message to the log file
func (l *Logger) Warnf(format string, args ...interface{}) {

	if err := l.ensureFile(); err != nil {
		fmt.Printf("Error ensuring log file: %v\n", err)
		return
	}
	log.Printf("[WARN] "+format, args...)
}

// Error writes an error message to the log file
func (l *Logger) Error(message string) {

	if err := l.ensureFile(); err != nil {
		fmt.Printf("Error ensuring log file: %v\n", err)
		return
	}
	log.Println("[ERROR]", message)
}

// Errorf writes a formatted error message to the log file
func (l *Logger) Errorf(format string, args ...interface{}) {

	if err := l.ensureFile(); err != nil {
		fmt.Printf("Error ensuring log file: %v\n", err)
		return
	}
	log.Printf("[ERROR] "+format, args...)
}

// Close closes the log file
func (l *Logger) Close() {

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		l.file.Close()
		l.file = nil
	}
}
