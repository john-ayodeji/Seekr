package sitemap_handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/john-ayodeji/Seekr/internal"
	"github.com/john-ayodeji/Seekr/services/sitemap_processor"
	"github.com/john-ayodeji/Seekr/utils"
)

var RBcfg internal.RabbitConfig

type sitemapRequest struct {
	URL string `json:"sitemap_url"`
}

type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"msg"`
}

type successResponse struct {
	Message string `json:"msg"`
}

func HandleSitemapSubmit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	var req sitemapRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Status:  "error",
			Message: "invalid request body",
		})
		return
	}

	if req.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Status:  "error",
			Message: "sitemap_url is required",
		})
		return
	}

	if !strings.HasSuffix(req.URL, "sitemap.xml") {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Status:  "error",
			Message: "url should contain '/sitemap.xml'",
		})
		return
	}

	normalizedURL, err := utils.NormalizeURL(req.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Status:  "error",
			Message: "invalid sitemap url",
		})
		return
	}

	sitemap, err := sitemap_processor.ParseSitemap(normalizedURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Status:  "error",
			Message: "failed to parse sitemap",
		})
		return
	}

	conn := RBcfg.Connection
	if conn == nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Status:  "error",
			Message: "rabbitmq connection not initialized",
		})
		return
	}

	ch, err := conn.Channel()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Status:  "error",
			Message: "failed to open rabbitmq channel",
		})
		return
	}
	defer ch.Close()

	for _, sm := range sitemap.UrlSet {
		nURL, err := utils.NormalizeURL(sm.Loc)
		if err != nil {
			continue
		}

		payload := struct {
			URL string `json:"url"`
		}{
			URL: nURL,
		}

		if err := internal.PublishToQueue(
			ch,
			RBcfg.Exchange,
			"url.fetch.success",
			payload,
		); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(errorResponse{
				Status:  "error",
				Message: "failed to enqueue sitemap urls",
			})
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(successResponse{
		Message: "your website has been added to the queue",
	})
}
