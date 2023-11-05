package domain

const (
	QueryCreatedAt    = "created_at"
	QueryUpdatedAt    = "updated_at"
	SortDirectionDesc = "desc"
	SortDirectionAsc  = "asc"
)

type PaginationQuery struct {
	Page          int64  `json:"page" query:"page" validate:"required,gt=0"`
	Limit         int64  `json:"limit" query:"limit" validate:"required,gte=-1,ne=0,lte=100"`
	Sort          string `json:"-" query:"-"`
	SortDirection string `json:"-" query:"-"`
}

type PaginationResponse[T any] struct {
	PaginationQuery
	TotalRows  int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
	Rows       []T   `json:"rows"`
}
