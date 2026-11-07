package app

import (
	"context"
	"errors"
)

type App struct{}

func NewApp(name string) *App {
	return &App{}
}

func (a *App) Run(ctx context.Context, args []string) error {
	return errors.New("not implemented")
}
