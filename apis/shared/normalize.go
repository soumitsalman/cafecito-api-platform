package shared

import (
	"github.com/stoewer/go-strcase"
)

func NormalizeTags(items []string) []string {
	for i, item := range items {
		items[i] = strcase.SnakeCase(item)
	}
	return items
}
