package dto

import (
	"fmt"
	"strings"
)

type LoginPayload struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r *LoginPayload) Validate() error {
	if len(strings.TrimSpace(r.Login)) < 2 {
		return fmt.Errorf("login must be at least 2 characters long")
	}

	if len(strings.TrimSpace(r.Password)) < 4 {
		return fmt.Errorf("password must be at least 6 characters long")
	}

	return nil
}
