package controllers

import (
	"context"

	"github.com/varik-08/gw_chat/config"
	"github.com/varik-08/gw_chat/server/http/middlewares"
)

func GetAppFromContext(ctx context.Context) *config.App {
	if app, ok := ctx.Value(middlewares.AppKey).(*config.App); ok {
		return app
	}
	return nil
}
