package sql

type Pagination struct {
	Limit  Int64Nullable
	Offset Int64Nullable
	Sort   StringNullable
}

type PaginationResult struct {
	Data      interface{}
	Total     int64
	TotalPage Int64Nullable
	Page      Int64Nullable
	Limit     Int64Nullable
}
