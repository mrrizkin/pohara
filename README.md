# Pohara Starter Kit Documentation

## Introduction

Welcome to Pohara ("A Lot" in Sundanese), a powerful and modular Go starter kit for building modern web applications. Pohara is designed to be battery-included while maintaining clean architecture and modularity through dependency injection.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Tech Stack](#tech-stack)
3. [Project Structure](#project-structure)
4. [Creating Your First Module](#creating-your-first-module)
5. [Dependency Injection Guide](#dependency-injection-guide)
6. [Configuration](#configuration)
7. [Frontend Development](#frontend-development)
8. [Advanced Topics](#advanced-topics)

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Node.js 18 or higher
- PNPM package manager
- PostgreSQL (or your preferred database)

### Installation Steps

1. Clone the repository:

    ```bash
    git clone https://github.com/mrrizkin/pohara.git
    cd pohara
    ```

2. Install dependencies:

    ```bash
    go mod download
    pnpm install
    ```

3. Copy environment file:

    ```bash
    cp .env.example .env
    ```

4. Start development servers:

    ```bash
    # Terminal 1: Frontend
    pnpm dev

    # Terminal 2: Backend
    air
    ```

## Tech Stack

### Backend

- **Go & Gofiber**: Fast and lightweight web framework
- **Uber Go fx**: Dependency injection system
- **GORM**: Object-relational mapping
- **Air**: Live reload for development

### Frontend

- **Inertia.js**: Modern SPA architecture
- **Vite**: Frontend tooling
- **TailwindCSS**: Utility-first styling
- **PNPM**: Package management

## Project Structure

```
pohara/
├── app/                    # Application modules
├── bootstrap/              # Application bootstrapping
├── config/                 # Configuration files
├── models/                 # Database models
├── module/                 # Core system modules
├── public/                 # Static files
├── resources/              # Frontend resources
└── storage/                # Local storage
```

### Key Directories Explained

#### 1. `/app` - Application Modules

Each module follows a consistent structure:

```
app/user/
├── handler.go      # HTTP request handling
├── repository.go   # Data access
├── scheduler.go    # Scheduler
├── router.go       # Route definitions
├── service.go      # Business logic
└── user.go         # Module registration
```

#### 2. `/bootstrap` - Application Setup

Centralizes module registration and dependency injection:

```go
func App() *fx.App {
    return fx.New(
        config.New(),

        // Core services
        fx.Provide(
            logger.New,
            database.New,
            // ...
        ),

        // Application modules
        user.Module,
        welcome.Module,

        // Infrastructure
        models.Module,
        scheduler.Module,
        server.Module,
    )
}
```

## Creating Your First Module

### 1. Create Module Structure

```bash
mkdir -p app/mymodule
touch app/mymodule/{handler,repository,service,scheduler,router,mymodule}.go
```

### 2. Implement Components

#### Handler (handler.go)

```go
type MyHandler struct {
    log       *logger.Logger
    validator *validator.Validator
    service   *MyService
}

type HandlerDependencies struct {
    fx.In
    Log       *logger.Logger
    Validator *validator.Validator
    Service   *MyService
}

type HandlerResult struct {
    fx.Out
    Handler *MyHandler
}

func Handler(deps HandlerDependencies) HandlerResult {
    return HandlerResult{
        Handler: &MyHandler{
            log:       deps.Log,
            validator: deps.Validator,
            service:   deps.Service,
        },
    }
}
```

#### Repository (repository.go)

```go
type MyRepository struct {
    db *database.Database
}

type RepositoryDependencies struct {
    fx.In
    Db *database.Database
}

type RepositoryResult struct {
    fx.Out
    Repository *MyRepository
}

func Repository(deps RepositoryDependencies) RepositoryResult {
    return RepositoryResult{
        Repository: &MyRepository{
            db: deps.Db,
        },
    }
}
```

#### Service (service.go)

```go
type MyService struct {
    repo *MyRepository
}

type ServiceDependencies struct {
    fx.In
    Repository *MyRepository
}

type ServiceResult struct {
    fx.Out
    Service *MyService
}

func Service(deps ServiceDependencies) ServiceResult {
    return ServiceResult{
        Service: &MyService{
            repo: deps.Repository,
        },
    }
}
```

#### Scheduler (scheduler.go)

```go
type MyScheduler struct {
    log *logger.Logger
}

type SchedulerDependencies struct {
    fx.In
    Log *logger.Logger
}

type SchedulerResult struct {
    fx.Out
    Scheduler *MyScheduler
}

func Scheduler(deps SchedulerDependencies) SchedulerResult {
    return SchedulerResult{
        Scheduler: &MyScheduler{
            log: deps.Log,
        },
    }
}

func (ms *MyScheduler) Schedule(s *scheduler.Scheduler) {
    s.Add("0 0 * * *", func() {
        ms.log.Info("Running daily backup...")
    })
}
```

#### Router (router.go)

```go
func webRouter(h *MyHandler) server.WebRouter {
    return server.WebRouter{
        Prefix: "/mymodule",
        Router: func(r fiber.Router) {
            r.Get("/", h.Index)
            r.Get("/:id", h.Show)
        },
        Name: "mymodule.",
    }
}

func apiRouter(h *MyHandler) server.ApiRouter {
    return server.ApiRouter{
        Prefix: "/mymodule",
        Router: func(r fiber.Router) {
            r.Get("/", h.List)
            r.Post("/", h.Create)
        },
        Name:    "api.mymodule.",
        Version: "v1",
    }
}
```

#### Module Registration (mymodule.go)

```go
var Module = fx.Module("mymodule",
    fx.Provide(
        Repository,
        Service,
        Handler,
        Scheduler,
        server.AsWebRouter(webRouter),
        server.AsApiRouter(apiRouter),
        scheduler.AsSchedule(func(ms *MyScheduler) scheduler.Schedule {
            return ms
        }),
    ),
)
```

### 3. Register in Bootstrap

Add your module to `bootstrap/bootstrap.go`:

```go
func App() *fx.App {
    return fx.New(
        // ... existing modules
        mymodule.Module,
        // ...
    )
}
```

## Dependency Injection Guide

### Key Concepts

1. **Dependencies (fx.In)**

    - Used to declare what a component needs
    - Automatically injected by fx

2. **Results (fx.Out)**

    - Declares what a component provides
    - Available for injection into other components

3. **Modules (fx.Module)**
    - Groups related components
    - Provides namespace isolation

### Best Practices

1. **Keep Dependencies Explicit**

    ```go
    type Dependencies struct {
        fx.In
        Logger    *logger.Logger
        Database  *database.Database
    }
    ```

2. **Use Constructor Functions**

    ```go
    func New(deps Dependencies) Result {
        return Result{
            Component: &Component{
                logger: deps.Logger,
                db:     deps.Database,
            },
        }
    }
    ```

3. **Module Independence**
    - Each module should be self-contained
    - Use interfaces for cross-module communication

## Configuration

### Configuration Files

```go
config/
├── app.go
├── database.go
├── session.go
├── config.go
└── load.go
```

### Configuration Definition

Configuration structs use `env` tags to map environment variables:

```go
type Database struct {
    DRIVER       string `env:"DB_DRIVER,default=sqlite"`
    HOST         string `env:"DB_HOST"`
    PORT         int    `env:"DB_PORT,default=5432"`
    NAME         string `env:"DB_NAME"`
    USERNAME     string `env:"DB_USERNAME,default=root"`
    PASSWORD     string `env:"DB_PASSWORD,default=root"`
    SSLMODE      string `env:"DB_SSLMODE,default=disable"`
    AUTO_MIGRATE bool   `env:"DB_AUTO_MIGRATE,default=true"`
}
```

### Tag Format

- Basic mapping: `env:"ENV_VAR_NAME"`
- With default: `env:"ENV_VAR_NAME,default=value"`
- Required field: `env:"ENV_VAR_NAME,required"`

### Configuration Loading

The configuration system is initialized in `config.go`:

```go
func New() fx.Option {
    configs := []interface{}{
        &App{},
        &Database{},
        &Session{},
    }

    for _, config := range configs {
        load(config)
    }

    return fx.Supply(configs...)
}
```

### Environment Variables

Create a `.env` file in your project root:

```env
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=myapp
DB_USERNAME=postgres
DB_PASSWORD=secret
DB_SSLMODE=disable
DB_AUTO_MIGRATE=true
```

## Frontend Development

### Directory Structure

```
resources/
├── js/
│   ├── pages/
│   └── main.jsx
├── css/
│   └── index.css
└── views/
    └── root.html
```

### Creating Pages

1. Add new page component in `resources/js/pages/`
2. Register route in your module's router
3. Use Inertia.js for navigation

### Building Assets

```bash
# Development
pnpm dev

# Production
pnpm build
```

## Advanced Topics

### Custom Middleware

Add middleware in your module's router:

```go
func apiRouter(h *MyHandler) server.ApiRouter {
    return server.ApiRouter{
        Router: func(r fiber.Router) {
            r.Use(middleware.Auth())
            // routes...
        },
    }
}
```

### Database Migrations

1. Create model in `/models`
2. Register in `models/models.go`
3. Migrations run automatically on startup

### Task Scheduling

Use the scheduler module for periodic tasks:

```go
func Schedule(scheduler *scheduler.Scheduler) {
    scheduler.Every(1).Hour().Do(myTask)
}
```

## Best Practices

### Configuration

1. **Environment Variables**

    - Use meaningful names
    - Group related variables
    - Document default values
    - Mark security-sensitive fields as required

2. **Type Safety**

    - Use appropriate types in structs
    - Add validation where needed
    - Handle conversion errors gracefully

3. **Organization**
    - Group related configurations
    - Keep config structs focused
    - Document special requirements

### Scheduling

1. **Task Design**

    - Keep tasks idempotent
    - Include error handling
    - Log task execution
    - Avoid long-running tasks

2. **Schedule Planning**

    - Consider timezone effects
    - Avoid resource contention
    - Plan for failure recovery
    - Document schedule reasoning

3. **Maintenance**
    - Monitor task execution
    - Review schedules periodically
    - Keep tasks modular
    - Document dependencies

## Troubleshooting

### Common Issues

1. **Module Not Loading**

    - Check module registration in bootstrap
    - Verify dependency injection setup

2. **Database Connection Issues**

    - Verify environment variables
    - Check database credentials

3. **Frontend Asset Issues**
    - Run `pnpm install`
    - Clear Vite cache
    - Check for JavaScript errors

### Configuration Issues

1. **Missing Environment Variables**

    - Check .env file exists
    - Verify variable names match
    - Check for typos in tags

2. **Type Conversion Errors**
    - Verify value formats
    - Check default values
    - Review type definitions

### Scheduler Issues

1. **Tasks Not Running**

    - Verify cron expression
    - Check scheduler registration
    - Review log output
    - Confirm module initialization

2. **Performance Problems**
    - Review task duration
    - Check resource usage
    - Monitor concurrent execution
    - Verify logging impact

### Getting Help

- Open issues on GitHub
- Check existing documentation
- Join the community discord

## Contributing

1. Fork the repository
2. Create feature branch
3. Commit changes
4. Push to branch
5. Create Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
