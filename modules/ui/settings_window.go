package ui

import (
	"errors"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"poebuy/config"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	nativeDialog "github.com/sqweek/dialog"
)

type SettingsWindow struct {
	fyne.Window

	confirmButton *widget.Button
	workingCfg    *config.Config
}

func NewSettingsWindow(app fyne.App, cfg *config.Config) *SettingsWindow {

	sw := &SettingsWindow{}

	workingCfg := new(config.Config)
	*workingCfg = *cfg
	sw.workingCfg = workingCfg

	settingsWindow := app.NewWindow("Settings")
	settingsWindow.SetFixedSize(true)
	settingsWindow.Resize(fyne.NewSize(405, 405))
	settingsWindow.CenterOnScreen()
	sw.Window = settingsWindow

	visitDelayBind := binding.BindInt(&sw.workingCfg.Trade.VisitDelay)
	visitDelayBindText := binding.IntToString(visitDelayBind)
	visitDelayEntry := widget.NewEntryWithData(visitDelayBindText)
	visitDelayEntry.Move(fyne.NewPos(10, 20))
	visitDelayEntry.Resize(fyne.NewSize(73, 36))
	visitDelayEntry.Validator = validateSettings

	visitDelayLabel := widget.NewLabel("Visit time (in seconds)")
	visitDelayLabel.Move(fyne.NewPos(80, 20))

	readLogBind := binding.BindBool(&sw.workingCfg.Trade.ReadLog)
	readLogCheckbox := widget.NewCheckWithData("Read game logs (to determine the load)", readLogBind)
	readLogCheckbox.Move(fyne.NewPos(5, 80))
	readLogCheckbox.Resize(fyne.NewSize(20, 20))

	gamePathLabel := widget.NewLabel("Path of Exile folder")
	gamePathLabel.Move(fyne.NewPos(115, 120))

	gamePathChosenLabel := widget.NewLabel("Current path:")
	gamePathChosenLabel.Move(fyne.NewPos(5, 150))

	gamePathChosenBind := binding.BindString(&sw.workingCfg.Trade.GamePath)
	gamePathChosenLabel2 := widget.NewLabelWithData(gamePathChosenBind)
	gamePathChosenLabel2.SizeName = theme.SizeNameCaptionText
	gamePathChosenLabel2.Move(fyne.NewPos(100, 153))

	gamePathButton := widget.NewButtonWithIcon("Browse...", theme.FolderIcon(), func() {
		tempWindow := app.NewWindow("Choose game folder")
		tempWindow.Resize(fyne.NewSize(DefaultWindowWidth, DefaultWindowHeight))
		tempWindow.SetFixedSize(true)
		tempWindow.Show()
		fileDialog := dialog.NewFolderOpen(func(lu fyne.ListableURI, err error) {
			if lu != nil && isValidPoEFolder(lu.Path()) {
				sw.workingCfg.Trade.GamePath = lu.Path()
				gamePathChosenBind.Reload()
			}
			tempWindow.Close()
			if !isValidPoEFolder(lu.Path()) {
				nativeDialog.Message("This folder doesn't contain Path of Exile.\n\nPlease select the folder with PathOfExile.exe or PathOfExileSteam.exe").Title("Error").Error()
			}
		}, tempWindow)
		fileDialog.Show()
		fileDialog.Resize(fyne.NewSize(DefaultWindowWidth, DefaultWindowHeight))
	})
	gamePathButton.Move(fyne.NewPos(10, 125))
	gamePathButton.Resize(fyne.NewSize(100, 25))

	hideoutSettingsLabel := widget.NewLabel("Hideout travel settings:")
	hideoutSettingsLabel.SizeName = theme.SizeNameCaptionText
	hideoutSettingsLabel.Move(fyne.NewPos(20, -11))

	hideoutSettingsRectangle := canvas.NewRectangle(nil)
	hideoutSettingsRectangle.Move(fyne.NewPos(5, 5))
	hideoutSettingsRectangle.Resize(fyne.NewSize(390, 180))
	hideoutSettingsRectangle.StrokeWidth = 1
	hideoutSettingsRectangle.StrokeColor = color.RGBA{R: 60, G: 60, B: 60, A: 255}
	hideoutSettingsRectangle.CornerRadius = 5
	hideoutSettingsRectangle.Refresh()

	updateBind := binding.BindBool(&sw.workingCfg.General.AutoUpdates)
	updateCheckbox := widget.NewCheckWithData("Automatic app updates", updateBind)
	updateCheckbox.Move(fyne.NewPos(5, 310))
	updateCheckbox.Resize(fyne.NewSize(20, 20))

	var updateStatus string
	if cfg.Service.UpdateRequired {
		updateStatus = "(Update is ready)"
	} else {
		updateStatus = "(No updates available)"
	}
	updateStatusLabel := widget.NewLabel(updateStatus)
	updateStatusLabel.Alignment = fyne.TextAlignTrailing
	updateStatusLabel.Move(fyne.NewPos(390, 303))

	appSettingsLabel := widget.NewLabel("App settings:")
	appSettingsLabel.SizeName = theme.SizeNameCaptionText
	appSettingsLabel.Move(fyne.NewPos(20, 284))

	appSettingsRectangle := canvas.NewRectangle(nil)
	appSettingsRectangle.Move(fyne.NewPos(5, 300))
	appSettingsRectangle.Resize(fyne.NewSize(390, 40))
	appSettingsRectangle.StrokeWidth = 1
	appSettingsRectangle.StrokeColor = color.RGBA{R: 60, G: 60, B: 60, A: 255}
	appSettingsRectangle.CornerRadius = 5
	appSettingsRectangle.Refresh()

	saveButton := widget.NewButton("Save", nil)
	sw.confirmButton = saveButton
	saveButton.Move(fyne.NewPos(10, 350))
	saveButton.Resize(fyne.NewSize(180, 40))

	cancelButton := widget.NewButton("Cancel", func() { sw.Close() })
	cancelButton.Move(fyne.NewPos(210, 350))
	cancelButton.Resize(fyne.NewSize(180, 40))

	sw.SetContent(container.NewWithoutLayout(
		saveButton,
		cancelButton,
		visitDelayEntry,
		visitDelayLabel,
		readLogCheckbox,
		gamePathButton,
		gamePathLabel,
		gamePathChosenLabel,
		gamePathChosenLabel2,
		hideoutSettingsRectangle,
		hideoutSettingsLabel,
		appSettingsRectangle,
		appSettingsLabel,
		updateCheckbox,
		updateStatusLabel,
	))

	return sw
}

func (w *SettingsWindow) OnSaveSettings(f func()) {
	w.confirmButton.OnTapped = f
}

func textFromSettings(delay int64) string {
	if delay == 0 {
		return ""
	} else {
		return fmt.Sprint(delay)
	}
}

func validateSettings(s string) error {
	if s == "" {
		return nil
	}
	i, err := strconv.Atoi(s)
	if i <= 0 {
		return errors.New("number must be greater than 0")
	}
	return err
}

func isValidPoEFolder(path string) bool {

	files := []string{
		"PathOfExile.exe",
		"PathOfExile_x64.exe",
		"PathOfExileSteam.exe",
		"PathOfExile_x64Steam.exe",
	}

	for _, file := range files {
		if _, err := os.Stat(filepath.Join(path, file)); err == nil {
			return true
		}
	}

	return false
}
