package model

// Response API响应结构
type Response struct {
	Code    int         `json:"code"`           // 状态码，0表示成功
	Message string      `json:"message"`        // 响应消息
	Data    interface{} `json:"data,omitempty"` // 响应数据
}

func Success(data interface{}) Response {
	return Response{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

func Fail(msg string) Response {
	return Response{
		Code:    1,
		Message: msg,
	}
}

// PageResponse 分页响应结构
type PageResponse struct {
	Data     interface{} `json:"data"`      // 数据列表
	Total    int         `json:"total"`     // 总数量
	Page     int         `json:"page"`      // 当前页码
	PageSize int         `json:"page_size"` // 每页数量
}
