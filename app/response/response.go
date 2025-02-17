package response

import "github.com/mrrizkin/pohara/modules/common/sql"

type Pagination struct {
	Total     int64             `json:"total"`
	TotalPage sql.Int64Nullable `json:"total_page,omitempty"`
	Page      sql.Int64Nullable `json:"page,omitempty"`
	Limit     sql.Int64Nullable `json:"limit,omitempty"`
}

type GeneralResponse[T any] struct {
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Detail     string      `json:"detail,omitempty"`
	Data       T           `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

func Success[T any](msg string, data T) GeneralResponse[T] {
	return GeneralResponse[T]{
		Status:  "success",
		Message: msg,
		Data:    data,
	}
}

func SuccessMsg(msg string) GeneralResponse[any] {
	return GeneralResponse[any]{
		Status:  "success",
		Message: msg,
	}
}

func SuccessPaginate[T any](msg string, data T, pagination *Pagination) GeneralResponse[T] {
	return GeneralResponse[T]{
		Status:     "success",
		Message:    msg,
		Data:       data,
		Pagination: pagination,
	}
}

func Error[T any](msg string, detail string, data T) GeneralResponse[T] {
	return GeneralResponse[T]{
		Status:  "error",
		Message: msg,
		Detail:  detail,
		Data:    data,
	}
}

func ErrorMsg(msg string, detail string) GeneralResponse[any] {
	return GeneralResponse[any]{
		Status:  "error",
		Message: msg,
		Detail:  detail,
	}
}

func ErrorPagination[T any](
	msg string,
	detail string,
	data T,
	pagination *Pagination,
) GeneralResponse[T] {
	return GeneralResponse[T]{
		Status:     "error",
		Message:    msg,
		Detail:     detail,
		Data:       data,
		Pagination: pagination,
	}
}
