package tui

import (
	"fmt"

	"github.com/drksbr/lg2/pkg/parser"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TUI struct {
	App         *tview.Application
	Grid        *tview.Grid
	Logo        *tview.TextView
	Shortcuts   *tview.TextView
	PeersList   *tview.List
	Content     *tview.TextView
	Peers       []parser.Peer
	CurrentPeer int
}

var (
	mwLogo = `  ▒▒  ▒▓ ▒▒▒▒▓  ▓▒▒▒▒▒
  ▓  ▓ ▒ ▒ ▒   ▓▒    ▒ 
 ▓   ▓▓   ▒ ▒  ▒ ▒  ▒  
 ▓▒▒▒▓▒ ▓▒▒▒▓ ▒▒▒  ▒▒    
    MultiGlass V0.1`
)

// NewTUI configures and returns an instance of terminal user interface.
func NewTUI(peers []parser.Peer) *TUI {
	peersList := tview.NewList().ShowSecondaryText(false)
	tui := &TUI{
		App:         tview.NewApplication(),
		Logo:        tview.NewTextView().SetTextAlign(tview.AlignCenter).SetDynamicColors(true).SetText(mwLogo),
		Shortcuts:   tview.NewTextView().SetDynamicColors(true).SetText("Change ['Tab'] / Quit ['q']\nFind ['f'] / Query ['n']\nNav [←][→] / Select [↓][↑]"),
		PeersList:   peersList,
		Content:     tview.NewTextView().SetDynamicColors(true).SetWrap(false),
		Peers:       peers,
		CurrentPeer: 0,
	}

	// Configure Peer List
	tui.PeersList.SetBackgroundColor(tcell.ColorDefault)
	tui.PeersList.SetMainTextStyle(tcell.StyleDefault)
	tui.PeersList.SetHighlightFullLine(true)
	tui.PeersList.SetSelectedBackgroundColor(tcell.ColorDarkGrey)
	tui.PeersList.SetBorderColor(tcell.ColorDefault)
	tui.PeersList.SetTitleColor(tcell.ColorDefault)
	tui.PeersList.SetTitle(fmt.Sprintf(" %s(%d) ", "Peers", len(peers)))
	tui.PeersList.SetBorder(true)

	// Run function on item selected on peerlist
	// Run function when selection changes with arrow keys
	tui.PeersList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		tui.CurrentPeer = index
		tui.updateContent()
	})

	// Tui Logo
	tui.Logo.SetBackgroundColor(tcell.ColorDefault)
	tui.Logo.SetTextColor(tcell.ColorDefault)
	tui.Logo.SetBorder(true).SetBorderColor(tcell.ColorDefault)

	// Tui Shortcuts
	tui.Shortcuts.SetTextColor(tcell.ColorDefault)
	tui.Shortcuts.SetBackgroundColor(tcell.ColorDefault)
	tui.Shortcuts.SetBorderColor(tcell.ColorDefault)
	tui.Shortcuts.SetTitleColor(tcell.ColorDefault)
	tui.Shortcuts.SetTitle(" Shortcuts ").SetBorder(true)

	// Configure Content
	tui.Content.SetBackgroundColor(tcell.ColorDefault)
	tui.Content.SetTextColor(tcell.ColorDefault)
	tui.Content.SetTitle(" Looking Glass Details ").SetBorder(true).SetBorderColor(tcell.ColorDefault)
	tui.Content.SetTitleColor(tcell.ColorDefault)
	tui.Content.SetWrap(true)

	// Setup Focus Cycling
	SetupFocusCycling(tui.App, tui.PeersList, tui.Content, tui.Grid)

	// Configure Peers List
	for i, peer := range peers {
		tui.PeersList.AddItem(fmt.Sprintf("[ %d ] %s", i+1, peer.PeerName), "", 0, func(index int) func() {
			return func() {
				tui.CurrentPeer = index
				tui.updateContent()
			}
		}(i))
	}

	if len(peers) > 0 {
		tui.updateContent()
	}

	// Configure Grid Layout
	left := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tui.Logo, 7, 1, false).
		AddItem(tui.Shortcuts, 5, 1, false).
		AddItem(tui.PeersList, 0, 1, true)

	tui.Grid = tview.NewGrid().SetRows(0).SetColumns(30, 0).
		SetBorders(false).
		AddItem(left, 0, 0, 1, 1, 0, 0, true).
		AddItem(tui.Content, 0, 1, 1, 1, 0, 0, false)

	// Configure Keybindings
	tui.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab: // Cycle through the content box
			CycleFocus(tui.App, tui.PeersList, tui.Content, tui.Grid)

		case tcell.KeyBacktab: // Cycle through the focusable elements
			CycleFocus(tui.App, tui.PeersList, tui.Content, tui.Grid)

		case tcell.KeyRune:
			if event.Rune() == 'q' || event.Rune() == 'Q' {
				tui.App.Stop()
			}
		}
		return event
	})

	return tui
}

// updateContent updates the details of the selected peer.
func (tui *TUI) updateContent() {
	if len(tui.Peers) == 0 {
		tui.Content.SetText("Nenhum peer encontrado.")
		return
	}

	peer := tui.Peers[tui.CurrentPeer]
	details := buildPeerDetails(&peer)

	// Set the text of the content box
	tui.Content.SetText(details)
}

// CycleFocus alterna o foco entre os quadros de lista de peers e o conteúdo.
func CycleFocus(app *tview.Application, list *tview.List, content *tview.TextView, grid *tview.Grid) {
	focused := app.GetFocus()

	// Verificar qual elemento atualmente tem o foco e alternar
	if focused == list {
		app.SetFocus(content)
	} else if focused == content {
		app.SetFocus(list)
	}
}

// Exemplo de uso no programa principal
func SetupFocusCycling(app *tview.Application, list *tview.List, content *tview.TextView, grid *tview.Grid) {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			CycleFocus(app, list, content, grid)
			return nil
		}
		return event
	})
}

// Start starts terminal user interface application.
func (tui *TUI) Start() error {
	return tui.App.SetRoot(tui.Grid, true).Run()
}
