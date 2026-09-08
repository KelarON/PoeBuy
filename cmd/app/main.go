package main

import (
	"poebuy/config"
	"poebuy/modules/bot"
	"poebuy/modules/ui"
	"poebuy/utils"

	_ "embed"
)

//go:embed PoeBuyUpdater.exe
var updaterFile []byte

func main() {

	// Create a logger instance
	logger := utils.NewLogger()
	defer logger.Close()

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil && err != config.ErrorNoConfigFile {
		logger.Errorf("config load failed: %v", err)
		return
	}

	// Create a new bot instance
	botInstance, err := bot.NewBot(cfg, logger)
	if err != nil {
		logger.Errorf("error app initialisation: %v", err)
		return
	}

	// Create a new UI instance
	ui := ui.NewUI(cfg, logger, botInstance)

	// Create an updater instance
	updater := utils.NewUpdater(logger)

	// Check for updates
	if cfg.General.AutoUpdates {
		go updater.PrepareUpdate(updaterFile, &cfg.Service.UpdateRequired, ui.ShowUpdateNotification)
	}

	// Run the UI
	ui.Run()

	// If an update is prepared, trigger the update process
	if cfg.Service.UpdateRequired && cfg.General.AutoUpdates {
		updater.TriggerUpdate()
	}
}
