package router

import (
	"context"
	"flag"
	"net/http"
)

type Router interface {
	ConfigureRouter(context.Context, flag.FlagSet) (http.Handler, error)
}
