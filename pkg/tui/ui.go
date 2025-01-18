package tui

import (
	"fmt"

	"github.com/drksbr/lg2/pkg/config"
	"github.com/drksbr/lg2/pkg/parser"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TUI struct {
	// Main TUI
	App *tview.Application

	// Layout
	Grid       *tview.Grid
	LeftPannel *tview.Flex

	// Blocks
	Logo      *tview.TextView
	Shortcuts *tview.TextView
	PeersList *tview.List
	Content   *tview.TextView

	// Shortcut Actions
	SearchForm   *tview.Form
	NewQueryForm *tview.Form

	// States
	IsSearching bool
	IsQuerying  bool
	CurrentPeer int

	// Data
	originalPeers []parser.Peer
	filteredPeers []parser.Peer
}

var (
	mwLogo = `  [::b]▒▒  ▒▓ ▒▒▒▒▓  ▓▒▒▒▒▒
  ▓  ▓ ▒ ▒ ▒   ▓▒    ▒ 
 ▓   ▓▓   ▒ ▒  ▒ ▒  ▒  
 ▓▒▒▒▓▒ ▓▒▒▒▓ ▒▒▒  ▒▒[::-]    
    MultiGlass v%s`
)

// NewTUI configures and returns an instance of terminal user interface.
func NewTUI(queryString string) *TUI {

	// Blank Peer List
	peers := []parser.Peer{}

	// Create TUI
	tui := &TUI{
		App:           tview.NewApplication(),
		Logo:          tview.NewTextView(),
		Shortcuts:     tview.NewTextView().SetDynamicColors(true),
		PeersList:     tview.NewList().ShowSecondaryText(false),
		Content:       tview.NewTextView().SetDynamicColors(true).SetWrap(false),
		originalPeers: peers,
		filteredPeers: peers,
		SearchForm:    tview.NewForm(),
		NewQueryForm:  tview.NewForm(),
		IsSearching:   false,
		IsQuerying:    false,
		CurrentPeer:   0,
	}

	go func() {
		// Make query to API
		peers, err := tui.GetDataFromAPI(queryString)
		if err != nil {
			tui.App.QueueUpdateDraw(func() {
				tui.Content.SetText(fmt.Sprintf("Error: %s", err))
				tui.PeersList.SetTitle(" Peers(0) ")
			})
			return
		}

		tui.App.QueueUpdateDraw(func() {
			tui.originalPeers = peers
			tui.filteredPeers = peers
			tui.PeersList.SetTitle(fmt.Sprintf(" Peers(%d) ", len(peers)))
			for i, peer := range peers {
				tui.PeersList.AddItem(fmt.Sprintf("[%02d] %s", i+1, peer.PeerName), "", 0, func(index int) func() {
					return func() {
						tui.CurrentPeer = index
						tui.updateContent()
					}
				}(i))
			}
			if len(peers) > 0 {
				tui.updateContent()
			}
		})
	}()

	// Configure Logo
	tui.Logo.SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetText(fmt.Sprintf(mwLogo, config.Version))

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
	tui.Shortcuts.SetText("Change ['Tab'] / Quit ['q']\nFind ['f'] / Query ['n']\nNav [←][→] / Select [↓][↑]")

	// Criar Search Box
	tui.SearchForm.SetBackgroundColor(tcell.ColorDefault)
	tui.SearchForm.SetBorderColor(tcell.ColorDefault)
	tui.SearchForm.SetTitleColor(tcell.ColorDefault)
	tui.SearchForm.SetFieldBackgroundColor(tcell.ColorDarkGray)
	tui.SearchForm.SetFieldTextColor(tcell.ColorDefault)
	tui.SearchForm.SetTitle(" Find Peer ").SetBorder(true)

	// Configurar searchInput
	searchInput := tview.NewInputField()
	searchInput.SetChangedFunc(func(text string) {
		tui.filterAndUpdatePeersList(text)
	})

	captionBack := tview.NewTextView()
	captionBack.SetText("Back ['ESQ']").
		SetSize(1, 20).SetTextAlign(tview.AlignLeft)

	tui.SearchForm.AddFormItem(searchInput)
	tui.SearchForm.AddFormItem(captionBack)

	// Criar New Query Form
	tui.NewQueryForm.SetBackgroundColor(tcell.ColorDefault)
	tui.NewQueryForm.SetBorderColor(tcell.ColorDefault)
	tui.NewQueryForm.SetTitleColor(tcell.ColorDefault)
	tui.NewQueryForm.SetFieldBackgroundColor(tcell.ColorDarkGray)
	tui.NewQueryForm.SetFieldTextColor(tcell.ColorDefault)
	tui.NewQueryForm.SetTitle(" New Query ").SetBorder(true)

	// Configurar Query Form
	newQueryInput := tview.NewInputField()
	newQueryInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			tui.updateTUIWithNewQuery(newQueryInput.GetText())
			tui.App.SetFocus(newQueryInput) // Keep focus on the input field
			return
		}
	})

	tui.NewQueryForm.AddFormItem(newQueryInput)
	tui.NewQueryForm.AddFormItem(captionBack)

	// Configure Content
	tui.Content.SetBackgroundColor(tcell.ColorDefault)
	tui.Content.SetTextColor(tcell.ColorDefault)
	tui.Content.SetTitle(" Looking Glass Details ").SetBorder(true).SetBorderColor(tcell.ColorDefault)
	tui.Content.SetTitleColor(tcell.ColorDefault)
	tui.Content.SetWrap(false)

	// Setup Focus Cycling
	SetupFocusCycling(tui.App, tui.PeersList, tui.Content, tui.Grid)

	// Configure Peers List
	for i, peer := range peers {
		tui.PeersList.AddItem(fmt.Sprintf("[%02d] %s", i+1, peer.PeerName), "", 0, func(index int) func() {
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
	tui.LeftPannel = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tui.Logo, 7, 1, false).
		AddItem(tui.Shortcuts, 5, 1, false).
		AddItem(tui.PeersList, 0, 1, true)

	tui.Grid = tview.NewGrid().SetRows(0).SetColumns(30, 0).
		SetBorders(false).
		AddItem(tui.LeftPannel, 0, 0, 1, 1, 0, 0, true).
		AddItem(tui.Content, 0, 1, 1, 1, 0, 0, false)

	// Inicializar BoxComponent
	// tui.BoxComponent = NewBoxComponent(tui.Grid)

	// Adicionar setup de keyboard shortcuts
	tui.SetupKeyboardShortcuts()

	return tui
}

// Modificar a função de captura de eventos
func (tui *TUI) SetupKeyboardShortcuts() {
	tui.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if (event.Rune() == 'f' || event.Rune() == 'F') && !tui.IsSearching {
			tui.IsSearching = true
			// Substitui diretamente no Flex
			tui.LeftPannel.RemoveItem(tui.Shortcuts)
			tui.LeftPannel.RemoveItem(tui.PeersList)
			tui.LeftPannel.AddItem(tui.SearchForm, 7, 1, true)
			tui.LeftPannel.AddItem(tui.PeersList, 0, 1, false)

			tui.App.SetFocus(tui.SearchForm)
			return nil
		}

		if (event.Rune() == 'n' || event.Rune() == 'N') && !tui.IsSearching {
			tui.IsSearching = true
			// Substitui diretamente no Flex
			tui.LeftPannel.RemoveItem(tui.Shortcuts)
			tui.LeftPannel.RemoveItem(tui.PeersList)
			tui.LeftPannel.AddItem(tui.NewQueryForm, 7, 1, true)
			tui.LeftPannel.AddItem(tui.PeersList, 0, 1, false)

			tui.App.SetFocus(tui.NewQueryForm)
			return nil
		}

		if event.Key() == tcell.KeyEsc && tui.IsSearching {
			tui.IsSearching = false
			tui.LeftPannel.RemoveItem(tui.SearchForm)
			tui.LeftPannel.RemoveItem(tui.PeersList)
			tui.LeftPannel.RemoveItem(tui.NewQueryForm)
			tui.LeftPannel.AddItem(tui.Shortcuts, 5, 1, false)
			tui.LeftPannel.AddItem(tui.PeersList, 0, 1, true)

			tui.App.SetFocus(tui.PeersList)
			return nil
		}

		// Quit the application when 'q' or 'Q' and tui.IsSearching is false is pressed
		if (event.Rune() == 'q' || event.Rune() == 'Q') && !tui.IsSearching {
			tui.App.Stop()
			return nil
		}

		// Navigate through boxes
		if event.Key() == tcell.KeyTab {
			CycleFocus(tui.App, tui.PeersList, tui.Content, tui.Grid)
			return nil
		}

		return event
	})
}

// updateContent updates the details of the selected peer.
func (tui *TUI) updateContent() {
	if len(tui.filteredPeers) == 0 {
		tui.Content.SetText("Nenhum peer encontrado.")
		return
	}

	peer := tui.filteredPeers[tui.CurrentPeer]
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
