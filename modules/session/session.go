package session

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"go.uber.org/fx"

	dbConfig "github.com/mrrizkin/pohara/modules/database/config"
	"github.com/mrrizkin/pohara/modules/session/config"
	"github.com/mrrizkin/pohara/modules/session/provider"
)

type SessionProvider interface {
	Setup() (fiber.Storage, error)
}

type Store struct {
	*session.Store
	storage fiber.Storage
}

func NewSession(
	lc fx.Lifecycle,
	config *config.Config,
	dbConfig *dbConfig.Config,
) (*Store, error) {
	var driver SessionProvider

	switch config.Driver {
	case "database":
		driver = provider.NewDatabase(dbConfig)
	case "file":
		driver = provider.NewFile()
	case "redis", "valkey", "memory":
		driver = provider.NewMemory(config)
	default:
		driver = provider.NewFile()
	}

	storage, err := driver.Setup()
	if err != nil {
		return nil, err
	}

	store := Store{
		Store: session.New(session.Config{
			Storage:        storage,
			Expiration:     24 * time.Hour,
			KeyLookup:      fmt.Sprintf("cookie:%s_session_key", config.CookieName),
			CookieHTTPOnly: config.HttpOnly,
			CookieSecure:   config.Secure,
			CookieSameSite: config.SameSite,
		}),
		storage: storage,
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return store.Stop()
		},
	})

	return &store, nil
}

func (s *Store) Stop() error {
	return s.storage.Close()
}
