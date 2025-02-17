package migration

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"go.uber.org/fx"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/mrrizkin/pohara/app/config"
	"github.com/mrrizkin/pohara/modules/database/db"
	"github.com/mrrizkin/pohara/modules/database/migration/dialect"
	"github.com/mrrizkin/pohara/modules/database/model"
	"github.com/mrrizkin/pohara/modules/database/repository"
)

//go:embed template/*
var templateFS embed.FS

type Migrator struct {
	mhRepo     *repository.MigrationHistoryRepository
	config     *config.Config
	tmpl       *template.Template
	migrations []Migration
}

type MigratorDeps struct {
	fx.In

	MHRepo *repository.MigrationHistoryRepository
	Db     *db.Database
	Config *config.Config
}

func NewMigrator(deps MigratorDeps) *Migrator {
	err := deps.Db.AutoMigrate(model.MigrationHistory{})
	if err != nil {
		panic(fmt.Sprintf("failed to migrate the __migration_history__: %v", err))
	}

	tmpl, err := template.New("template").ParseFS(templateFS, "template/*.go.tmpl")
	if err != nil {
		panic(fmt.Sprintf("failed to parse template: %v", err))
	}

	return &Migrator{
		mhRepo: deps.MHRepo,
		config: deps.Config,
		tmpl:   tmpl,
	}
}

func (m *Migrator) AddMigration(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

func (m *Migrator) Migrate() error {
	dialect, err := m.getDialect()
	if err != nil {
		return err
	}

	schema := NewSchema(dialect)
	migrations := make(map[string]Migration)
	histories, err := m.mhRepo.GetMigrationHistory("id ASC")
	if err != nil {
		return fmt.Errorf("failed to get migration history: %v", err)
	}

	executedMigrations := make(map[string]model.MigrationHistory)
	for _, h := range histories {
		executedMigrations[h.ID] = h
	}

	for _, migration := range m.migrations {
		id := migration.ID()
		if _, exists := executedMigrations[id]; !exists {
			migrations[id] = migration
		}
	}

	if len(migrations) == 0 {
		fmt.Println("Nothing to migrate")
		return nil
	}

	// get next batch number
	batchNumber, err := m.mhRepo.GetNextBatchNumber()
	if err != nil {
		return fmt.Errorf("failed to get next batch number: %v", err)
	}

	var history []*model.MigrationHistory
	for id, migration := range migrations {
		migration.Up(schema)
		history = append(history, &model.MigrationHistory{
			ID:    id,
			Batch: batchNumber,
		})
	}

	statements := schema.statement()
	err = m.mhRepo.MigrationMigrate(statements, history)
	if err != nil {
		return fmt.Errorf("failed migrate: %v", err)
	}

	fmt.Printf("Batch %d complete. %d migration executed.\n", batchNumber, len(migrations))
	return nil
}

func (m *Migrator) RollbackLastBatch() error {
	lastBatch, err := m.mhRepo.GetNextBatchNumber()
	if err != nil {
		return err
	}
	lastBatch--

	histories, err := m.mhRepo.GetMigrationHistoryByBatch(lastBatch, "id DESC")
	if err != nil {
		return fmt.Errorf("failed to get migration history: %v", err)
	}

	return m.rollbackMigrations(histories)
}

func (m *Migrator) RollbackAll() error {
	histories, err := m.mhRepo.GetMigrationHistory("batch DESC, id DESC")
	if err != nil {
		return fmt.Errorf("failed to get migration history: %v", err)
	}

	return m.rollbackMigrations(histories)
}

func (m *Migrator) rollbackMigrations(histories []model.MigrationHistory) error {
	dialect, err := m.getDialect()
	if err != nil {
		return err
	}

	schema := NewSchema(dialect)
	for _, history := range histories {
		for _, m := range m.migrations {
			if m.ID() == history.ID {
				m.Down(schema)
			}
		}
	}

	statements := schema.statement()
	err = m.mhRepo.MigrationRollback(statements, histories)
	if err != nil {
		return fmt.Errorf("failed rollback: %v", err)
	}

	return nil
}

func (m *Migrator) getDialect() (Dialect, error) {
	switch m.config.Database.Driver {
	default:
		return nil, fmt.Errorf("unsupported driver: %s", m.config.Database.Driver)
	case "pgsql", "postgres", "postgresql":
		return &dialect.PostgresDialect{}, nil
	case "mysql", "mariadb", "maria":
		return &dialect.MySqlDialect{}, nil
	case "sqlite", "sqlite3", "file":
		return &dialect.SQLiteDialect{}, nil
	}
}

func (m *Migrator) Status() error {
	histories, err := m.mhRepo.GetMigrationHistory("id ASC")
	if err != nil {
		return fmt.Errorf("failed to get migration history: %v", err)
	}

	executedMigrations := make(map[string]model.MigrationHistory)
	for _, h := range histories {
		executedMigrations[h.ID] = h
	}

	// Determine max column widths
	maxMigrationLen := len("Migration")
	maxBatchLen := len("Batch")
	maxExecutedLen := len("Executed At")

	for _, migration := range m.migrations {
		id := migration.ID()
		if len(id) > maxMigrationLen {
			maxMigrationLen = len(id)
		}
	}

	for _, history := range histories {
		batchStr := fmt.Sprintf("%d", history.Batch)
		executedStr := history.CreatedAt.Format("2006-01-02 15:04:05")

		if len(batchStr) > maxBatchLen {
			maxBatchLen = len(batchStr)
		}
		if len(executedStr) > maxExecutedLen {
			maxExecutedLen = len(executedStr)
		}
	}

	if len("Pending") > maxBatchLen {
		maxBatchLen = len("Pending")
	}

	if len("Not Executed") > maxExecutedLen {
		maxExecutedLen = len("Not Executed")
	}

	// Print table header
	fmt.Println("\nMigration Status:")
	fmt.Printf(
		"+-%s-+-%s-+-%s-+\n",
		strings.Repeat("-", maxMigrationLen),
		strings.Repeat("-", maxBatchLen),
		strings.Repeat("-", maxExecutedLen),
	)
	fmt.Printf(
		"| %-*s | %-*s | %-*s |\n",
		maxMigrationLen,
		"Migration",
		maxBatchLen,
		"Batch",
		maxExecutedLen,
		"Executed At",
	)
	fmt.Printf(
		"+-%s-+-%s-+-%s-+\n",
		strings.Repeat("-", maxMigrationLen),
		strings.Repeat("-", maxBatchLen),
		strings.Repeat("-", maxExecutedLen),
	)

	// Print migration statuses
	for _, migration := range m.migrations {
		id := migration.ID()
		if history, exists := executedMigrations[id]; exists {
			fmt.Printf(
				"| %-*s | %-*d | %-*s |\n",
				maxMigrationLen,
				id,
				maxBatchLen,
				history.Batch,
				maxExecutedLen,
				history.CreatedAt.Format("2006-01-02 15:04:05"),
			)
		} else {
			fmt.Printf("| %-*s | %-*s | %-*s |\n",
				maxMigrationLen, id, maxBatchLen, "Pending", maxExecutedLen, "Not Executed")
		}
	}

	// Print table footer
	fmt.Printf(
		"+-%s-+-%s-+-%s-+\n",
		strings.Repeat("-", maxMigrationLen),
		strings.Repeat("-", maxBatchLen),
		strings.Repeat("-", maxExecutedLen),
	)
	return nil
}

func (m *Migrator) CreateMigration(name string) error {
	// if the name is empty, generate a random name
	if name == "" {
		name = fmt.Sprintf("migration_%d", time.Now().UnixNano())
	}

	// add timestamp to migration name and make it lowercase
	filename := fmt.Sprintf("%s_%s", time.Now().Format("20060102150405"), strings.ToLower(name))

	// check if the directory exist
	if _, err := os.Stat(filepath.Join("database", "migration")); os.IsNotExist(err) {
		if err := os.Mkdir(filepath.Join("database", "migration"), 0755); err != nil {
			return err
		}
	}

	// check if the file exist
	if _, err := os.Stat(filepath.Join("database", "migration", fmt.Sprintf("%s.go", filename))); err == nil {
		return fmt.Errorf("migration %s already exist", filename)
	}

	// create a file in the migrations directory
	file, err := os.Create(filepath.Join("database", "migration", fmt.Sprintf("%s.go", filename)))
	if err != nil {
		return err
	}
	defer file.Close()

	// resolve name from snake_case to CamelCase
	name = cases.Title(language.English, cases.NoLower).String(strings.ReplaceAll(name, "_", " "))
	name = strings.ReplaceAll(name, " ", "")

	// write the migration template to the file
	err = m.tmpl.ExecuteTemplate(file, "new_migration", map[string]interface{}{
		"Name":     name,
		"Filename": filename,
	})
	if err != nil {
		return err
	}

	migrationEntry, err := os.Open(filepath.Join("database", "migration", "migration.go"))
	if err != nil {
		return err
	}

	// replace the /** PLACEHOLDER **/ with the struct name
	content, err := io.ReadAll(migrationEntry)
	if err != nil {
		return err
	}
	migrationEntry.Close()

	// replace the /** PLACEHOLDER **/ with the struct name
	content = bytes.ReplaceAll(
		content,
		[]byte("/** PLACEHOLDER **/"),
		[]byte("&"+name+"{},\n\t\t/** PLACEHOLDER **/"),
	)

	// create a new file in the migrations directory
	migrationEntry, err = os.Create(filepath.Join("database", "migration", "migration.go"))
	if err != nil {
		return err
	}
	defer migrationEntry.Close()

	// write the updated content to the file
	_, err = migrationEntry.Write(content)
	if err != nil {
		return err
	}

	return nil
}
