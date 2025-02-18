package migration

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"go.uber.org/fx"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/mrrizkin/pohara/modules/database/config"
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
	if err := deps.Db.AutoMigrate(model.MigrationHistory{}); err != nil {
		panic(fmt.Sprintf("failed to migrate the __migration_history__: %v", err))
	}

	tmpl, err := template.ParseFS(templateFS, "template/*.go.tmpl")
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
	histories, err := m.mhRepo.GetMigrationHistory("id ASC")
	if err != nil {
		return fmt.Errorf("failed to get migration history: %v", err)
	}

	executedMigrations := make(map[string]model.MigrationHistory)
	for _, h := range histories {
		executedMigrations[h.ID] = h
	}

	pendingMigrations := m.getPendingMigrations(executedMigrations)
	if len(pendingMigrations) == 0 {
		fmt.Println("Nothing to migrate")
		return nil
	}

	batchNumber, err := m.mhRepo.GetNextBatchNumber()
	if err != nil {
		return fmt.Errorf("failed to get next batch number: %v", err)
	}

	history := m.applyMigrations(schema, pendingMigrations, batchNumber)
	statements := schema.statement()

	if err := m.mhRepo.MigrationMigrate(statements, history); err != nil {
		return fmt.Errorf("failed migrate: %v", err)
	}

	fmt.Printf("Batch %d complete. %d migration executed.\n", batchNumber, len(pendingMigrations))
	return nil
}

func (m *Migrator) getPendingMigrations(
	executedMigrations map[string]model.MigrationHistory,
) map[string]Migration {
	pendingMigrations := make(map[string]Migration)
	for _, migration := range m.migrations {
		id := migration.ID()
		if _, exists := executedMigrations[id]; !exists {
			pendingMigrations[id] = migration
		}
	}
	return pendingMigrations
}

func (m *Migrator) applyMigrations(
	schema *Schema,
	migrations map[string]Migration,
	batchNumber int,
) []*model.MigrationHistory {
	var history []*model.MigrationHistory
	for id, migration := range migrations {
		migration.Up(schema)
		history = append(history, &model.MigrationHistory{
			ID:    id,
			Batch: batchNumber,
		})
	}
	return history
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
		for _, migration := range m.migrations {
			if migration.ID() == history.ID {
				migration.Down(schema)
			}
		}
	}

	statements := schema.statement()
	if err := m.mhRepo.MigrationRollback(statements, histories); err != nil {
		return fmt.Errorf("failed rollback: %v", err)
	}

	return nil
}

func (m *Migrator) getDialect() (Dialect, error) {
	switch m.config.Driver {
	case "pgsql", "postgres", "postgresql":
		return &dialect.PostgresDialect{}, nil
	case "mysql", "mariadb", "maria":
		return &dialect.MySqlDialect{}, nil
	case "sqlite", "sqlite3", "file":
		return &dialect.SQLiteDialect{}, nil
	default:
		return nil, fmt.Errorf("unsupported driver: %s", m.config.Driver)
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

	maxMigrationLen, maxBatchLen, maxExecutedLen := m.calculateColumnWidths(histories)
	m.printMigrationStatusTable(executedMigrations, maxMigrationLen, maxBatchLen, maxExecutedLen)

	return nil
}

func (m *Migrator) calculateColumnWidths(histories []model.MigrationHistory) (int, int, int) {
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

	return maxMigrationLen, maxBatchLen, maxExecutedLen
}

func (m *Migrator) printMigrationStatusTable(
	executedMigrations map[string]model.MigrationHistory,
	maxMigrationLen, maxBatchLen, maxExecutedLen int,
) {
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

	fmt.Printf(
		"+-%s-+-%s-+-%s-+\n",
		strings.Repeat("-", maxMigrationLen),
		strings.Repeat("-", maxBatchLen),
		strings.Repeat("-", maxExecutedLen),
	)
}

func (m *Migrator) CreateMigration(name string) error {
	if name == "" {
		name = fmt.Sprintf("migration_%d", time.Now().UnixNano())
	}

	filename := fmt.Sprintf("%s_%s", time.Now().Format("20060102150405"), strings.ToLower(name))
	migrationDir := filepath.Join("database", "migration")
	migrationFile := filepath.Join(migrationDir, fmt.Sprintf("%s.go", filename))

	if err := os.MkdirAll(migrationDir, 0755); err != nil {
		return err
	}

	if _, err := os.Stat(migrationFile); err == nil {
		return fmt.Errorf("migration %s already exist", filename)
	}

	file, err := os.Create(migrationFile)
	if err != nil {
		return err
	}
	defer file.Close()

	name = cases.Title(language.English, cases.NoLower).String(strings.ReplaceAll(name, "_", " "))
	name = strings.ReplaceAll(name, " ", "")

	if err := m.tmpl.ExecuteTemplate(file, "new_migration", map[string]interface{}{
		"Name":     name,
		"Filename": filename,
	}); err != nil {
		return err
	}

	migrationEntryPath := filepath.Join(migrationDir, "migration.go")
	content, err := os.ReadFile(migrationEntryPath)
	if err != nil {
		return err
	}

	content = bytes.ReplaceAll(
		content,
		[]byte("/** PLACEHOLDER **/"),
		[]byte("&"+name+"{},\n\t\t/** PLACEHOLDER **/"),
	)

	if err := os.WriteFile(migrationEntryPath, content, 0644); err != nil {
		return err
	}

	return nil
}
