package shared

import (
	"strings"
	"time"

	"github.com/stoewer/go-strcase"
)

func NormalizeTags(items []string) []string {
	for i, item := range items {
		items[i] = strcase.SnakeCase(item)
	}
	return items
}

func NormalizeText(item string) string {
	return strings.ToLower(strings.TrimSpace(item))
}

func NormalizeTexts(items []string) []string {
	for i, value := range items {
		items[i] = NormalizeText(value)
	}
	return items
}

// // NormalizeEndOfDay returns 00:00:00 UTC of the calendar day after day.
// // Use with a strict-less-than comparison so a YYYY-MM-DD `to` includes that whole UTC date.
func NormalizeEndOfDay(date time.Time) time.Time {
	if date.IsZero() {
		return date
	}
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location()).AddDate(0, 0, 1)
}
