package validator

import (
	"errors"
	"regexp"
	"strings"
)

var controllerIDRe = regexp.MustCompile(`^[A-Za-z0-9_-]{3,64}$`)

func ValidateControllerID(id string) error {
	if id == "" {
		return errors.New("invalid id: id wajib diisi")
	}

	if strings.TrimSpace(id) != id {
		return errors.New("invalid id: id tidak boleh mengandung spasi/whitespace (termasuk di awal/akhir)")
	}
	if strings.ContainsAny(id, " \t\n\r\v\f") {
		return errors.New("invalid id: id tidak boleh mengandung spasi/whitespace")
	}

	if !controllerIDRe.MatchString(id) {
		return errors.New("invalid id: format id tidak valid (hanya boleh A-Z a-z 0-9, '_' atau '-', panjang 3-64)")
	}
	return nil
}
