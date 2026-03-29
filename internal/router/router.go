package router

import (
	"context"
	"net/http"

	config "github.com/Pklerik/gophKeep/internal/config/server"
)

type Router interface {
	ConfigureRouter(context.Context, config.Config) (http.Handler, error)
}
