package paginator

import "errors"
import "fmt"

const (
	// ErrMsgNumberBelow1 represents the error message corresponding to
	// a number being "< 1"
	ErrMsgNumberBelow1 = "cannot be <= 0"

	// ErrInvalidValue represents the error message corresponding to
	// a number being invalid
	ErrInvalidValue = "invalid value"
)

// HandlerParams represents the params needed to handle a pagination
type HandlerParams struct {
	// Page represents the number of the current page
	Page int `from:"query" json:"page" default:"1"`

	// PerPage represents the maximum number of element we want per page
	PerPage int `from:"query" json:"per_page" default:"100"`
}

// IsValid checks if a paginator is Valid
func (params *HandlerParams) IsValid() (isValid bool, fieldName string, err error) {
	p := New(params.Page, params.PerPage)
	if p.IsValid() {
		return true, "", nil
	}

	if p.currentPage <= 0 {
		return false, "page", errors.New(ErrMsgNumberBelow1)
	}

	if p.perPage > p.MaxPerPage {
		return false, "per_page", fmt.Errorf("cannot be > %d", p.MaxPerPage)
	}

	if p.perPage <= 0 {
		return false, "per_page", errors.New(ErrMsgNumberBelow1)
	}

	return false, "page/per_page", errors.New(ErrInvalidValue)
}

// Paginator returns a Paginator from an HandlerParams
func (params HandlerParams) Paginator() *Paginator {
	return New(params.Page, params.PerPage)
}

// Paginator represents a pagination
type Paginator struct {
	currentPage int
	perPage     int
	MaxPerPage  int
}

// New creates a new Paginator
func New(page int, perPage int) *Paginator {
	return &Paginator{
		currentPage: page,
		perPage:     perPage,
		MaxPerPage:  100,
	}
}

// IsValid checks if the paginator is valid
func (p *Paginator) IsValid() bool {
	// Also Update HandlerParams.IsValid() to return a specific error message
	// to the users
	return (p.currentPage > 0) && (p.perPage <= p.MaxPerPage) && (p.perPage > 0)
}

// Offset returns a valid SQL offset value
func (p *Paginator) Offset() int {
	return (p.currentPage - 1) * p.perPage
}

// Limit returns a valid SQL limit value
func (p *Paginator) Limit() int {
	return p.perPage
}

// CurrentPage returns the current page index
func (p *Paginator) CurrentPage() int {
	return p.currentPage
}

// PerPage returns a the number of data a page should contain
func (p *Paginator) PerPage() int {
	return p.perPage
}
