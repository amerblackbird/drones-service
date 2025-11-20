package domain

type PaginationOption[T any] struct {
	Limit  int `json:"limit" validate:"min=1,max=100"`
	Offset int `json:"offset" validate:"min=0"`
	Filter *T  `json:"filter,omitempty"`
}

type FindOneOption[T any] struct {
	Filter *T `json:"filter,omitempty"`
}

type Pagination[T any] struct {
	Data       []*T `json:"data"`
	Total      int  `json:"total"`
	Page       int  `json:"page"`
	PageSize   int  `json:"page_size"`
	TotalPages int  `json:"total_pages"`
}

// Helper function to create pagination options
func NewPaginationOption[T any](limit, offset int, filter *T) *PaginationOption[T] {
	return &PaginationOption[T]{
		Limit:  limit,
		Offset: offset,
		Filter: filter,
	}
}

// Helper function to create find one options
func NewFindOneOption[T any](filter *T) *FindOneOption[T] {
	return &FindOneOption[T]{
		Filter: filter,
	}
}
