package delivery

import (
	"github.com/mrrizkin/pohara/internal/server"
)

func WebRouter(h *DashboardHandler) server.WebRouter {
	return server.NewWebRouter("/_/", func(r *server.Router) {
		r.Get("/", h.Index).Name("index")
	}, "dashboard")
}
