package neoweb

import (
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/modules/neoweb/inertia"
	"github.com/mrrizkin/pohara/modules/neoweb/templ"
	"github.com/mrrizkin/pohara/modules/neoweb/vite"
)

var Module = fx.Module("neoweb",
	fx.Provide(inertia.New, vite.New, templ.New),
)
