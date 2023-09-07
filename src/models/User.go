package models

import (
	"errors"
	"strings"
	"time"
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

	u.format()

	return nil
}

func (u *User) validate(stage string) error {

	if u.Name == "" || u.Nick == "" || u.Email == "" || (u.Password == "" && stage == "create") {
		return errors.New("all fields are required")
	}

	return nil
}

func (u *User) format() {
	u.CreatedAt = time.Now()
	u.Name = strings.TrimSpace(u.Name)
	u.Nick = strings.TrimSpace(u.Nick)
	u.Email = strings.TrimSpace(u.Email)
}
