package repository

import (
	"fmt"

	"github.com/mrrizkin/pohara/modules/database/db"
	"github.com/mrrizkin/pohara/modules/database/model"
	"github.com/mrrizkin/pohara/modules/database/utils"
	"go.uber.org/fx"
)

type MigrationHistoryRepository struct {
	db *db.Database
}

type MigrationHistoryRepositoryDeps struct {
	fx.In

	Database *db.Database
}

func NewMigrationHistoryRepository(
	deps MigrationHistoryRepositoryDeps,
) *MigrationHistoryRepository {
	return &MigrationHistoryRepository{
		db: deps.Database,
	}
}

func (m *MigrationHistoryRepository) GetNextBatchNumber() (int, error) {
	var lastBatch struct {
		MaxBatch int
	}

	err := m.db.Model(model.MigrationHistory{}).
		Select("COALESCE(MAX(batch), 0) as max_batch").
		Scan(&lastBatch).
		Error
	if err != nil {
		return 0, err
	}

	return lastBatch.MaxBatch + 1, nil
}

func (m *MigrationHistoryRepository) MigrationMigrate(
	statements []string,
	histories []*model.MigrationHistory,
) error {
	tx := m.db.Begin()
	for _, statement := range statements {
		if err := tx.Exec(statement).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration: %v", err)
		}
	}

	if err := tx.Create(histories).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record migration: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit migration: %v", err)
	}

	return nil
}

func (m *MigrationHistoryRepository) MigrationRollback(
	statements []string,
	histories []model.MigrationHistory,
) error {
	tx := m.db.Begin()
	for _, statement := range statements {
		if err := tx.Exec(statement).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to rollback migration: %v", err)
		}
	}

	if err := tx.Create(histories).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record migration: %v", err)
	}

	listHistoryId := utils.Pluck(histories, func(h model.MigrationHistory) string { return h.ID })
	if err := tx.Delete(&model.MigrationHistory{}, "id IN (?)", listHistoryId).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove migration record: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit migration: %v", err)
	}

	listBatch := utils.KeyBy(histories, func(h model.MigrationHistory) int { return h.Batch })
	for batch := range listBatch {
		fmt.Printf("Rolled back migration: (Batch %d)\n", batch)
	}

	return nil
}

func (m *MigrationHistoryRepository) GetMigrationHistory(
	order string,
) ([]model.MigrationHistory, error) {
	var histories []model.MigrationHistory
	result := m.db.Order(order).Find(&histories)
	if result.Error != nil {
		return nil, result.Error
	}

	return histories, nil
}

func (m *MigrationHistoryRepository) GetMigrationHistoryByBatch(
	batch int,
	order string,
) ([]model.MigrationHistory, error) {
	var histories []model.MigrationHistory
	result := m.db.Where("batch = ?", batch).
		Order(order).
		Find(&histories)
	if result.Error != nil {
		return nil, result.Error
	}

	return histories, nil
}
