package service

import (
	"context"
	"fmt"

	"github.com/Satoshi-Tb/go_todo_app/store"
)

type LoginTask struct {
	DB             store.Queryer
	Repo           UserGetter
	TokenGenerator TokenGenerator
}

func (lt *LoginTask) Login(ctx context.Context, name, pw string) (string, error) {
	u, err := lt.Repo.GetUser(ctx, lt.DB, name)
	if err != nil {
		return "", fmt.Errorf("failed to list: %w", err)
	}

	if err = u.ComparePassword(pw); err != nil {
		return "", fmt.Errorf("wrong password: %w", err)
	}

	jwt, err := lt.TokenGenerator.GenerateToken(ctx, *u)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}
	return string(jwt), nil
}
