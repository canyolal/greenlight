package data

import (
	"math"
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

// limit determines page size
func (f Filters) limit() int {
	return f.PageSize
}

// offset returns starting index of pagination
func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

// Metadata struct to hold pagination metadata.
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

// The calculateMetadata() function calculates the appropriate pagination metadata
// values given the total number of records, current page, and page size values. Note
// that the last page value is calculated using the math.Ceil() function, which rounds
// up a float to the nearest integer. So, for example, if there were 12 records in total
// and a page size of 5, the last page value would be math.Ceil(12/5) = 3.
func calculateMetadata(totalRecords, page, pageSize int) Metadata {

	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}

}
