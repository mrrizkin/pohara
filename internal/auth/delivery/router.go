package delivery

import (
	"github.com/mrrizkin/pohara/internal/server"
)

func WebRouter(h *AuthHandler) server.WebRouter {
	return server.NewWebRouter("/_/auth", func(r *server.Router) {
		r.Get("/login", h.LoginPage).Name("login")
	}, "auth.")
}

func ApiRouter(h *AuthHandler) server.ApiRouter {
	return server.NewApiRouter("v1", "/auth", func(r *server.Router) {
		r.Get("/login", h.Login).Name("login")
	}, "api.auth.")
}
