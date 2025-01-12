package response

import "github.com/mrrizkin/pohara/modules/common/sql"

type Pagination struct {
	Total     int64             `json:"total"`
	TotalPage sql.Int64Nullable `json:"total_page,omitempty"`
	Page      sql.Int64Nullable `json:"page,omitempty"`
	Limit     sql.Int64Nullable `json:"limit,omitempty"`
}

type GeneralResponse struct {
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Detail     string      `json:"detail,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

func Success(msg string, data interface{}) GeneralResponse {
	return GeneralResponse{
		Status:  "success",
		Message: msg,
		Data:    data,
	}
}

func SuccessMsg(msg string) GeneralResponse {
	return GeneralResponse{
		Status:  "success",
		Message: msg,
	}
}

func SuccessPaginate(msg string, data interface{}, pagination *Pagination) GeneralResponse {
	return GeneralResponse{
		Status:     "success",
		Message:    msg,
		Data:       data,
		Pagination: pagination,
	}
}

func Error(msg string, detail string, data interface{}) GeneralResponse {
	return GeneralResponse{
		Status:  "error",
		Message: msg,
		Detail:  detail,
		Data:    data,
	}
}

func ErrorMsg(msg string, detail string) GeneralResponse {
	return GeneralResponse{
		Status:  "error",
		Message: msg,
		Detail:  detail,
	}
}

func ErrorPagination(
	msg string,
	detail string,
	data interface{},
	pagination *Pagination,
) GeneralResponse {
	return GeneralResponse{
		Status:     "error",
		Message:    msg,
		Detail:     detail,
		Data:       data,
		Pagination: pagination,
	}
}
