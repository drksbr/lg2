package fetch

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	debug      bool   = false
	saveSample bool   = true
	url        string = "https://lg.ring.nlnog.net/prefix"
)

// Func to save html data to a file
func saveHTML(htmlData string, query string) {
	// Create samples directory if it doesn't exist
	if err := os.MkdirAll("samples", 0755); err != nil {
		log.Fatal(err)
	}
	// timestamp := time.Now().Format("20060102-150405")
	f, err := os.Create(fmt.Sprintf("samples/sample-%s.html", query))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err = f.WriteString(htmlData)
	if err != nil {
		log.Fatal(err)
	}
}

// GetLookingGlassData makes an HTTP request to the Looking Glass.
func GetLookingGlassData(query string) (string, error) {
	if debug {
		url = "http://localhost:3000/sample"
	}

	client := &http.Client{}
	reqURL := url + "?q=" + query + "&match=exact&peer=all"
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", err
	}

	headers := map[string]string{
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"Accept-Language":           "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7,es;q=0.6",
		"Cache-Control":             "max-age=0",
		"Connection":                "keep-alive",
		"Sec-Fetch-Dest":            "document",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "none",
		"Sec-Fetch-User":            "?1",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
		"sec-ch-ua":                 "\"Google Chrome\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"",
		"sec-ch-ua-mobile":          "?0",
		"sec-ch-ua-platform":        "\"macOS\"",
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Save sample HTML data to a file
	if saveSample {
		// Replace incompatible characters in query with '_'
		safeQuery := strings.NewReplacer(
			"/", "_",
			"\\", "_",
			":", "_",
			"*", "_",
			"?", "_",
			"\"", "_",
			"<", "_",
			">", "_",
			"|", "_",
		).Replace(query)
		saveHTML(string(body), safeQuery)
	}

	return string(body), nil
}
