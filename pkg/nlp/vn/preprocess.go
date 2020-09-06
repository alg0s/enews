package vn

import "regexp"

func removeMissingContent(c string) string {
	return ``
}

func removeNewline(c string) string {
	return ``
}

// cleanEntity removes unwanted characters in a token,
// and replaces whitespace with underscore
func cleanEntity(e string) string {
	var patterns = map[string]string{
		`\s|-`:      `_`,
		`~!@#$%^&*`: ``,
		`^_+|_+$`:   ``,
	}
	for old, new := range patterns {
		m := regexp.MustCompile(old)
		e = m.ReplaceAllString(e, new)
	}
	return e
}
