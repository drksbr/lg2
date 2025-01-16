package tui

import (
	"fmt"
	"strings"

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
		Shortcuts:   tview.NewTextView().SetDynamicColors(true).SetText("Select [↓][↑] / Quit ['q']\nFind ['f'] / Query ['n']\nNav [←][→] / Change ['Tab']"),
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
	tui.PeersList.SetTitle(" Peers ")
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
	var details strings.Builder
	details.WriteString(fmt.Sprintf("Peer Nome: %s\n", peer.PeerName))
	details.WriteString(fmt.Sprintf("AS-PATH: %s\n", formatASPath(peer.Path)))
	details.WriteString("Sequência:\n")
	for i, as := range peer.Path {
		details.WriteString(fmt.Sprintf("  [ %d ] %d - %s\n", i+1, as.AsNumber, as.AsName))
	}
	tui.Content.SetText(details.String())
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

// formatASPath converte o caminho de ASNs para uma string formatada.
func formatASPath(path []parser.AsPath) string {
	var formattedPath []string
	count := 1

	// Process consecutive repeated ASNs
	for i := 0; i < len(path); i++ {
		current := path[i].AsNumber

		// Count repetitions
		for i+1 < len(path) && path[i+1].AsNumber == current {
			count++
			i++
		}

		// Format the ASN with repetition count if needed
		if count > 1 {
			formattedPath = append(formattedPath, fmt.Sprintf("( %d x %d )", current, count))
		} else {
			formattedPath = append(formattedPath, fmt.Sprintf("%d", current))
		}
		count = 1
	}

	// Format with line breaks if too long
	if len(formattedPath) > 10 {
		var result strings.Builder
		for i, as := range formattedPath {
			if i > 0 && i%10 == 0 {
				result.WriteString("\n         ")
			}
			if i > 0 {
				result.WriteString(" « ")
			}
			result.WriteString(as)
		}
		return result.String()
	}

	return strings.Join(formattedPath, " « ")
}

// Start starts terminal user interface application.
func (tui *TUI) Start() error {
	return tui.App.SetRoot(tui.Grid, true).Run()
}
