package data

import (
	"strings"

	"github.com/canyolal/greenlight/internal/validator"
)

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

// Check that the client-provided Sort field matches one of the entries in our safelist
// and if it does, extract the column name from the Sort field by stripping the leading
// hyphen character (if one exists).
func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafeList {
		if safeValue == f.Sort {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + f.Sort)
}

// sortDirection returns the direction of sorting for movie filtering depending on '+' or '-'
func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}
