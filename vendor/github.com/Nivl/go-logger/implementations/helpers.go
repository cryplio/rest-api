// Package implementations contains helpers to help implementing a logger
package implementations

import "fmt"
import "strings"

// FormatError is an helper to format the output of an error
func FormatError(staticData []string, msg string) string {
	static := strings.Join(staticData, ", ")
	if static != "" {
		static = ", " + static
	}
	return fmt.Sprintf(`level: "ERROR"%s, %s"`, static, msg)
}
