package main

import (
	"fmt"

	"github.com/drksbr/lg2/pkg/cli"
	"github.com/drksbr/lg2/pkg/fetch"
	"github.com/drksbr/lg2/pkg/parser"
	tui "github.com/drksbr/lg2/pkg/ui"
)

func main() {
	// Obter os argumentos da CLI
	query := cli.ParseArgs()

	// Buscar dados do Looking Glass
	htmlData, err := fetch.GetLookingGlassData(query)
	if err != nil {
		fmt.Printf("Erro ao buscar dados: %v\n", err)
		return
	}

	// Processar os resultados
	results, _ := parser.ParsePeersFromHTML(htmlData)

	// Print na tela os resultados
	// for _, peer := range results {
	// 	fmt.Printf("Peer Nome: %s\n", peer.PeerName)
	// 	fmt.Printf("Caminho de ASNs:\n")
	// 	for _, asPath := range peer.Path {
	// 		fmt.Printf("  ASN: %d, Nome: %s\n", asPath.AsNumber, asPath.AsName)
	// 	}
	// 	fmt.Println()
	// }

	// // Exibir interface interativa
	t := tui.NewTUI(results)
	t.Start()
}
