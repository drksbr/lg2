package cli

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
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
	// Add default mask if not provided
	if !strings.Contains(prefix, "/") {
		ip := net.ParseIP(prefix)
		if ip != nil {
			if ip.To4() != nil {
				prefix += "/24"
			} else {
				prefix += "/48"
			}
		}
	}

	if !isValidPrefix(prefix) {
		fmt.Fprintf(os.Stderr, "Erro: prefixo inválido '%s'\n", prefix)
		os.Exit(1)
	}

	return prefix
}

func isValidPrefix(prefix string) bool {
	_, network, err := net.ParseCIDR(prefix)
	return err == nil && network != nil
}
