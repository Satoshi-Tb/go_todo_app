package service

import (
	"context"
	"fmt"

	"github.com/Satoshi-Tb/go_todo_app/entity"
	"github.com/Satoshi-Tb/go_todo_app/store"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUser struct {
	DB   store.Execer
	Repo UserRegister
}

func (r *RegisterUser) RegisterUser(ctx context.Context, name, passord, role string) (*entity.User, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(passord), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: ")
	}

	u := &entity.User{
		Name:     name,
		Password: string(pw),
		Role:     role,
	}

	if err := r.Repo.RegisterUser(ctx, r.DB, u); err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}

	return u, nil
}
