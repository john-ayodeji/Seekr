package crawler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTML struct {
	HTML string
}

var (
	ErrRedirect    = errors.New("redirect response")
	ErrForbidden   = errors.New("forbidden")
	ErrNotFound    = errors.New("not found")
	ErrRateLimited = errors.New("rate limited")
	ErrServerError = errors.New("server error")
)

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

	html, err := FetchHTML(res)
	if err != nil {
		return HTML{}, err
	}
	return html, nil
}

func FetchHTML(res *http.Response) (HTML, error) {
	switch res.StatusCode {

	case 200:
		htmlBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return HTML{}, err
		}
		return HTML{HTML: string(htmlBytes)}, nil

	case 301, 302, 307, 308:
		// follow redirect elsewhere, don’t index this URL
		return HTML{}, ErrRedirect

	case 401, 403:
		return HTML{}, ErrForbidden

	case 404, 410:
		return HTML{}, ErrNotFound

	case 429:
		return HTML{}, ErrRateLimited

	case 500, 502, 503, 504:
		return HTML{}, ErrServerError

	default:
		return HTML{}, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
}
