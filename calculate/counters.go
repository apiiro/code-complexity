package calculate

import "fmt"

type CodeSummary struct {
	NumberOfFiles      float64                       `json:"number_of_files"`
	CountersByLanguage map[Language]*SummaryCounters `json:"counters_by_language"`
}

type SummaryCounters struct {
	Total   *CodeCounters `json:"total"`
	Average *CodeCounters `json:"average"`
}

type CodeCounters struct {
	Lines                      float64 `json:"-"`
	LinesOfCode                float64 `json:"lines_of_code"`
	Keywords                   float64 `json:"-"`
	Indentations               float64 `json:"-"`
	IndentationsNormalized     float64 `json:"-"`
	IndentationsDiff           float64 `json:"-"`
	IndentationsDiffNormalized float64 `json:"-"`
	KeywordsComplexity         float64 `json:"keywords_complexity"`
	IndentationsComplexity     float64 `json:"indentations_complexity"`
	IndentationsDiffComplexity float64 `json:"indentations_diff_complexity"`
}

func (counters *CodeCounters) inc(other *CodeCounters) {
	counters.Lines += other.Lines
	counters.LinesOfCode += other.LinesOfCode
	counters.Keywords += other.Keywords
	counters.Indentations += other.Indentations
	counters.IndentationsNormalized += other.IndentationsNormalized
	counters.IndentationsDiff += other.IndentationsDiff
	counters.IndentationsDiffNormalized += other.IndentationsDiffNormalized
	counters.KeywordsComplexity += other.KeywordsComplexity
	counters.IndentationsComplexity += other.IndentationsComplexity
	counters.IndentationsDiffComplexity += other.IndentationsDiffComplexity
}

func (counters *CodeCounters) average(by float64) *CodeCounters {
	averaged := &CodeCounters{}
	if by == 0 {
		return averaged
	}
	averaged.Lines = counters.Lines / by
	averaged.LinesOfCode = counters.LinesOfCode / by
	averaged.Keywords = counters.Keywords / by
	averaged.Indentations = counters.Indentations / by
	averaged.IndentationsNormalized = counters.IndentationsNormalized / by
	averaged.IndentationsDiff = counters.IndentationsDiff / by
	averaged.IndentationsDiffNormalized = counters.IndentationsDiffNormalized / by
	averaged.KeywordsComplexity = counters.KeywordsComplexity / by
	averaged.IndentationsComplexity = counters.IndentationsComplexity / by
	averaged.IndentationsDiffComplexity = counters.IndentationsDiffComplexity / by
	return averaged
}

func (counters *CodeCounters) String() string {
	return fmt.Sprintf("loc=%v,Keywords=%v,indent=%v", counters.LinesOfCode, counters.Keywords, counters.IndentationsNormalized)
}

func safeDivide(a float64, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}
