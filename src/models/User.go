package models

import (
	"errors"
	"friendfy-api/src/security"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

type User struct {
	ID        uint64    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Nick      string    `json:"nick,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (u *User) Prepare(stage string) error {
	if err := u.validate(stage); err != nil {
		return err
	}

	if err := u.format(stage); err != nil {
		return err
	}

	return nil
}

func (u *User) validate(stage string) error {

	if u.Name == "" || u.Nick == "" || u.Email == "" || (u.Password == "" && stage == "create") {
		return errors.New("all fields are required")
	}

	if err := checkmail.ValidateFormat(u.Email); err != nil {
		return errors.New("invalid email")
	}

	return nil
}

func (u *User) format(stage string) error {
	u.CreatedAt = time.Now()
	u.Name = strings.TrimSpace(u.Name)
	u.Nick = strings.TrimSpace(u.Nick)
	u.Email = strings.TrimSpace(u.Email)

	if stage == "create" {
		hashedPW, err := security.Hash(u.Password)
		if err != nil {
			return err
		}

		u.Password = string(hashedPW)
	}

	return nil
}
