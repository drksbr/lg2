package tui

import (
	"fmt"
	"strings"

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
