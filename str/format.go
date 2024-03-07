package str

import (
	"fmt"
	"regexp"
)

// Sprintf safely formats strings with placeholders
func Sprintf(format string, args ...interface{}) string {
	placeholderRegex := regexp.MustCompile(`%[+\-#0 ]?\d*(\.\d+)?[bcdefgosuxXp]`)
	if placeholderRegex.MatchString(format) {
		return fmt.Sprintf(format, args...)
	}
	return format
}
