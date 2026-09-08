package utils

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	searchPhrase   = "[LOADING SCREEN]"
	poeLogFolder   = "logs"
	poeLogFile     = "LatestClient.txt"
	loadingTimeout = time.Minute
)

// PoeLogMonitor manages log reading and trigger waiting
type PoeLogMonitor struct {
	logger       *Logger
	gamePath     string
	triggerCh    chan struct{}
	foundCh      chan struct{}
	updatePathCh chan string
	stopCh       chan struct{}
	isRunning    bool
}

// NewPoeLogMonitor creates and starts the log monitor in a background goroutine
func NewPoeLogMonitor(logger *Logger, gamePath string) *PoeLogMonitor {

	monitor := &PoeLogMonitor{
		logger:       logger,
		gamePath:     gamePath,
		triggerCh:    make(chan struct{}, 1),
		foundCh:      make(chan struct{}, 1),
		updatePathCh: make(chan string, 1),
		stopCh:       make(chan struct{}, 1),
	}

	return monitor
}

// Start starts reading in background
func (m *PoeLogMonitor) Start() {

	if !m.isRunning {
		go m.startReading()
		m.isRunning = true
	}
}

// UpdateGamePath safely requests the monitor to switch to a new game path
func (m *PoeLogMonitor) UpdateGamePath(newPath string) {

	if m.isRunning {
		select {
		case m.updatePathCh <- newPath:
		default:
			// If the channel is full, an update is already pending.
			select {
			case <-m.updatePathCh:
			default:
			}
			m.updatePathCh <- newPath
		}
	} else {
		m.gamePath = newPath
	}
}

// Stop gracefully shuts down the background reader
func (m *PoeLogMonitor) Stop() {

	if m.isRunning {
		select {
		case m.stopCh <- struct{}{}:
			m.isRunning = false
		default:
			// If the channel is full, reader is already stoping
		}
	}
}

// RequestLoadingScreen sends a request to search for the next searchPhrase
// and blocks until it is found (or until a timeout occurs)
func (m *PoeLogMonitor) RequestLoadingScreen() bool {

	// Do nothing, if reader not running
	if !m.isRunning {
		return false
	}

	// Clear the channel in case of old triggers (race condition protection)
	select {
	case <-m.foundCh:
	default:
	}

	// Send a signal to the log reader
	select {
	case m.triggerCh <- struct{}{}:
	default:
		// If the channel is full, a request is already being processed
	}

	// Wait for a response with a timeout
	select {
	case <-m.foundCh:
		return true // Successfully found
	case <-time.After(loadingTimeout):
		m.logger.Warn("Timeout waiting for location loading")
		return false // Timeout
	}
}

// startReading is the internal method that runs in a background goroutine
func (m *PoeLogMonitor) startReading() {

	var file *os.File
	var reader *bufio.Reader
	var lastSize int64
	var logPath string
	isWaitingForTrigger := true

	if m.gamePath != "" {
		logPath = filepath.Join(m.gamePath, poeLogFolder, poeLogFile)
	}

	openFile := func() error {
		if file != nil {
			file.Close()
			file = nil // Prevent double close
		}
		var err error
		file, err = os.Open(logPath)
		if err != nil {
			return err
		}
		pos, err := file.Seek(0, io.SeekEnd)
		if err != nil {
			return err
		}
		lastSize = pos
		reader = bufio.NewReader(file)
		return nil
	}

	// Initial opening attempt
	_ = openFile()

	for {
		select {
		case newPath := <-m.updatePathCh:
			m.gamePath = newPath
			logPath = filepath.Join(m.gamePath, poeLogFolder, poeLogFile)

			// Close old file and open the new one
			if err := openFile(); err != nil {
				m.logger.Errorf("Failed to open new log file: %v", err)
			}

			// Reset state for the new file
			isWaitingForTrigger = true
			continue // Skip to the next iteration to re-evaluate file state

		case <-m.stopCh:
			if file != nil {
				file.Close()
			}
			m.logger.Info("Log monitor stopped gracefully")
			return

		case <-m.triggerCh:
			isWaitingForTrigger = false

		default:
		}

		// If the file is not open (e.g., initial open failed or update failed), wait and retry
		if file == nil {
			time.Sleep(1 * time.Second)
			_ = openFile()
			continue
		}

		// Check for file recreation
		currentInfo, err := os.Stat(logPath)
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		if currentInfo.Size() < lastSize {
			if err := openFile(); err != nil {
				time.Sleep(time.Second)
				continue
			}
			isWaitingForTrigger = false
			m.logger.Info("Log file recreated and reopened")
		}

		// Read the line
		line, err := reader.ReadString('\n')
		if err == nil {
			line = strings.TrimSpace(line)

			// If we are NOT waiting for a trigger and found the target phrase
			if !isWaitingForTrigger && strings.Contains(line, searchPhrase) {
				select {
				case m.foundCh <- struct{}{}:
				default:
				}
				isWaitingForTrigger = true
			}

			lastSize, _ = file.Seek(0, io.SeekCurrent)
			continue
		}

		if err == io.EOF {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		m.logger.Errorf("Read error: %v", err)
		time.Sleep(time.Second)
	}
}
