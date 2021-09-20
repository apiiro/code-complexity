package calculate

import (
	"code-complexity/options"
	"fmt"
	"github.com/gobwas/glob"
	"golang.org/x/net/html/charset"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type context struct {
	CodeSummary
	includePatterns  []glob.Glob
	excludePatterns  []glob.Glob
	verboseLogging   bool
	maxFileSizeBytes int64
}

func Complexity(opts *options.Options) (*CodeSummary, error) {

	includePatterns, err := compileGlobs(opts.IncludePatterns)
	if err != nil {
		return nil, fmt.Errorf("failed to compile include patterns: %v", err)
	}
	excludePatterns, err := compileGlobs(opts.ExcludePatterns)
	if err != nil {
		return nil, fmt.Errorf("failed to compile exclude patterns: %v", err)
	}
	ctx := &context{
		CodeSummary: CodeSummary{
			CountersByLanguage: make(map[Language]*CodeCounters),
			AveragesByLanguage: make(map[Language]*CodeCounters),
		},
		includePatterns:  includePatterns,
		excludePatterns:  excludePatterns,
		verboseLogging:   opts.VerboseLogging,
		maxFileSizeBytes: opts.MaxFileSizeBytes,
	}

	err = filepath.Walk(
		opts.CodePath,
		func(path string, info fs.FileInfo, _ error) error {
			return ctx.visitPath(path, info)
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to walk files under '%v': %v", opts.CodePath, err)
	}

	for language, counters := range ctx.CountersByLanguage {
		ctx.AveragesByLanguage[language] = counters.average(ctx.NumberOfFiles)
	}

	return &ctx.CodeSummary, nil
}

func (ctx *context) visitPath(path string, info fs.FileInfo) error {

	if info.IsDir() {
		if ctx.isExcluded(path) {
			ctx.verboseLog("--- dir '%v' is excluded by patterns", path)
			return filepath.SkipDir
		}
		return nil
	}
	if !info.Mode().IsRegular() {
		ctx.verboseLog("--- file '%v' is not regular", path)
		return nil
	}
	if info.Size() > ctx.maxFileSizeBytes {
		ctx.verboseLog("--- file '%v' is too large (%v MB)", path, info.Size()/(1024*1024))
		return nil
	}

	fileExtension := filepath.Ext(path)
	language, matched := tryGetLanguage(fileExtension)
	if !matched {
		ctx.verboseLog("--- file '%v' was not mapped to any supported language", path)
		return nil
	}

	if ctx.isExcluded(path) || !ctx.isIncluded(path) {
		ctx.verboseLog("--- file '%v' is not matching patterns", path)
		return nil
	}

	counters, err := ctx.getCounters(path, language)
	if err != nil {
		return fmt.Errorf("failed to count at %v: %v", path, err)
	}
	ctx.verboseLog("+++ '%v': %v", path, counters)

	totalCounters, found := ctx.CountersByLanguage[language]
	if !found {
		totalCounters = &CodeCounters{}
		ctx.CountersByLanguage[language] = totalCounters
	}
	totalCounters.inc(counters)
	ctx.NumberOfFiles++

	return nil
}

func (ctx *context) getCounters(path string, language Language) (*CodeCounters, error) {
	content, err := ctx.readFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")

	counters := &CodeCounters{}

	minIndentations := float64(0)
	expectEndingComment := ""
	for _, line := range lines {
		counters.Lines++

		if len(expectEndingComment) > 0 {
			// in comment block
			if strings.Contains(line, expectEndingComment) {
				// comment block ended on this line
				expectEndingComment = ""
			}
			continue
		}

		cleanLine := strings.TrimSpace(line)
		if len(cleanLine) == 0 {
			continue
		}

		if strings.HasPrefix(cleanLine, "//") || strings.HasPrefix(cleanLine, "#") {
			// single line comment
			continue
		}

		if strings.HasPrefix(cleanLine, "/*") {
			expectEndingComment = "*/"
			cleanLine = cleanLine[2:]
		} else if language == "python" && strings.HasPrefix(cleanLine, "\"\"\"") {
			expectEndingComment = "\"\"\""
			cleanLine = cleanLine[3:]
		} else if language == "ruby" && strings.HasPrefix(cleanLine, "=begin") {
			expectEndingComment = "=end"
			continue
		} else if language == "ruby" && strings.HasPrefix(cleanLine, "<<-DOC") {
			expectEndingComment = "DOC"
			continue
		}

		if len(expectEndingComment) > 0 {
			// in comment block
			if strings.Contains(cleanLine, expectEndingComment) {
				// comment block ended on this line
				expectEndingComment = ""
			}
			continue
		}

		counters.LinesOfCode++

		counters.Indentations += float64(len(line) - len(trimSpaceLeft(line)))
		if minIndentations == 0 || minIndentations < counters.Indentations {
			minIndentations = counters.Indentations
		}

		counters.Keywords = countKeywords(cleanLine, language)
	}

	if minIndentations > 0 {
		counters.IndentationsNormalized = counters.Indentations / minIndentations
	}

	counters.IndentationsComplexity = safeDivide(counters.IndentationsNormalized, counters.LinesOfCode)
	counters.KeywordsComplexity = safeDivide(counters.Keywords, counters.LinesOfCode)

	return counters, nil
}

func (ctx *context) readFile(path string) (string, error) {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file at '%v': %v", path, err)
	}

	encoding, encodingName, _ := charset.DetermineEncoding(fileBytes, "")

	decodedBytes, err := encoding.NewDecoder().Bytes(fileBytes)
	if err != nil {
		return "", fmt.Errorf("failed to decode file at '%v' (detected as %v): %v", path, encodingName, err)
	}

	return string(decodedBytes), nil
}

func (ctx *context) isExcluded(path string) bool {
	if len(ctx.excludePatterns) > 0 && !matches(path, ctx.excludePatterns) {
		return true
	}
	return false
}

func (ctx *context) isIncluded(path string) bool {
	if len(ctx.includePatterns) == 0 || matches(path, ctx.includePatterns) {
		return true
	}
	return false
}

func (ctx *context) verboseLog(format string, v ...interface{}) {
	if ctx.verboseLogging {
		log.Printf(format, v...)
	}
}
