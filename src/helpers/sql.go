package helpers

import (
	"fmt"

	"github.com/lib/pq"
)

// SQLConflictError returns an SQL Conlift error
func SQLConflictError(fieldname, value string) error {
	conflictError := new(pq.Error)
	conflictError.Code = "23505"
	conflictError.Message = "error: duplicate field"
	conflictError.Detail = fmt.Sprintf("Key (%s)=(%s) already exists.", fieldname, value)
	return conflictError
}
