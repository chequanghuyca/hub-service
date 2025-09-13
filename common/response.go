package common

// Response is a generic API response wrapper for Swagger documentation
type Response struct {
	Status string      `json:"status" example:"success"`
	Data   interface{} `json:"data"`
}

// SuccessResponse is the standard API response structure
type SuccessResponse struct {
	Status string      `json:"status" example:"success"`
	Data   interface{} `json:"data"`
	Paging interface{} `json:"meta,omitempty"`
	Filter interface{} `json:"filter,omitempty"`
}

func NewSuccessResponse(data, paging, filter interface{}) *SuccessResponse {
	return &SuccessResponse{
		Status: "success",
		Data:   data,
		Paging: paging,
		Filter: filter,
	}
}

func SimpleSuccessResponse(data interface{}) *SuccessResponse {
	return NewSuccessResponse(data, nil, nil)
}
