package verification

import (
	"metertronik/pkg/config"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"metertronik/pkg/utils/template"
	"errors"
	emailClient"metertronik/pkg/verification/email"
)

func SendVerificationEmail(email string, code string) error {

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	from := mail.NewEmail(
		cfg.SendgridFromName,
		cfg.SendgridFromEmail,
	)
	
	to := mail.NewEmail("Metertronik New User", email)
	subject := "Your Metertronik Verification Code"
	htmlContent := template.VerificationEmailTemplate(code)
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)
	

	client := emailClient.NewSendgridClient(cfg)
	response, err := client.Send(message)


	if err != nil {
		return err
	}

	if response.StatusCode >= 400 {
		return errors.New("failed to send verification email, " + response.Body)
	}
	
	return nil
}