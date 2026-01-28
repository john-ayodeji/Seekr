package crawler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/john-ayodeji/Seekr/internal"
	amqp "github.com/rabbitmq/amqp091-go"
)

var StopWords = []string{
	"a", "about", "above", "after", "again", "against", "all", "am",
	"an", "and", "any", "are", "as", "at",

	"be", "because", "been", "before", "being", "below", "between",
	"both", "but", "by",

	"can",

	"did", "do", "does", "doing", "down", "during",

	"each",

	"few", "for", "from", "further",

	"had", "has", "have", "having", "he", "her", "here", "hers",
	"herself", "him", "himself", "his", "how",

	"i", "if", "in", "into", "is", "it", "its", "itself",

	"just",

	"me", "more", "most", "my", "myself",

	"no", "nor", "not", "now",

	"of", "off", "on", "once", "only", "or", "other", "our",
	"ours", "ourselves", "out", "over", "own",

	"same", "she", "should", "so", "some", "such",

	"than", "that", "the", "their", "theirs", "them",
	"themselves", "then", "there", "these", "they",
	"this", "those", "through", "to", "too",

	"under", "until", "up",

	"very",

	"was", "we", "were", "what", "when", "where", "which",
	"while", "who", "whom", "why", "with", "would",

	"you", "your", "yours", "yourself", "yourselves",
}

func ProcessTokens(conn *amqp.Connection) {
	StopWordsMap := BuildStopWordSet(StopWords)

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}

	msgs, err := internal.ConsumeFromQueue(conn, "parse.html.success")
	if err != nil {
		fmt.Println(err)
		return
	}

	type MSG struct {
		JobID       string `json:"Jobid"`
		URL         string `json:"Url"`
		Title       string `json:"Title"`
		Description string `json:"Description"`
		Headings    string `json:"Headings"`
		Paragraphs  string `json:"Paragraphs"`
		Links       string `json:"Links"`
	}

	type TokenData struct {
		JobID             string
		URL               string
		TitleTokens       []string
		DescriptionTokens []string
		HeadingTokens     []string
		ParagraphTokens   []string
	}

	for msg := range msgs {
		var p MSG
		if err := json.Unmarshal(msg.Body, &p); err != nil {
			fmt.Println(err)
			return
		}

		Title := TokenizeAndRemoveStopWords(p.Title, StopWordsMap)
		Description := TokenizeAndRemoveStopWords(p.Description, StopWordsMap)
		Headings := TokenizeAndRemoveStopWords(p.Headings, StopWordsMap)
		Paragraphs := TokenizeAndRemoveStopWords(p.Paragraphs, StopWordsMap)

		q := TokenData{
			JobID:             p.JobID,
			URL:               p.URL,
			TitleTokens:       Title,
			DescriptionTokens: Description,
			HeadingTokens:     Headings,
			ParagraphTokens:   Paragraphs,
		}

		if err := internal.Cfg.Db.MarkTokenized(context.Background(), p.URL); err != nil {
			fmt.Println(err)
			msg.Nack(false, true)
			continue
		}

		if err := internal.PublishToQueue(ch, internal.RabbitCfg.Exchange, "token.text.success", q); err != nil {
			fmt.Println(err)
			msg.Nack(false, true)
			continue
		}
	}
}
