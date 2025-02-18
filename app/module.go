package app

import (
	"github.com/mrrizkin/pohara/app/config"
	"github.com/mrrizkin/pohara/app/http/controllers"
	"github.com/mrrizkin/pohara/app/http/middleware"
	"github.com/mrrizkin/pohara/app/repository"
	"github.com/mrrizkin/pohara/app/routes"
	"github.com/mrrizkin/pohara/app/service"
	"go.uber.org/fx"
)

var Module = fx.Module("app",
	config.Provide,
	controllers.Provide,
	middleware.Provide,
	repository.Provide,
	service.Provide,
	routes.Provide,
)
