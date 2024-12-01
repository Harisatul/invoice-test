package model

type APIResponse struct {
	Status          int         `json:"status"`
	Message         string      `json:"message"`
	Data            interface{} `json:"data,omitempty"`
	Errors          interface{} `json:"error,omitempty"`
	PaginationIndex interface{} `json:"pagination_index,omitempty"`
}

type PaginationIndex struct {
	Page        int  `json:"page"`
	PageSize    int  `json:"page_size"`
	TotalCount  int  `json:"total_count"`
	TotalPages  int  `json:"total_pages"`
	HasPrevious bool `json:"has_previous"`
	HasNext     bool `json:"has_next"`
}
