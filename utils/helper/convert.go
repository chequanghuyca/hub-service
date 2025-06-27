package helper

import (
	"fmt"
	"regexp"
	"strings"
)

// normalizeString converts a string to lowercase, removes punctuation and extra spaces.
func NormalizeString(s string) string {
	s = strings.ToLower(s)
	re := regexp.MustCompile(`[^\w\s]`)
	s = re.ReplaceAllString(s, "")
	return strings.Join(strings.Fields(s), " ")
}

// toFloat64 safely converts an interface{} to float64, handling multiple numeric types.
func ToFloat64(v interface{}) (float64, error) {
	switch i := v.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int:
		return float64(i), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", v)
	}
}
