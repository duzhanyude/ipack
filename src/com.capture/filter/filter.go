package filter

import "strings"

type Filter interface {
	Handler(content string) bool
}

type ContentFilter struct {
	FilterContent string
}

func (c *ContentFilter) Handler(content string) bool {

	if content != "" && c.FilterContent != "" {
		if !strings.Contains(content, c.FilterContent) {
			return false
		}
	}
	return true
}
