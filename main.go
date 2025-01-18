package main

import (
	"fmt"

	"github.com/drksbr/lg2/pkg/cli"
	"github.com/drksbr/lg2/pkg/fetch"
	"github.com/drksbr/lg2/pkg/parser"
	"github.com/drksbr/lg2/pkg/tui"
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
	results, _ := parser.ParseHTML(htmlData, query)

	// // Exibir interface interativa
	t := tui.NewTUI(results)
	t.Start()
}
