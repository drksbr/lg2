package main

import (
	"github.com/drksbr/lg2/pkg/cli"
	"github.com/drksbr/lg2/pkg/tui"
)

func main() {
	// Obter os argumentos da CLI
	query := cli.ParseArgs()

	// // Exibir interface interativa
	t := tui.NewTUI(query)
	t.Start()
}
