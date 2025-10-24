package types

import (
	"errors"
	"regexp"
	"time"

	sT "github.com/jsbento/chess-server-v4/cmd/services/sessions/types"
)

type User struct {
	Id        string    `json:"id" gorm:"primaryKey;index"`
	Username  string    `json:"username" gorm:"not null"`
	Email     string    `json:"email" gorm:"not null;unique"`
	Password  string    `json:"password" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`

	Sessions []*sT.Session `json:"sessions,omitempty" gorm:"foreignKey:UserId;references:Id"`
}

type CreateUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *CreateUser) Validate() error {
	if u.Username == "" {
		return errors.New("username is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	} else if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(u.Email) {
		return errors.New("invalid email address")
	}
	if u.Password == "" {
		return errors.New("password is required")
	} else if len(u.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	return nil
}

type UpdateUser struct {
	Id       string  `json:"id"`
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
}

func (u *UpdateUser) Validate() error {
	if u.Id == "" {
		return errors.New("id is required")
	}
	if u.Username != nil {
		if *u.Username == "" {
			return errors.New("username is required")
		}
	}
	if u.Email != nil {
		if *u.Email == "" {
			return errors.New("email is required")
		} else if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(*u.Email) {
			return errors.New("invalid email address")
		}
	}
	return nil
}
