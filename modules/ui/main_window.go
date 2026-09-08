package ui

import (
	"fmt"
	"image/color"
	"poebuy/config"
	"poebuy/modules/connections/models"
	"poebuy/utils"
	"slices"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	DefaultWindowWidth  = 800
	DefaultWindowHeight = 605
)

type MainWindow struct {
	fyne.Window

	leagueSelect    *widget.Select
	nameEntry       *widget.Entry
	linkEntry       *widget.Entry
	addTradeButton  *widget.Button
	tradeTable      *widget.Table
	visitDelayEntry *widget.Entry
	linkCopyPopup   *widget.PopUp
	logoutButton    *widget.Button
	settingsButton  *widget.Button
}

func NewMainWindow(app fyne.App, info *models.TradeInfo, cfg *config.Config) *MainWindow {

	mw := &MainWindow{}

	mw.Window = app.NewWindow("PoeBuy")
	mw.SetFixedSize(true)
	mw.Resize(fyne.NewSize(DefaultWindowWidth, DefaultWindowHeight))

	mw.linkCopyPopup = widget.NewPopUp(widget.NewLabel("✔ Link copied"), mw.Canvas())

	leagueLabel := widget.NewLabel("League:")
	leagueLabel.Move(fyne.NewPos(15, 10))

	if !slices.Contains(info.GetLeagues(), cfg.Trade.League) {
		cfg.Trade.League = info.GetLeagues()[0]
		cfg.Save()
	}
	leagueBind := binding.BindString(&cfg.Trade.League)
	leagueBind.AddListener(binding.NewDataListener(cfg.Save))
	leagueSelect := widget.NewSelectWithData(info.GetLeagues(), leagueBind)
	mw.leagueSelect = leagueSelect
	leagueSelect.Move(fyne.NewPos(15, 50))
	leagueSelect.Resize(fyne.NewSize(350, 35))
	leagueSelect.PlaceHolder = "Select league"
	leagueSelect.Refresh()

	nicknameLabel := widget.NewLabel("Logged in as " + info.Nickname)
	nicknameLabel.Alignment = fyne.TextAlignTrailing
	nicknameLabel.Move(fyne.NewPos(700, 10))
	nicknameLabel.TextStyle = fyne.TextStyle{Bold: true}

	logoutButton := widget.NewButtonWithIcon("Logout", theme.LogoutIcon(), nil)
	mw.logoutButton = logoutButton
	logoutButton.Move(fyne.NewPos(700, 17))
	logoutButton.Resize(fyne.NewSize(80, 23))

	settingsButton := widget.NewButtonWithIcon("", theme.SettingsIcon(), nil)
	mw.settingsButton = settingsButton
	settingsButton.Move(fyne.NewPos(743, 90))
	settingsButton.Resize(fyne.NewSize(40, 40))

	addTradeLabel := widget.NewLabel("Add trade links:")
	addTradeLabel.Move(fyne.NewPos(15, 100))

	nameEntry := widget.NewEntry()
	mw.nameEntry = nameEntry
	nameEntry.Move(fyne.NewPos(15, 140))
	nameEntry.Resize(fyne.NewSize(765, 40))
	nameEntry.PlaceHolder = "Trade name"
	nameEntry.Refresh()

	linkEntry := widget.NewEntry()
	mw.linkEntry = linkEntry
	linkEntry.Move(fyne.NewPos(15, 190))
	linkEntry.Resize(fyne.NewSize(765, 40))
	linkEntry.PlaceHolder = "Trade link or code"
	linkEntry.Refresh()

	addTradeRectangle := canvas.NewRectangle(nil)
	addTradeRectangle.Move(fyne.NewPos(10, 135))
	addTradeRectangle.Resize(fyne.NewSize(775, 152))
	addTradeRectangle.StrokeWidth = 1
	addTradeRectangle.StrokeColor = color.RGBA{R: 60, G: 60, B: 60, A: 255}
	addTradeRectangle.CornerRadius = 5
	addTradeRectangle.Refresh()

	addTradeButton := widget.NewButton("Add", nil)
	mw.addTradeButton = addTradeButton
	addTradeButton.Move(fyne.NewPos(580, 240))
	addTradeButton.Resize(fyne.NewSize(200, 40))

	tradeTable := widget.NewTable(
		func() (int, int) {
			return len(cfg.Trade.Links), 5
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			icon := widget.NewIcon(nil)
			icon.Hide()
			return container.NewStack(label, icon)
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {

			container := o.(*fyne.Container)
			label := container.Objects[0].(*widget.Label)
			icon := container.Objects[1].(*widget.Icon)

			switch i.Col {
			case 0:
				label.Show()
				icon.Hide()
				label.SetText(cfg.Trade.Links[i.Row].Name)
			case 1:
				label.Show()
				icon.Hide()
				label.SetText("Click to copy")
			case 2:
				label.Hide()
				icon.Show()
				if cfg.Trade.Links[i.Row].IsActiv {
					icon.SetResource(theme.CheckButtonCheckedIcon())
				} else {
					icon.SetResource(theme.CheckButtonIcon())
				}
			case 3:
				label.Hide()
				icon.Show()
				icon.SetResource(theme.DeleteIcon())
			case 4:
				label.Show()
				icon.Hide()
				label.SetText(millisecondsToHumanReadable(cfg.Trade.Links[i.Row].Delay))
			}
		})
	mw.tradeTable = tradeTable
	tradeTable.Move(fyne.NewPos(15, 300))
	tradeTable.Resize(fyne.NewSize(765, 280))
	tradeTable.SetColumnWidth(0, 430)
	tradeTable.SetColumnWidth(1, 100)
	tradeTable.SetColumnWidth(2, 55)
	tradeTable.SetColumnWidth(3, 60)
	tradeTable.SetColumnWidth(4, 100)
	tradeTable.CreateHeader = func() fyne.CanvasObject {
		label := widget.NewLabel("")
		return label
	}
	tradeTable.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		label := o.(*widget.Label)
		switch id.Col {
		case 0:
			label.SetText("Name")
		case 1:
			label.SetText("Link")
		case 2:
			label.SetText("Active")
		case 3:
			label.SetText("Delete")
		case 4:
			label.SetText("Delay")
		}
	}
	tradeTable.ShowHeaderRow = true
	tradeTable.Refresh()

	versionLabel := widget.NewLabel("version: " + utils.GetCurrentVersion())
	versionLabel.SizeName = theme.SizeNameCaptionText
	versionLabel.Move(fyne.NewPos(700, 573))

	mw.SetContent(container.NewWithoutLayout(
		leagueLabel,
		leagueSelect,
		nicknameLabel,
		logoutButton,
		settingsButton,
		addTradeLabel,
		nameEntry,
		linkEntry,
		addTradeRectangle,
		addTradeButton,
		tradeTable,
		versionLabel,
	))

	return mw
}

func (w *MainWindow) OnAddTrade(f func()) {
	w.addTradeButton.OnTapped = f
}

func (w *MainWindow) OnTableCellClick(f func(id widget.TableCellID)) {
	w.tradeTable.OnSelected = f
}

func (w *MainWindow) OnLogout(f func()) {
	w.logoutButton.OnTapped = f
}

func (w *MainWindow) OnSettings(f func()) {
	w.settingsButton.OnTapped = f
}

func (w *MainWindow) ShowLinkCopyPopup() {
	fyne.Do(func() {
		w.linkCopyPopup.ShowAtPosition(fyne.NewPos(DefaultWindowWidth/2, DefaultWindowHeight-w.linkCopyPopup.Size().Height))
	})
	time.Sleep(1500 * time.Millisecond)
	fyne.Do(w.linkCopyPopup.Hide)
}

// millisecondsToHumanReadable converts milliseconds to a human-readable string
// If the input is 0, it returns "no delay"
func millisecondsToHumanReadable(ms int64) string {

	if ms == 0 {
		return "no delay"
	}

	t := time.Duration(ms) * time.Millisecond

	switch {
	case t < time.Second:
		return fmt.Sprintf("%v ms", float32(ms))
	case t < time.Minute:
		return fmt.Sprintf("%v s", float32(ms)/1000)
	case t < time.Hour:
		return fmt.Sprintf("%v min", float32(ms)/1000/60)
	default:
		return fmt.Sprintf("%v h", float32(ms)/1000/60/60)
	}
}
