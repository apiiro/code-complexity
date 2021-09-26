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

func newContext() *context {
	return &context{
		CodeSummary: CodeSummary{
			CountersByLanguage: make(map[Language]*SummaryCounters),
		},
	}
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
	ctx := newContext()
	ctx.includePatterns = includePatterns
	ctx.excludePatterns = excludePatterns
	ctx.verboseLogging = opts.VerboseLogging
	ctx.maxFileSizeBytes = opts.MaxFileSizeBytes

	err = filepath.Walk(
		opts.CodePath,
		func(path string, info fs.FileInfo, _ error) error {
			return ctx.visitPath(opts.CodePath, path, info)
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to walk files under '%v': %v", opts.CodePath, err)
	}

	for _, counters := range ctx.CountersByLanguage {
		counters.Average = counters.Total.average(counters.NumberOfFiles)
	}

	return &ctx.CodeSummary, nil
}

func (ctx *context) visitPath(rootPath string, path string, info fs.FileInfo) error {

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

	relativePath, err := filepath.Rel(rootPath, path)
	if err != nil {
		return fmt.Errorf("failed to relativize path %v: %v", path, err)
	}
	if ctx.isExcluded(relativePath) || !ctx.isIncluded(relativePath) {
		ctx.verboseLog("--- file '%v' is not matching patterns", path)
		return nil
	}

	fileCounters, err := ctx.getCountersForPath(path, language)
	if err != nil {
		return fmt.Errorf("failed to count at %v: %v", path, err)
	}
	ctx.verboseLog("+++ '%v': %v", path, fileCounters)

	summaryCounters, found := ctx.CountersByLanguage[language]
	if !found {
		summaryCounters = &SummaryCounters{
			Total:   &CodeCounters{},
			Average: &CodeCounters{},
		}
		ctx.CountersByLanguage[language] = summaryCounters
	}
	summaryCounters.Total.inc(fileCounters)
	summaryCounters.NumberOfFiles++

	return nil
}

func (ctx *context) getCountersForPath(path string, language Language) (*CodeCounters, error) {
	content, err := ctx.readFile(path)
	if err != nil {
		return nil, err
	}
	return ctx.getCountersForCode(content, language)
}

func (ctx *context) getCountersForCode(content string, language Language) (*CodeCounters, error) {

	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")

	counters := &CodeCounters{}

	minIndentation := float64(0)
	prevIndentation := float64(-1)
	expectEndingComment := ""
	for _, line := range lines {
		counters.Lines++

		if len(expectEndingComment) > 0 {
			// in comment block
			endCommentIndex := strings.Index(line, expectEndingComment)
			if endCommentIndex == -1 {
				continue
			} else {
				// comment block ended on this line
				line = strings.TrimSpace(line[(endCommentIndex + len(expectEndingComment)):])
				expectEndingComment = ""
			}
		}

		cleanLine := strings.TrimSpace(line)
		if len(cleanLine) == 0 {
			continue
		}

		if strings.HasPrefix(cleanLine, "//") || strings.HasPrefix(cleanLine, "#") {
			// single line comment
			continue
		}

		const pythonMultilineString = "\"\"\""
		postCommentLine := ""
		if strings.Contains(cleanLine, "/*") {
			expectEndingComment = "*/"
			commentIndex := strings.Index(cleanLine, "/*")
			postCommentLine = strings.TrimSpace(cleanLine[2+commentIndex:])
			cleanLine = strings.TrimSpace(cleanLine[:commentIndex])
		} else if language == "python" && strings.Contains(cleanLine, pythonMultilineString) {
			expectEndingComment = pythonMultilineString
			commentIndex := strings.Index(cleanLine, pythonMultilineString)
			postCommentLine = strings.TrimSpace(cleanLine[len(pythonMultilineString)+commentIndex:])
			cleanLine = strings.TrimSpace(cleanLine[:commentIndex])
		} else if language == "ruby" && strings.HasPrefix(cleanLine, "=begin") {
			expectEndingComment = "=end"
			continue
		} else if language == "ruby" && strings.HasPrefix(cleanLine, "<<-DOC") {
			expectEndingComment = "DOC"
			continue
		}

		if len(postCommentLine) > 0 {
			// in comment block
			endCommentIndex := strings.Index(postCommentLine, expectEndingComment)
			if endCommentIndex == -1 {
				continue
			} else {
				// comment block ended on this line
				expectEndingComment = ""
			}
		}

		if len(cleanLine) == 0 {
			continue
		}

		counters.LinesOfCode++

		indentation := float64(len(line) - len(trimSpaceLeft(line)))
		if indentation > 0 {
			counters.Indentations += indentation
			if minIndentation == 0 || indentation < minIndentation {
				minIndentation = indentation
			}
			if prevIndentation != -1 {
				indentationDiff := indentation - prevIndentation
				if indentationDiff > 0 {
					counters.IndentationsDiff += indentationDiff
				}
			}
		}

		prevIndentation = indentation

		counters.Keywords += countKeywords(cleanLine, language)
	}

	if minIndentation > 0 {
		counters.IndentationsNormalized = counters.Indentations / minIndentation
		counters.IndentationsDiffNormalized = counters.IndentationsDiff / minIndentation
	}

	counters.IndentationsComplexity = safeDivide(counters.IndentationsNormalized, counters.LinesOfCode)
	counters.IndentationsDiffComplexity = safeDivide(counters.IndentationsDiffNormalized, counters.LinesOfCode)
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
	if len(ctx.excludePatterns) > 0 && matches(path, ctx.excludePatterns) {
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
