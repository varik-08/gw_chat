package controllers

import (
	"context"

	"github.com/varik-08/gw_chat/config"
)

func GetAppFromContext(ctx context.Context) *config.App {
	if app, ok := ctx.Value("app").(*config.App); ok {
		return app
	}
	return nil
}
