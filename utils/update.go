//go:build windows

package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"golang.org/x/sys/windows"
)

const (
	UPDATE_CACHE_DIR  = "UpdateCache"
	UPDATER_FILE_NAME = "PoeBuyUpdater.exe"
	UPDATE_FILE_NAME  = "PoeBuy.exe"
)

type Updater struct {
	logger *Logger
}

func NewUpdater(logger *Logger) *Updater {
	return &Updater{
		logger: logger,
	}
}

// PrepareUpdate checks for updates and prepares the update files if a new version is available
func (u *Updater) PrepareUpdate(updater []byte, updateRequired *bool, showUpdateNotification func()) {

	err := cleanupUpdateDirectory()
	if err != nil {
		u.logger.Errorf("Failed to start update process: %v", err)
		return
	}

	u.logger.Info("Checking for updates...")

	// Get the current version of the application
	curVersion := u.getCurrentVersion()

	// If the current version is "dev", we don't need to check for updates
	if curVersion == "dev" {
		return
	}

	// Fetch the latest release information from GitHub
	release, err := GetLatestRelease()
	if err != nil {
		u.logger.Errorf("Failed to check for updates: %v", err)
		return
	}

	// Compare versions
	if compareVersions(curVersion, release.GetVersion()) {
		u.logger.Infof("An update is available. Current version: %s, Latest version: %s", curVersion, release.GetVersion())
	} else {
		u.logger.Info("App is up to date.")
		return
	}

	// Create the update directory if it doesn't exist
	if _, err := os.Stat(UPDATE_CACHE_DIR); os.IsNotExist(err) {
		err := u.createUpdateDirectory()
		if err != nil {
			u.logger.Errorf("Failed to create update directory: %v", err)
			return
		}
	}

	// Save the updater file to the update directory
	err = saveUpdater(updater)
	if err != nil {
		u.logger.Errorf("Failed to save updater file: %v", err)
		return
	}

	// Save the downloaded file to the update directory
	updateFilePath := UPDATE_CACHE_DIR + "/" + UPDATE_FILE_NAME
	err = SaveRelease(release, updateFilePath)
	if err != nil {
		u.logger.Errorf("Failed to save update file: %v", err)
		return
	}

	u.logger.Infof("Update file saved to: %s", updateFilePath)

	// Set the updateRequired flag to true
	*updateRequired = true

	// Show the update notification
	showUpdateNotification()

}

// getCurrentVersion returns the current version of the application
func (u *Updater) getCurrentVersion() string {
	if fyne.CurrentApp().Metadata().Version == "0.0.1" {
		return "dev"
	}
	return "v" + fyne.CurrentApp().Metadata().Version
}

// CreateUpdateDirectory creates the update directory if it doesn't exist
func (u *Updater) createUpdateDirectory() error {
	err := os.Mkdir(UPDATE_CACHE_DIR, 0755)
	if err != nil {
		return err
	}
	return nil
}

func (u *Updater) TriggerUpdate() {

	updaterPath := filepath.Join(UPDATE_CACHE_DIR, UPDATER_FILE_NAME)

	cmd := exec.Command(updaterPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: windows.CREATE_NEW_PROCESS_GROUP | windows.DETACHED_PROCESS,
	}

	_ = cmd.Start()
}

// compareVersions compares two version strings and returns true if the latest version is greater than the current version.
func compareVersions(currentVersion, latestVersion string) bool {

	// Clean the version strings to remove any leading "v" or whitespace
	currentVersion = cleanVersion(currentVersion)
	latestVersion = cleanVersion(latestVersion)

	// Split the version strings into parts
	cParts := strings.Split(currentVersion, ".")
	lParts := strings.Split(latestVersion, ".")

	// Determine the maximum length of the version parts
	maxLen := len(cParts)
	if len(lParts) > maxLen {
		maxLen = len(lParts)
	}

	// Compare each part of the version strings
	for i := 0; i < maxLen; i++ {
		cNum := 0
		if i < len(cParts) {
			cNum, _ = strconv.Atoi(cParts[i])
		}

		lNum := 0
		if i < len(lParts) {
			lNum, _ = strconv.Atoi(lParts[i])
		}

		if lNum > cNum {
			return true
		}
		if lNum < cNum {
			return false
		}
	}

	// If all parts are equal, the versions are the same
	return false
}

func cleanVersion(v string) string {

	v = strings.TrimSpace(v)       // Remove leading and trailing whitespace
	v = strings.TrimPrefix(v, "v") // Remove leading "v" if present
	v = strings.TrimPrefix(v, "V") // Remove leading "V" if present
	return v
}

// saveUpdater saves the updater file to the update directory
func saveUpdater(updater []byte) error {

	updateFilePath := UPDATE_CACHE_DIR + "/" + UPDATER_FILE_NAME
	err := os.WriteFile(updateFilePath, updater, 0644)
	if err != nil {
		return err
	}
	return nil
}

// cleanupUpdateDirectory removes the update directory and its contents
func cleanupUpdateDirectory() error {

	// Try to remove the update directory up to 3 times with a delay in between
	for i := 0; i < 3; i++ {
		err := os.RemoveAll(UPDATE_CACHE_DIR)
		if err == nil {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("failed to remove update directory after 3 attempts")
}
