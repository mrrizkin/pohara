package sql

type PaginationResult struct {
	Data      interface{}
	Total     int64
	TotalPage Int64Nullable
	Page      Int64Nullable
	Limit     Int64Nullable
}
