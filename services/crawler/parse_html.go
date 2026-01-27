package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/uuid"
	"github.com/john-ayodeji/Seekr/internal"
	"github.com/john-ayodeji/Seekr/internal/database"
	amqp "github.com/rabbitmq/amqp091-go"
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

func ProcessParseHTML(conn *amqp.Connection) {
	msgs, err := internal.ConsumeFromQueue(conn, "html_parser.jobs")
	if err != nil {
		fmt.Println(err)
		return
	}

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}

	type pgdata struct {
		ID    uuid.UUID `json:"ID"`
		JobID string    `json:"Jobid"`
		URL   string    `json:"Url"`
		HTML  string    `json:"Html"`
	}

	for msg := range msgs {
		var p pgdata
		if err := json.Unmarshal(msg.Body, &p); err != nil {
			fmt.Println(err)
			msg.Ack(false)
			continue
		}

		PageData, err := ParseHTML(p.HTML)
		if err != nil {
			fmt.Println(err)
			msg.Ack(false)
			continue
		}

		var dbData database.AddParsedHTMLParams
		dbData.ID = uuid.New()
		dbData.Jobid = p.JobID
		dbData.Url = p.URL
		dbData.Title = PageData.Title
		dbData.Description = PageData.Description
		for _, d := range PageData.Headings {
			dbData.Headings += fmt.Sprintf(", %v", d)
		}
		for _, e := range PageData.Paragraphs {
			dbData.Paragraphs += fmt.Sprintf(", %v", e)
		}
		for _, f := range PageData.Links {
			dbData.Links += fmt.Sprintf(", %v", f)
		}

		data, _ := internal.Cfg.Db.AddParsedHTML(context.Background(), dbData)
		_ = internal.Cfg.Db.MarkCrawled(context.Background(), p.URL)

		if err := internal.PublishToQueue(ch, internal.RabbitCfg.Exchange, "parse.html.success", data); err != nil {
			_ = internal.PublishToQueue(ch, internal.RabbitCfg.Exchange, "parse.html.success", data)
			_ = msg.Ack(false)
			continue
		}

		msg.Ack(false)
	}
}
