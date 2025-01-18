package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/drksbr/lg2/pkg/parser"
)

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

// Função para filtrar e atualizar lista
func (tui *TUI) filterAndUpdatePeersList(searchTerm string) {
	// Limpar lista atual
	tui.PeersList.Clear()

	// Filtrar peers
	tui.filteredPeers = []parser.Peer{}
	for _, peer := range tui.originalPeers {
		if strings.Contains(strings.ToLower(peer.PeerName), strings.ToLower(searchTerm)) {
			tui.filteredPeers = append(tui.filteredPeers, peer)
		}
	}

	// Atualizar lista na UI
	for i, peer := range tui.filteredPeers {
		tui.PeersList.AddItem(fmt.Sprintf("[%02d] %s", i+1, peer.PeerName), "", 0, func(index int) func() {
			return func() {
				tui.CurrentPeer = index
				tui.updateContent()
			}
		}(i))
	}

	// Atualizar título com quantidade
	tui.PeersList.SetTitle(fmt.Sprintf(" Peers(%d) ", len(tui.filteredPeers)))
}

func (tui *TUI) updateTUIWithNewQuery(queryString string) {
	done := make(chan bool)

	// Debug inicial
	// log.Printf("Iniciando busca para query: %s", queryString)
	tui.PeersList.Clear()
	tui.originalPeers = []parser.Peer{}
	tui.filteredPeers = []parser.Peer{}

	// Iniciar spinner
	go func() {
		spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-done:
				return
			default:
				spin := spinners[i%len(spinners)]
				tui.App.QueueUpdateDraw(func() {
					tui.PeersList.SetTitle(fmt.Sprintf(" Loading %s ", spin))
				})
				time.Sleep(100 * time.Millisecond)
				i++
			}
		}
	}()

	// Buscar dados em goroutine
	go func() {
		newPeers := GetDataFromAPI(queryString)
		// log.Printf("Dados recebidos: %d peers", len(newPeers))

		// Parar spinner
		done <- true

		// Atualizar UI na thread principal
		tui.App.QueueUpdateDraw(func() {
			// Limpar lista atual
			tui.PeersList.Clear()
			// log.Printf("Lista limpa")

			// Atualizar peers
			tui.originalPeers = newPeers
			tui.filteredPeers = newPeers

			// Adicionar novos items
			for i, peer := range newPeers {
				tui.PeersList.AddItem(fmt.Sprintf("[%02d] %s", i+1, peer.PeerName), "", 0, func(index int) func() {
					return func() {
						tui.CurrentPeer = index
						tui.updateContent()
					}
				}(i))
			}
			// log.Printf("Itens adicionados à lista")

			// Atualizar título
			tui.PeersList.SetTitle(fmt.Sprintf(" Peers(%d) ", len(newPeers)))

			// Forçar atualização do conteúdo
			if len(newPeers) > 0 {
				tui.CurrentPeer = 0
				tui.updateContent()
			}

			// log.Printf("UI atualizada")
		})
	}()
}
