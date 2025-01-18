package parser

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type AsPath struct {
	AsNumber int    // Número do AS
	AsName   string // Nome do AS
	Country  string // Sigla do país
}

type Peer struct {
	PeerName          string   // Nome do peer
	AsPath            []AsPath // Caminho de ASNs
	OriginValidation  string   // Estado de validação de origem
	AspaValidation    string   // Estado de validação ASPA
	OnlyToCustomerOTC string   // Informações de "Only To Customer"
	Origin            string   // Origem
	Med               string   // MED (Multi Exit Discriminator)
	LastUpdate        string   // Última atualização
	Communities       []string // Comunidades
	Prefix            string   // Prefixo
}

// ParseHTML parses the updated HTML format and extracts peer information
func ParseHTML(htmlData string, prefix string) ([]Peer, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlData))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	var peers []Peer

	doc.Find("div.peername").Each(func(index int, peerNode *goquery.Selection) {
		peer := Peer{}

		// Set Prefix
		peer.Prefix = prefix

		// Extract Peer Name
		peer.PeerName = strings.Split(strings.TrimSpace(peerNode.Find(".me-auto").Text()), " ")[1]

		// Navigate to the corresponding table for the peer
		peerTable := peerNode.NextFiltered("table")

		peerTable.Find("tr").Each(func(_ int, row *goquery.Selection) {
			header := strings.TrimSpace(row.Find("td:first-child").Text())
			data := row.Find("td")

			switch header {
			case "AS-Path":
				var asPath []AsPath
				data.Find("button").Each(func(_ int, btn *goquery.Selection) {
					number, _ := strconv.Atoi(btn.Find("a.whois.asn").Text())
					name := btn.AttrOr("title", "")

					// Extract Country from the name if available
					var country string
					if parts := strings.Split(name, ","); len(parts) > 1 {
						country = strings.ToUpper(strings.TrimSpace(parts[len(parts)-1])[:2])
						name = strings.TrimSpace(parts[0])
					}

					// Discard anything after <br> in the name
					if brIndex := strings.Index(name, "<br>"); brIndex != -1 {
						name = strings.TrimSpace(name[:brIndex])
					}

					asPath = append(asPath, AsPath{AsNumber: number, AsName: name, Country: country})
				})
				peer.AsPath = asPath

			case "Origin validation state":
				peer.OriginValidation = data.Text()

			case "ASPA validation state":
				peer.AspaValidation = data.Text()

			case "Only To Customer (OTC)":
				peer.OnlyToCustomerOTC = data.Text()

			case "Origin":
				peer.Origin = data.Text()

			case "MED":
				peer.Med = strings.TrimPrefix(data.Text(), "MED")

			case "Last update":
				peer.LastUpdate = strings.TrimSpace(strings.TrimPrefix(data.Text(), "Last update"))

			case "Communities":
				var communities []string
				data.Find("button").Each(func(_ int, btn *goquery.Selection) {
					communities = append(communities, btn.Text())
				})
				peer.Communities = communities
			}
		})

		peers = append(peers, peer)
	})

	return peers, nil
}

func GetNetworkFromPrefix(prefix string) (*net.IPNet, error) {
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

	_, network, err := net.ParseCIDR(prefix)
	if err != nil {
		return nil, fmt.Errorf("invalid prefix: %v", err)
	}
	return network, nil
}
