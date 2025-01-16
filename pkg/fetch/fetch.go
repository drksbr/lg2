package fetch

import (
	"bytes"
	"io"
	"net/http"
)

var (
	debug bool   = false
	url   string = "https://bgp.tools/super-lg"
)

// GetLookingGlassData faz a requisição HTTP para o Looking Glass.
func GetLookingGlassData(query string) (string, error) {

	if debug {
		url = "http://localhost:3000/lg"
	}
	data := "q=" + query + "&asnmatch="
	headers := map[string]string{
		"accept":       "text/html",
		"content-type": "application/x-www-form-urlencoded",
		"origin":       "https://bgp.tools",
		"user-agent":   "Mozilla/5.0 (compatible; LG2/1.0)",
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		return "", err
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

	return string(body), nil
}
