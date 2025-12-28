package template

import (
	_ "embed"
	"strings"
)

//go:embed verification-page.html
var verificationTemplate string

func VerificationEmailTemplate(code string) string {
	// Ganti placeholder {{CODE}} dengan kode verifikasi yang sebenarnya
	html := strings.ReplaceAll(verificationTemplate, "{{CODE}}", code)
	return html
}
