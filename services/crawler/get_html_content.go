package crawler

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTML struct {
	HTML string
}

func GetHTML(url string) (HTML, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return HTML{}, err
	}
	req.Header.Set("User-Agent", "SeekrBot/1.0 (+https://seekr.tech/bot-info")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("From", "crawler@seekr.tech")
	req.Header.Set("Cache-Control", "no-cache")

	client := http.Client{
		Timeout: 20 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return HTML{}, err
	}

	if res.StatusCode == 200 || res.StatusCode == 301 || res.StatusCode == 302 {
		htmlBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return HTML{}, err
		}
		return HTML{
			HTML: string(htmlBytes),
		}, nil
	} else {
		return HTML{}, fmt.Errorf("crawling is not allowed on this page")
	}
}
