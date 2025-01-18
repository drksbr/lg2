package tui

import (
	"fmt"
	"strings"

	"github.com/drksbr/lg2/pkg/parser"
)

func buildPeerDetails(peer *parser.Peer) string {
	var details strings.Builder

	// Build the header string with prefix and peer
	details.WriteString(fmt.Sprintf("[::b]Info:[::-] %s / %s\n\n", peer.Prefix, peer.PeerName))

	// Build the details string
	details.WriteString(fmt.Sprintf("[::b]AS-PATH:[::-] %s\n\n", formatASPath(peer.AsPath)))

	// Append AS path details
	details.WriteString("[::b]Sequential:[::-]\n")
	for i, as := range peer.AsPath {
		details.WriteString(fmt.Sprintf("     [%02d] |%s| [::b]AS%d[::-] (%s)\n", i+1, as.Country, as.AsNumber, as.AsName))
	}
	details.WriteString("\n")

	if len(peer.Communities) > 0 {
		details.WriteString("[::b]Communities:[::-]\n     ")
		for i, community := range peer.Communities {
			if i > 0 {
				details.WriteString(" | ")
			}
			if i > 0 && i%4 == 0 {
				details.WriteString("\n     ")
			}
			details.WriteString(community)
		}
		details.WriteString("\n\n")
	} 

	// Append MED
	details.WriteString("[::b]MED:[::-] ")
	details.WriteString(peer.Med)
	details.WriteString("\n\n")

	// Append Last Update
	details.WriteString("[::b]Last Update:[::-] ")
	details.WriteString(peer.LastUpdate)


	return details.String()
}
