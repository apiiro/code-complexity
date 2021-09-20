package calculate

import "fmt"

type CodeSummary struct {
	CountersByLanguage map[Language]*CodeCounters
	AveragesByLanguage map[Language]*CodeCounters
	NumberOfFiles      float64
}

type CodeCounters struct {
	Lines                  float64
	LinesOfCode            float64
	Keywords               float64
	Indentations           float64
	IndentationsNormalized float64
	KeywordsComplexity     float64
	IndentationsComplexity float64
}

func (counters *CodeCounters) inc(other *CodeCounters) {
	counters.Lines += other.Lines
	counters.LinesOfCode += other.LinesOfCode
	counters.Keywords += other.Keywords
	counters.Indentations += other.Indentations
	counters.IndentationsNormalized += other.IndentationsNormalized
	counters.KeywordsComplexity += other.KeywordsComplexity
	counters.IndentationsComplexity += other.IndentationsComplexity
}

func (counters *CodeCounters) average(by float64) *CodeCounters {
	averaged := &CodeCounters{}
	if by == 0 {
		return averaged
	}
	averaged.Lines = counters.Lines / by
	averaged.Lines = counters.LinesOfCode / by
	averaged.Lines = counters.Keywords / by
	averaged.Lines = counters.Indentations / by
	averaged.Lines = counters.IndentationsNormalized / by
	averaged.Lines = counters.KeywordsComplexity / by
	averaged.Lines = counters.IndentationsComplexity / by
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
