package repository

import (
	"math"

	"github.com/mrrizkin/pohara/modules/common/sql"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Module = fx.Module("repository",
	fx.Provide(
		NewUserRepository,
		NewSettingRepository,
	),
)

type PaginationResult[T any] struct {
	Data      []T
	Total     int64
	TotalPage sql.Int64Nullable
	Page      sql.Int64Nullable
	Limit     sql.Int64Nullable
}

type QueryPaginateParams struct {
	Page  sql.Int64Nullable
	Limit sql.Int64Nullable
}

func QueryPaginate[T any](db *gorm.DB, m []T, params QueryPaginateParams) (result *PaginationResult[T], err error) {
	var total int64
	err = db.Session(&gorm.Session{NewDB: true}).Count(&total).Error
	if err != nil {
		return
	}

	if params.Limit.Valid {
		db.Limit(int(params.Limit.Int64))
	}

	offset := sql.Int64Null()
	if params.Page.Valid && params.Page.Int64 != 0 {
		offset.Valid = true
		offset.Int64 = (params.Page.Int64 - 1)
		if params.Limit.Int64 != 0 {
			offset.Int64 = params.Page.Int64 * params.Limit.Int64
		}

		db.Offset(int(offset.Int64))
	}

	err = db.Find(&m).Error
	if err != nil {
		return
	}

	totalPage := sql.Int64Null()
	if offset.Valid && params.Limit.Valid && params.Limit.Int64 != 0 {
		totalPage.Valid = true
		totalPage.Int64 = int64(math.Ceil(float64(total)/float64(params.Limit.Int64))) + 1
	}

	result = &PaginationResult[T]{
		Data:      m,
		Total:     total,
		TotalPage: totalPage,
		Page:      params.Page,
		Limit:     params.Limit,
	}

	return
}
