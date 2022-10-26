package data

import "github.com/canyolal/greenlight/internal/validator"

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

// Validates page, page_size and sort.
func ValidateFilters(v *validator.Validator, filter *Filters) {
	v.Check(filter.Page >= 1 && filter.Page <= 10000000, "page", "must be in between 1-10,000,000")
	v.Check(filter.PageSize >= 1 && filter.PageSize <= 100, "page_size", "must be in between 1-100")

	// Only allowed strings can be sorted.
	v.Check(validator.In(filter.Sort, filter.SortSafeList...), "sort", "invalid sort value")
}
