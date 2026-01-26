package crawler

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
	Title       string
	Description string
	Headings    []string
	Paragraphs  []string
	Links       []string
}

func ParseHTML(html string) (*PageData, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	page := &PageData{}

	page.Title = strings.TrimSpace(doc.Find("title").Text())

	if desc, exists := doc.Find(`meta[name="description"]`).Attr("content"); exists {
		page.Description = strings.TrimSpace(desc)
	}

	doc.Find("h1, h2, h3").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			page.Headings = append(page.Headings, text)
		}
	})

	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			page.Paragraphs = append(page.Paragraphs, text)
		}
	})

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			href = strings.TrimSpace(href)
			if href != "" {
				page.Links = append(page.Links, href)
			}
		}
	})

	return page, nil
}
