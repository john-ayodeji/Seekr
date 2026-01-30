package search_handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/john-ayodeji/Seekr/internal"
	"github.com/john-ayodeji/Seekr/internal/database"
)

type SearchResponse struct {
	Query    string      `json:"query"`
	Limit    int32       `json:"limit"`
	Offset   int32       `json:"offset"`
	Count    int32       `json:"count"`
	Next     *string     `json:"next,omitempty"`
	Previous *string     `json:"previous,omitempty"`
	Results  interface{} `json:"results"`
}

func buildPageURL(r *http.Request, limit, offset int32) string {
	q := r.URL.Query()
	q.Set("limit", strconv.Itoa(int(limit)))
	q.Set("offset", strconv.Itoa(int(offset)))
	return r.URL.Path + "?" + q.Encode()
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func buildTsQuery(search string) string {
	words := strings.Fields(search)
	for i, w := range words {
		w = strings.TrimSpace(w)
		if w != "" {
			words[i] = w + ":*"
		}
	}
	return strings.Join(words, " & ")
}

func HandleSearch(w http.ResponseWriter, r *http.Request) {
	qv := r.URL.Query()

	Query := qv.Get("q")
	withSnippet := qv.Get("snippet") == "true"

	if Query == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "missing query",
		})
		return
	}

	// pagination defaults
	limit := int32(20)
	offset := int32(0)

	if l := qv.Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = int32(v)
		}
	}

	if o := qv.Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = int32(v)
		}
	}

	ctx := r.Context()

	var (
		results interface{}
		err     error
		count   int32
	)

	if withSnippet {
		var rows []database.SearchParsedHTMLWithSnippetRow
		rows, err = internal.Cfg.Db.SearchParsedHTMLWithSnippet(ctx, database.SearchParsedHTMLWithSnippetParams{
			Similarity: Query,
			ToTsquery:  buildTsQuery(Query),
			Limit:      limit,
			Offset:     offset,
		})
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}
		results = rows
		count = int32(len(rows))
	} else {
		var rows []database.SearchParsedHTMLRow
		rows, err = internal.Cfg.Db.SearchParsedHTML(ctx, database.SearchParsedHTMLParams{
			Similarity: Query,
			ToTsquery:  buildTsQuery(Query),
			Limit:      limit,
			Offset:     offset,
		})
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}
		results = rows
		count = int32(len(rows))
	}

	// next / previous logic
	var next *string
	var prev *string

	if count == limit {
		url := buildPageURL(r, limit, offset+limit)
		next = &url
	}

	if offset > 0 {
		prevOffset := offset - limit
		if prevOffset < 0 {
			prevOffset = 0
		}
		url := buildPageURL(r, limit, prevOffset)
		prev = &url
	}

	writeJSON(w, http.StatusOK, SearchResponse{
		Query:    Query,
		Limit:    limit,
		Offset:   offset,
		Count:    count,
		Next:     next,
		Previous: prev,
		Results:  results,
	})
}
