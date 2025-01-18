package main

import (
	"os"

	"github.com/mrrizkin/pohara/bootstrap"
	"go.uber.org/fx"
)

// @title						Pohara API
// @version						1.0
// @description					Pohara API provides a comprehensive set of endpoints for managing user data and related operations.
// @termsOfService				https://www.example.com/terms/
// @contact.name				API Support Team
// @contact.url					https://www.example.com/support
// @contact.email				support@example.com
// @license.name				MIT License
// @license.url					https://opensource.org/licenses/MIT
// @host						localhost:3000
// @BasePath					/api/v1
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						pohara-api-token
// @externalDocs.description	OpenAPI Specification
// @externalDocs.url			https://swagger.io/specification/
func main() {
	var app *fx.App

	if len(os.Args) > 1 && os.Args[1] == "cli" {
		os.Args = append(os.Args[:1], os.Args[2:]...)
		app = bootstrap.Console()
	} else {
		app = bootstrap.App()
	}

	if err := app.Err(); err != nil {
		panic(err)
	}

	app.Run()
}
