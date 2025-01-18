package tui

import (
	"fmt"
	"time"

	"github.com/drksbr/lg2/pkg/fetch"
	"github.com/drksbr/lg2/pkg/parser"
)

func (tui *TUI) GetDataFromAPI(queryString string) ([]parser.Peer, error) {
	done := make(chan bool)
	tui.IsQuerying = true

	// Iniciar spinner
	go func() {
		spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-done:
				tui.IsQuerying = false
				return
			default:
				if tui.IsQuerying {
					spin := spinners[i%len(spinners)]
					tui.App.QueueUpdateDraw(func() {
						tui.PeersList.SetTitle(fmt.Sprintf(" Loading %s ", spin))
					})
					time.Sleep(100 * time.Millisecond)
					i++
				}
			}
		}
	}()

	// Check if queryString is valid
	network, err := parser.GetNetworkFromPrefix(queryString)
	if err != nil {
		done <- true
		return nil, err
	}

	// Check if network is valid
	if network == nil {
		done <- true
		return nil, fmt.Errorf("invalid network")
	}

	// Fetch data
	data, err := fetch.GetLookingGlassData(queryString)
	if err != nil {
		done <- true
		return nil, err
	}

	// Parse results
	peers, err := parser.ParseHTML(data, queryString)
	if err != nil {
		done <- true
		return nil, err
	}

	// Check results
	if len(peers) == 0 {
		done <- true
		return nil, fmt.Errorf("no peers found")
	}

	// Stop spinner and return
	done <- true
	return peers, nil
}

func (tui *TUI) updateTUIWithNewQuery(queryString string) {
	// Verificar se já está em consulta
	if tui.IsQuerying {
		return
	}

	done := make(chan bool)
	tui.IsQuerying = true

	// Iniciar spinner
	go func() {
		spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-done:
				return
			default:
				if tui.IsQuerying {
					spin := spinners[i%len(spinners)]
					tui.App.QueueUpdateDraw(func() {
						tui.PeersList.SetTitle(fmt.Sprintf(" Loading %s ", spin))
					})
					time.Sleep(100 * time.Millisecond)
					i++
				}
			}
		}
	}()

	// Buscar dados em goroutine
	go func() {
		defer func() {
			tui.IsQuerying = false
			done <- true
		}()

		newPeers, err := tui.GetDataFromAPI(queryString)
		if err != nil {
			tui.App.QueueUpdateDraw(func() {
				tui.Content.SetText(fmt.Sprintf("Error: %s", err))
				tui.PeersList.SetTitle(" Peers(0) ")
			})
			return
		}

		tui.App.QueueUpdateDraw(func() {
			tui.originalPeers = newPeers
			tui.filteredPeers = newPeers
			tui.PeersList.Clear()

			for i, peer := range newPeers {
				tui.PeersList.AddItem(fmt.Sprintf("[%02d] %s", i+1, peer.PeerName), "", 0, func(index int) func() {
					return func() {
						tui.CurrentPeer = index
						tui.updateContent()
					}
				}(i))
			}

			tui.PeersList.SetTitle(fmt.Sprintf(" Peers(%d) ", len(newPeers)))
			if len(newPeers) > 0 {
				tui.CurrentPeer = 0
				tui.updateContent()
			}
		})
	}()
}
