package users

import (
	"fmt"
	"strings"
	"time"
)

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleManager Role = "manager"
	RoleWorker  Role = "worker"
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"update_at"`
	Status    string    `json:"status"`
}

func NewUser(name, lastName, email string, role Role) (*User, error) {
	u := &User{
		Name:     strings.TrimSpace(name),
		LastName: strings.TrimSpace(lastName),
		Email:    strings.TrimSpace(strings.ToLower(email)),
		Role:     role,
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) Validate() error {
	if u.Name == "" {
		return fmt.Errorf("name is required")
	}
	if u.LastName == "" {
		return fmt.Errorf("last_name is required")
	}
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}

	switch u.Role {
	case RoleAdmin, RoleManager, RoleWorker:
		return nil
	default:
		return fmt.Errorf("invalid role")
	}
}
