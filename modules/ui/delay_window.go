package ui

import (
	"errors"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type DelayWindow struct {
	fyne.Window

	delayEntry    *widget.Entry
	confirmButton *widget.Button
	linkID        int
}

func NewDelayWindow(app fyne.App, delay int64, linkId int) *DelayWindow {

	dw := &DelayWindow{linkID: linkId}

	delayWindow := app.NewWindow("Delay change")
	delayWindow.SetFixedSize(true)
	delayWindow.Resize(fyne.NewSize(405, 150))
	delayWindow.CenterOnScreen()
	dw.Window = delayWindow

	delayEntry := widget.NewEntry()
	dw.delayEntry = delayEntry
	delayEntry.Move(fyne.NewPos(10, 20))
	delayEntry.Resize(fyne.NewSize(380, 40))
	delayEntry.SetPlaceHolder("Enter delay in milliseconds")
	delayEntry.SetText(textFromDelay(delay))
	delayEntry.Validator = validateDelay
	delayEntry.Refresh()

	confirmButton := widget.NewButton("OK", nil)
	dw.confirmButton = confirmButton
	confirmButton.Move(fyne.NewPos(10, 80))
	confirmButton.Resize(fyne.NewSize(180, 40))

	cancelButton := widget.NewButton("Cancel", func() { dw.Close() })
	cancelButton.Move(fyne.NewPos(210, 80))
	cancelButton.Resize(fyne.NewSize(180, 40))

	dw.SetContent(container.NewWithoutLayout(
		delayEntry,
		confirmButton,
		cancelButton,
	))

	return dw
}

func (w *DelayWindow) OnConfirmDelay(f func()) {
	w.confirmButton.OnTapped = f
}

func textFromDelay(delay int64) string {
	if delay == 0 {
		return ""
	} else {
		return fmt.Sprint(delay)
	}
}

func validateDelay(s string) error {
	if s == "" {
		return nil
	}
	i, err := strconv.Atoi(s)
	if i < 0 {
		return errors.New("number must be positive")
	}
	return err
}
