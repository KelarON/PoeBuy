package ui

import (
	"errors"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SettingsWindow struct {
	fyne.Window

	delayEntry    *widget.Entry
	confirmButton *widget.Button
	cancelButton  *widget.Button
}

func NewSettingsWindow(app fyne.App) *SettingsWindow {

	dw := &SettingsWindow{}

	SettingsWindow := app.NewWindow("Settings")
	SettingsWindow.SetFixedSize(true)
	SettingsWindow.Resize(fyne.NewSize(405, 405))
	SettingsWindow.CenterOnScreen()
	dw.Window = SettingsWindow

	delayEntry := widget.NewEntry()
	dw.delayEntry = delayEntry
	delayEntry.Move(fyne.NewPos(10, 20))
	delayEntry.Resize(fyne.NewSize(380, 40))
	delayEntry.SetPlaceHolder("Enter delay in milliseconds")
	delayEntry.Validator = validateSettings
	delayEntry.Refresh()

	confirmButton := widget.NewButton("OK", nil)
	dw.confirmButton = confirmButton
	confirmButton.Move(fyne.NewPos(210, 80))
	confirmButton.Resize(fyne.NewSize(180, 40))

	cancelButton := widget.NewButton("Cancel", func() { dw.Close() })
	dw.cancelButton = cancelButton
	cancelButton.Move(fyne.NewPos(10, 80))
	cancelButton.Resize(fyne.NewSize(180, 40))

	dw.SetContent(container.NewWithoutLayout(
		delayEntry,
		confirmButton,
		cancelButton,
	))

	return dw
}

func (w *SettingsWindow) OnConfirmSettings(f func()) {
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
	if i < 0 {
		return errors.New("number must be positive")
	}
	return err
}
