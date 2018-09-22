package service

import (
	"strings"
)

func stringToSlug(title string) string {
	// Filter and conver to lowercase
	slug := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return r + 'a' - 'A'
		case r >= 'a' && r <= 'z':
			return r
		case r >= '0' && r <= '9':
			return r
		}
		return ' '
	}

	// Convert whitespace to hyphen '-'
	str := strings.Map(slug, title)
	strarray := strings.Fields(str)
	str = strings.Join(strarray, "-")

	return str
}
