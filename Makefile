# Define commands
GO_BUILD = go build -v -trimpath -o ./bin/app ./main.go
SWAGGER_GENERATE = swag init -g ./main.go -ot json -o ./public/docs --pd --parseInternal
NPM_INSTALL = cd resources/inertia && pnpm install
GO_MOD_TIDY = go mod tidy
VITE_DEV = cd resources/inertia && pnpm dev
VITE_PREVIEW = cd resources/inertia && pnpm preview
VITE_BUILD = cd resources/inertia && pnpm build
TEMPL_BUILD = templ generate
TEMPL_WATCH = templ generate --watch
AIR = air

# Ensure required tools are available
REQUIRED_EXECUTABLES = pnpm go air templ
K := $(foreach exec,$(REQUIRED_EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "$(exec) not found in PATH")))

# Targets
.PHONY: setup dev build build-assets build-app dev-assets dev-app generate-swagger lint preview

# Setup development environment
setup:
	$(NPM_INSTALL)
	$(GO_MOD_TIDY)

# Development mode - run frontend and backend concurrently
dev:
	@echo "Starting development environment..."
	@(trap 'kill 0' SIGINT; \
		$(VITE_DEV) & \
		$(TEMPL_WATCH) & \
		sleep 5 && $(AIR) & \
		wait)

# Build everything (frontend first, then backend)
build: build-assets build-app

# Build frontend assets
build-assets:
	$(VITE_BUILD)

# Build backend application
build-app:
	$(TEMPL_BUILD) && $(GO_BUILD)

# Generate Swagger documentation
swagger:
	$(SWAGGER_GENERATE)

# Lint the codebase
lint:
	cd resources/inertia && pnpm lint

# Preview the built application
preview:
	$(VITE_PREVIEW)
