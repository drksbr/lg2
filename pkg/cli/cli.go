package cli

import (
	"fmt"
	"os"

	"github.com/drksbr/lg2/pkg/config"
	"github.com/drksbr/lg2/pkg/tui"
	"github.com/spf13/cobra"
)

var (
	Banner = `
    ▒▒  ▒▓ ▒▒▒▒▓  ▓▒▒▒▒▒
   ▓  ▓ ▒ ▒ ▒   ▓▒    ▒ 
  ▓   ▓▓   ▒ ▒  ▒ ▒  ▒     Multiglass %s // (c) 2024
 ▓▒▒▒▓▒ ▓▒▒▒▓ ▒▒▒  ▒▒      https://github/drksbr/lg2

-> A TUI interface for nlnog.net looking glass, essential for querying and analyzing BGP routes, providing insights on AS-PATH, origin validations, and communities in an interactive and intuitive manner. Compatible with Windows, Mac, and Linux, it offers high performance and an optimized interface for network engineers and administrators.

	IPv4 usage example: lg 1.1.1.0/24
	IPv6 usage example: lg 2001:db8::/32`

	showVersion bool

	rootCmd = &cobra.Command{
		Use:   "lg [flags] [prefix]",
		Short: "Looking Glass CLI for querying BGP prefixes",
		Long:  fmt.Sprintf(Banner, config.Version),
		Run:   run,
	}
)

func run(cmd *cobra.Command, args []string) {
	if showVersion {
		fmt.Printf("lg version %s\n", config.Version)
		return
	}

	if err := searchPrefix(cmd, args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func searchPrefix(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("é necessário informar um prefixo")
	}

	// // Exibir interface interativa
	t := tui.NewTUI(args[0])
	t.Start()

	return nil
}

func Execute() {
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "mostra a versão do lg")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
