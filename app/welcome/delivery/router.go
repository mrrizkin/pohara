package delivery

import (
	"github.com/mrrizkin/pohara/internal/server"
)

func WebRouter(h *WelcomeHandler) server.WebRouter {
	return server.NewWebRouter("/", func(r *server.Router) {
		r.Get("/", h.Index).Name("index")
	}, "welcome")
}
