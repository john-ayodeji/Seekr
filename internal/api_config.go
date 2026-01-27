package internal

import (
	"github.com/john-ayodeji/Seekr/internal/database"
)

type ApiConfig struct {
	Port         int
	RabbitMQ_Url string
	Db           *database.Queries
}

var Cfg *ApiConfig
