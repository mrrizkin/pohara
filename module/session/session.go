package session

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"go.uber.org/fx"

	"github.com/mrrizkin/pohara/config"
	"github.com/mrrizkin/pohara/module/session/provider"
)

type SessionProvider interface {
	Setup() (fiber.Storage, error)
}

type Session struct {
	*session.Store
	storage fiber.Storage
}

func New(
	lc fx.Lifecycle,
	config *config.App,
	dbConfig *config.Database,
	sessionConfig *config.Session,
) (*Session, error) {
	var driver SessionProvider

	switch sessionConfig.DRIVER {
	case "database":
		driver = provider.NewDatabase(dbConfig)
	case "file":
		driver = provider.NewFile()
	case "redis", "valkey", "memory":
		driver = provider.NewMemory(sessionConfig)
	default:
		driver = provider.NewFile()
	}

	storage, err := driver.Setup()
	if err != nil {
		return nil, err
	}

	sess := Session{
		Store: session.New(session.Config{
			Storage:        storage,
			Expiration:     24 * time.Hour,
			KeyLookup:      fmt.Sprintf("cookie:%s_session_key", config.APP_NAME),
			CookieHTTPOnly: sessionConfig.HTTP_ONLY,
			CookieSecure:   sessionConfig.SECURE,
			CookieSameSite: sessionConfig.SAME_SITE,
		}),
		storage: storage,
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return sess.Stop()
		},
	})

	return &sess, nil
}

func (s *Session) Stop() error {
	return s.storage.Close()
}
