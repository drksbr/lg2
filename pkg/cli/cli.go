package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/drksbr/lg2/pkg/parser"
)

// ParseArgs lida com argumentos da linha de comando.
func ParseArgs() string {
	flag.Usage = func() {
		fmt.Println("Uso: lg <prefix>")
		fmt.Println("Exemplo IPv4: 1.1.1.0/24 ou 1.1.1.0")
		fmt.Println("Exemplo IPv6: 2001:db8::/32 ou 2001:db8::")
		flag.PrintDefaults()
	}

	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	prefix := args[0]

	network, err := parser.GetNetworkFromPrefix(prefix)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if network == nil {
		fmt.Println("Prefixo inválido")
		os.Exit(1)
	}

	return prefix
}
