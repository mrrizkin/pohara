package database

import (
	"context"
	"fmt"
	"math"

	"github.com/mrrizkin/pohara/config"
	"github.com/mrrizkin/pohara/internal/common/sql"
	"github.com/mrrizkin/pohara/internal/infrastructure/database/provider"
	"github.com/mrrizkin/pohara/internal/ports"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type GormDatabase struct {
	db *gorm.DB
}

type DatabaseDriver interface {
	DSN() string
	Connect(cfg *config.Database) (*gorm.DB, error)
}

type GormDatabaseDependencies struct {
	fx.In

	ConfigDB *config.Database
	Log      ports.Logger
}

type GormDatabaseResult struct {
	fx.Out

	Database *GormDatabase
}

var Module = fx.Module("gormdb",
	fx.Provide(NewGormDB),
	fx.Decorate(loadMigrationModels),
	fx.Provide(func(db *GormDatabase) ports.Database { return db }),
)

func AsGormMigration(model interface{}) any {
	return fx.Annotate(
		func() interface{} { return model },
		fx.ResultTags(`group:"gorm_migrate"`),
	)
}

func NewGormDB(
	lc fx.Lifecycle,
	deps GormDatabaseDependencies,
) (GormDatabaseResult, error) {
	var driver DatabaseDriver
	switch deps.ConfigDB.DRIVER {
	case "mysql", "mariadb", "maria":
		driver = provider.NewMysql(deps.ConfigDB)
	case "pgsql", "postgres", "postgresql":
		driver = provider.NewPostgres(deps.ConfigDB)
	case "sqlite", "sqlite3", "file":
		driver = provider.NewSqlite(deps.ConfigDB, deps.Log)
	default:
		return GormDatabaseResult{}, fmt.Errorf("unknown database driver: %s", deps.ConfigDB.DRIVER)

	}

	db, err := driver.Connect(deps.ConfigDB)
	if err != nil {
		return GormDatabaseResult{}, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			sqlDB, err := db.DB()
			if err != nil {
				return err
			}

			return sqlDB.Close()
		},
	})

	return GormDatabaseResult{
		Database: &GormDatabase{db},
	}, nil
}

func (g *GormDatabase) Create(value interface{}) error {
	return g.db.Create(value).Error
}

func (g *GormDatabase) First(dest interface{}, conds ...interface{}) error {
	return g.db.First(dest, conds...).Error
}

func (g *GormDatabase) Find(dest interface{}, conds ...interface{}) error {
	return g.db.Find(dest, conds...).Error
}

func (g *GormDatabase) Save(value interface{}) error {
	return g.db.Save(value).Error
}

func (g *GormDatabase) Delete(value interface{}, conds ...interface{}) error {
	return g.db.Delete(value, conds...).Error
}

func (g *GormDatabase) FindAll(
	dest interface{},
	p ports.Pagination,
	conds ...interface{},
) (*ports.FindAllResult, error) {
	var total int64
	query := g.db.Model(dest)

	if len(conds) > 0 {
		query = query.Where(conds[0], conds[1:]...)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if p.Sort.Valid && p.Sort.String != "" {
		query = query.Order(p.Sort.String)
	}

	if p.Limit.Valid {
		query = query.Limit(int(p.Limit.Int64))
	}

	if p.Offset.Valid {
		query = query.Offset(int(p.Offset.Int64))
	}

	if err := query.Find(dest).Error; err != nil {
		return nil, err
	}

	page := sql.Int64Nullable{
		Valid: false,
	}
	if p.Offset.Valid && p.Limit.Valid && p.Limit.Int64 != 0 {
		page.Valid = true
		page.Int64 = int64(math.Ceil(float64(p.Offset.Int64)/float64(p.Limit.Int64))) + 1
	}

	totalPage := sql.Int64Nullable{
		Valid: false,
	}
	if p.Offset.Valid && p.Limit.Valid && p.Limit.Int64 != 0 {
		totalPage.Valid = true
		totalPage.Int64 = int64(math.Ceil(float64(p.Offset.Int64)/float64(p.Limit.Int64))) + 1
	}

	return &ports.FindAllResult{
		Data:      dest,
		Total:     total,
		TotalPage: totalPage,
		Page:      page,
		Limit:     p.Limit,
	}, nil
}

func (g *GormDatabase) Count(model interface{}, conds ...interface{}) (int64, error) {
	var count int64
	query := g.db.Model(model)

	if len(conds) > 0 {
		query = query.Where(conds[0], conds[1:]...)
	}

	return count, query.Count(&count).Error
}

func (g *GormDatabase) Transaction(fn func(tx ports.Database) error) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		txDB := &GormDatabase{db: tx}
		return fn(txDB)
	})
}

func (g *GormDatabase) Raw(dest interface{}, statement string, values ...interface{}) error {
	return g.db.Raw(statement, values...).Scan(dest).Error
}

func (g *GormDatabase) Exec(statement string, values ...interface{}) error {
	return g.db.Exec(statement, values...).Error
}

type LoadMigrationModelsDependencies struct {
	fx.In

	Log    ports.Logger
	GormDB *GormDatabase
	Models []interface{} `group:"gorm_migrate"`
}

func loadMigrationModels(deps LoadMigrationModelsDependencies) *GormDatabase {
	deps.Log.Info("loading migration models", "count", len(deps.Models))
	deps.GormDB.db.AutoMigrate(deps.Models...)
	return deps.GormDB
}
