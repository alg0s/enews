package vn

import "regexp"

func removeMissingContent(c string) string {
	return c
}

func removeNewline(c string) string {
	return c
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

// chunkIDs groups a list of items into smaller chunks
func chunkIDs(size int, items []int32) [][]int32 {
	total := (len(items) + size - 1) / size
	batches := make([][]int32, 0, total)
	for size < len(items) {
		items, batches = items[size:], append(batches, items[0:size])
	}
	batches = append(batches, items)
	return batches
}
