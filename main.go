package main

import (
	"github.com/mrrizkin/pohara/bootstrap"
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
	app := bootstrap.App()
	if err := app.Err(); err != nil {
		panic(err)
	}

	app.Run()
}
