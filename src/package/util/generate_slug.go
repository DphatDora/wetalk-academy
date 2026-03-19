package util

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func removeVietnameseTones(s string) string {
	// normalize unicode
	t := norm.NFD.String(s)

	// remove diacritics
	result := make([]rune, 0, len(t))
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		result = append(result, r)
	}

	return string(result)
}

func GenerateSlug(topic string) string {
	slug := strings.ToLower(topic)
	slug = removeVietnameseTones(slug)

	reg := regexp.MustCompile(`[^a-z0-9\s-]`)
	slug = reg.ReplaceAllString(slug, "")

	slug = strings.TrimSpace(slug)
	slug = regexp.MustCompile(`\s+`).ReplaceAllString(slug, "-")

	suffix := time.Now().Unix()

	return fmt.Sprintf("%s-%d", slug, suffix)
}
