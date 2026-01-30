package main

import (
	"net/http"

	"github.com/john-ayodeji/Seekr/handlers/search_handler"
	"github.com/john-ayodeji/Seekr/handlers/sitemap_handler"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/sitemap", sitemap_handler.HandleSitemapSubmit)
	mux.HandleFunc("GET /api/v1/search", search_handler.HandleSearch)
}
