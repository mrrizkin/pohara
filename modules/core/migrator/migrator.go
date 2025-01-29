package migrator

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"go.uber.org/fx"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/mrrizkin/pohara/app/config"
	"github.com/mrrizkin/pohara/modules/cli"
	"github.com/mrrizkin/pohara/modules/core/migration"
	"github.com/mrrizkin/pohara/modules/database"
)

type Migration interface {
	Up(schema *migration.Schema)
	Down(schema *migration.Schema)
	ID() string
}

type MigrationHistory struct {
	ID        string    `db:"id"`
	Batch     int       `db:"batch"`
	CreatedAt time.Time `db:"created_at"`
}

func (MigrationHistory) TableName() string {
	return "__migration_history__"
}

type MigrationLoaderDependencies struct {
	fx.In

	Migrator   *Migrator
	Migrations []Migration `group:"migration"`
}

type Migrator struct {
	db         *database.Database
	config     *config.Config
	migrations []Migration
}

type MigratorDependencies struct {
	fx.In

	Db     *database.Database
	Config *config.Config
}

var Module = fx.Module("migrator",
	fx.Provide(
		NewMigrator,
		cli.AsCommand(NewMigratorCommand),
	),
	fx.Decorate(LoadMigration),
)

func NewMigrator(deps MigratorDependencies) *Migrator {
	return &Migrator{
		db:     deps.Db,
		config: deps.Config,
	}
}

func LoadMigration(deps MigrationLoaderDependencies) *Migrator {
	sortedMigrations := make([]Migration, len(deps.Migrations))
	copy(sortedMigrations, deps.Migrations)

	sort.Slice(sortedMigrations, func(i, j int) bool {
		return sortedMigrations[i].ID() < sortedMigrations[j].ID()
	})

	for _, migration := range sortedMigrations {
		deps.Migrator.AddMigration(migration)
	}

	return deps.Migrator
}

func ProvideMigration(migrations ...Migration) fx.Option {
	var options []fx.Option

	for _, migration := range migrations {
		m := migration
		options = append(options, fx.Provide(fx.Annotate(
			func() Migration { return m },
			fx.ResultTags(`group:"migration"`),
		)))
	}

	return fx.Options(options...)
}

func (m *Migrator) AddMigration(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

func (m *Migrator) getNextBatchNumber() (int, error) {
	var lastBatch struct {
		MaxBatch int
	}

	err := m.db.Get(&lastBatch, "COALESCE(MAX(batch), 0) as max_batch")
	if err != nil {
		return 0, err
	}

	return lastBatch.MaxBatch + 1, nil
}

func (m *Migrator) Migrate() error {
	dialect, err := m.getDialect()
	if err != nil {
		return err
	}

	schema := migration.NewSchema(dialect)
	createMigrationTable(schema)

	migrations := make(map[string]Migration)
	var histories []MigrationHistory

	if err := m.db.Select(&histories, "SELECT * FROM __migration_history__ ORDER BY id ASC"); err != nil {
		return fmt.Errorf("failed to get migration history: %v", err)
	}

	executedMigrations := make(map[string]MigrationHistory)
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
	batchNumber, err := m.getNextBatchNumber()
	if err != nil {
		return fmt.Errorf("failed to get next batch number: %v", err)
	}

	var history []*MigrationHistory
	for id, migration := range migrations {
		migration.Up(schema)
		history = append(history, &MigrationHistory{
			ID:    id,
			Batch: batchNumber,
		})
	}

	statements := schema.Statement()
	tx, err := m.db.Beginx()
	if err != nil {
		return err
	}
	for _, statement := range statements {
		if _, err := tx.Exec(statement); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration: %v", err)
		}
	}

	_, err = m.db.Builder().
		Table("__migration_history__").
		Insert().
		Values(history).
		WithTx(tx).
		Exec()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record migration: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration: %v", err)
	}

	fmt.Printf("Batch %d complete. %d migration executed.\n", batchNumber, len(migrations))
	return nil
}

func (m *Migrator) RollbackLastBatch() error {
	lastBatch, err := m.getNextBatchNumber()
	if err != nil {
		return err
	}
	lastBatch--

	var histories []MigrationHistory
	err = m.db.Builder().
		Table("__migration_history__").
		Where("batch = ?", lastBatch).
		OrderBy("id", "DESC").
		All(&histories)
	if err != nil {
		return fmt.Errorf("failed to get migration history: %v", err)
	}

	return m.rollbackMigrations(histories)
}

func (m *Migrator) RollbackAll() error {
	var histories []MigrationHistory
	err := m.db.Select(
		&histories,
		"SELECT * FROM __migration_history__ ORDER BY batch DESC, id DESC",
	)
	if err != nil {
		return fmt.Errorf("failed to get migration history: %v", err)
	}

	return m.rollbackMigrations(histories)
}

func (m *Migrator) rollbackMigrations(histories []MigrationHistory) error {
	dialect, err := m.getDialect()
	if err != nil {
		return err
	}

	schema := migration.NewSchema(dialect)
	for _, history := range histories {
		for _, m := range m.migrations {
			if m.ID() == history.ID {
				m.Down(schema)
			}
		}
	}

	statements := schema.Statement()
	tx, err := m.db.Beginx()
	if err != nil {
		return err
	}

	for _, statement := range statements {
		if _, err := tx.Exec(statement); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to rollback migration: %v", err)
		}
	}

	for _, history := range histories {
		if _, err := tx.Exec("DELETE FROM __migration_history__ WHERE id = ?", history.ID); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to remove migration record: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit rollback: %v", err)
	}

	for _, history := range histories {
		fmt.Printf("Rolled back migration: (Batch %d)\n", history.Batch)
	}

	return nil
}

func (m *Migrator) getDialect() (migration.Dialect, error) {
	switch m.config.Database.Driver {
	default:
		return nil, fmt.Errorf("unsupported driver: %s", m.config.Database.Driver)
	case "pgsql", "postgres", "postgresql":
		return &migration.PostgresDialect{}, nil
	case "mysql", "mariadb", "maria":
		return &migration.MySqlDialect{}, nil
	case "sqlite", "sqlite3", "file":
		return &migration.SQLiteDialect{}, nil
	}
}

func (m *Migrator) Status() error {
	var histories []MigrationHistory
	err := m.db.Builder().
		Table("__migration_history__").
		OrderBy("id", "ASC").
		All(&histories)
	if err != nil {
		return fmt.Errorf("failed to get migration history: %v", err)
	}

	fmt.Println("\nMigration Status:")
	fmt.Println("+-----------------+--------+---------------------+")
	fmt.Println("| Migration       | Batch  | Executed At        |")
	fmt.Println("+-----------------+--------+---------------------+")

	executedMigrations := make(map[string]MigrationHistory)
	for _, h := range histories {
		executedMigrations[h.ID] = h
	}

	for _, migration := range m.migrations {
		id := migration.ID()
		if history, exists := executedMigrations[id]; exists {
			fmt.Printf("| %-15s | %-6d | %-19s |\n",
				id, history.Batch, history.CreatedAt.Format("2006-01-02 15:04:05"))
		} else {
			fmt.Printf("| %-15s | Pending | Not Executed        |\n", id)
		}
	}

	fmt.Println("+-----------------+--------+---------------------+")
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
	_, err = fmt.Fprintf(file, migrationTemplate, name, name, filename, name, name)
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

func createMigrationTable(schema *migration.Schema) {
	schema.CreateIfNotExist("__migration_history__", func(table *migration.Blueprint) {
		table.Text("id").Primary()
		table.Integer("batch")
		table.Timestamp("created_at")
	})
}

// migrationTemplate is the template for a migration file
const migrationTemplate = `package migration

import (
	"github.com/mrrizkin/pohara/modules/core/migration"
)

type %s struct{}

func (m *%s) ID() string {
	return "%s"
}

func (m *%s) Up(schema *migration.Schema) {
	// your migration schema here
}

func (m *%s) Down(schema *migration.Schema) {
	// your rollback migration here
}
`
