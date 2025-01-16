package parser

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// AsPath representa um item de um caminho AS (Autonomous System).
type AsPath struct {
	AsNumber int    // Número do AS
	AsName   string // Nome do AS
}

// Peer representa os dados de um peer.
type Peer struct {
	PeerName string   // Nome do peer (extraído de "[ ]")
	Path     []AsPath // Caminho de ASNs
}

// ParsePeersFromHTML processa o HTML, divide os peers em partes e preenche as structs.
func ParsePeersFromHTML(htmlData string) ([]Peer, error) {
	var peers []Peer

	// Parse o HTML em uma árvore de nós
	doc, err := html.Parse(strings.NewReader(htmlData))
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar HTML: %v", err)
	}

	// Função auxiliar para encontrar o nó com ID bgpresults
	var findNodeByID func(*html.Node, string) *html.Node
	findNodeByID = func(n *html.Node, id string) *html.Node {
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if attr.Key == "id" && attr.Val == id {
					return n
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if found := findNodeByID(c, id); found != nil {
				return found
			}
		}
		return nil
	}

	// Encontrar o nó com ID "bgpresults"
	bgpResults := findNodeByID(doc, "bgpresults")
	if bgpResults == nil {
		return nil, fmt.Errorf("div com id='bgpresults' não encontrada")
	}

	// Função auxiliar para renderizar o HTML interno de um nó
	renderHTML := func(n *html.Node) string {
		var sb strings.Builder
		html.Render(&sb, n)
		return sb.String()
	}

	// Renderizar o HTML interno de bgpresults
	bgpResultsHTML := renderHTML(bgpResults)

	// Dividir o conteúdo da div em partes que iniciam com "unicast ["
	parts := strings.Split(bgpResultsHTML, "unicast [")

	// Processar cada parte
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Extrair PeerName
		peerNameEnd := strings.Index(part, "]")
		if peerNameEnd == -1 {
			continue
		}
		rawPeerName := strings.TrimSpace(part[:peerNameEnd])

		// Melhor tratamento do nome do peer
		peerName := strings.ReplaceAll(rawPeerName, "0000-00-00", "")
		peerName = strings.ReplaceAll(peerName, "\n", "")
		peerName = strings.TrimSpace(peerName)

		peerNameParts := strings.Split(peerName, "-")
		if len(peerNameParts) > 1 {
			peerName = strings.TrimSpace(peerNameParts[0]) + " - " + strings.TrimSpace(peerNameParts[1])
		}

		// Criar uma lista de AsPath
		var path []AsPath
		subDoc, err := html.Parse(strings.NewReader(part))
		if err != nil {
			return nil, fmt.Errorf("erro ao processar parte do peer: %v", err)
		}

		// Encontrar todos os elementos <abbr> na parte atual
		var findAbbr func(*html.Node)
		findAbbr = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "abbr" {
				var asNumber int
				var asName string
				for _, attr := range n.Attr {
					if attr.Key == "title" {
						asName = attr.Val
					}
				}
				if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
					asNumberStr := strings.TrimSpace(n.FirstChild.Data)
					asNumber, _ = strconv.Atoi(asNumberStr)
				}
				if asNumber != 0 && asName != "" {
					path = append(path, AsPath{
						AsNumber: asNumber,
						AsName:   asName,
					})
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				findAbbr(c)
			}
		}

		findAbbr(subDoc)

		// Adicionar o Peer à lista
		peers = append(peers, Peer{
			PeerName: peerName,
			Path:     path,
		})
	}

	return peers, nil
}
