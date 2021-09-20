package calculate

import (
	"github.com/gobwas/glob"
	"strings"
	"unicode"
	"unicode/utf8"
)

func compileGlobs(patterns []string) ([]glob.Glob, error) {
	patterns = expandPatternsIfNeeded(patterns)
	globs := make([]glob.Glob, len(patterns))
	for i, pattern := range patterns {
		compiled, err := glob.Compile(pattern)
		if err != nil {
			return nil, err
		}
		globs[i] = compiled
	}
	return globs, nil
}

func expandPatternsIfNeeded(patterns []string) []string {
	for _, pattern := range patterns {
		if strings.HasPrefix(pattern, "*/") {
			patterns = append(patterns, strings.Replace(pattern, "*/", "", 1))
		}
		if strings.HasPrefix(pattern, "**/") {
			patterns = append(patterns, strings.Replace(pattern, "**/", "", 1))
		}
	}
	return patterns
}

func matches(path string, patterns []glob.Glob) bool {
	for _, pattern := range patterns {
		if pattern.Match(path) {
			return true
		}
	}
	return false
}

func trimSpaceLeft(s string) string {
	// adapted from strings.TrimSpace
	start := 0
	for ; start < len(s); start++ {
		c := s[start]
		if c >= utf8.RuneSelf {
			// If we run into a non-ASCII byte, fall back to the
			// slower unicode-aware method on the remaining bytes
			return strings.TrimFunc(s[start:], unicode.IsSpace)
		}
		if asciiSpace[c] == 0 {
			break
		}
	}
	if start == 0 {
		return s
	}
	return s[start:]
}

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}
