package email


import (
	"metertronik/pkg/config"

	"github.com/sendgrid/sendgrid-go"
)

func NewSendgridClient(cfg *config.Config) *sendgrid.Client {
	apiKey := cfg.SendgridAPIKey
	
	return sendgrid.NewSendClient(apiKey)
}