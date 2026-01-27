package sitemap_processor

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Sitemap struct {
	//XMLName xml.Name `xml:"urlset"`
	UrlSet []URL `xml:"url"`
}

type URL struct {
	Loc      string  `xml:"loc"`
	Priority float64 `xml:"priority"`
}

func ParseSitemap(url string) (Sitemap, string, error) {
	var sitemap Sitemap

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("error from req: %v", err)
		return Sitemap{}, "", err
	}

	req.Header.Set("User-Agent", "SeekrBot/1.0 (+https://seekr.tech/bot-info")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("From", "crawler@seekr.tech")
	req.Header.Set("Cache-Control", "no-cache")

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("client.Do error: %v\n", err)
		return Sitemap{}, "", err
	}
	defer res.Body.Close()

	XMLBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return Sitemap{}, "", err
	}

	if err := xml.Unmarshal(XMLBytes, &sitemap); err != nil {
		fmt.Println(err)
		return Sitemap{}, "", err
	}

	return sitemap, uuid.New().String(), nil
}
