package common

// Response is a generic API response wrapper for Swagger documentation
type Response struct {
	Status string      `json:"status" example:"success"`
	Data   interface{} `json:"data"`
}

type successResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Paging interface{} `json:"meta,omitempty"`
	Filter interface{} `json:"filter,omitempty"`
}

func NewSuccessResponse(data, paging, filter interface{}) *successResponse {
	return &successResponse{
		Status: "success",
		Data:   data,
		Paging: paging,
		Filter: filter,
	}
}

func SimpleSuccessResponse(data interface{}) *successResponse {
	return NewSuccessResponse(data, nil, nil)
}
