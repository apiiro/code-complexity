package calculate

import (
	"code-complexity/test_resources"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

// test_resources patterns, encoding

func getCountersForCode(code string, language Language) (*CodeCounters, error) {
	ctx := newContext()
	return ctx.getCountersForCode(code, language)
}

func TestCountersForEmptyInput(t *testing.T) {
	r := require.New(t)

	counters, err := getCountersForCode("", "java")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.Lines)
	r.Equal(float64(0), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), counters.IndentationsNormalized)
	r.Equal(float64(0), counters.IndentationsDiff)
	r.Equal(float64(0), counters.IndentationsDiffNormalized)
	r.Equal(float64(0), counters.KeywordsComplexity)
	r.Equal(float64(0), counters.IndentationsComplexity)
	r.Equal(float64(0), counters.IndentationsDiffComplexity)
}

func TestCountersForJava(t *testing.T) {
	r := require.New(t)

	// language=java
	code := `
// comment
int x = 3;
/* another comment */
`
	counters, err := getCountersForCode(code, "java")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=java
	code = `
/*
multiline comment
*/
int x = 3;
`
	counters, err = getCountersForCode(code, "java")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=java
	code = `
if (x > 3) {
	x = 4;
}
else {
	if (x > 7) {
			x = 8;
	}
}
`
	counters, err = getCountersForCode(code, "java")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(8), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(6), counters.Indentations)
	r.Equal(float64(6), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(4), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	// language=java
	code = `
if (x > 3) {
  x = 4;
}
else {
  if (x > 7) {
     x = 8;
   }
}
`
	counters, err = getCountersForCode(code, "java")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(8), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(12), counters.Indentations)
	r.Equal(float64(6), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(7), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	// language=java
	code = `
public class Classic {
	public Classic() {
		for (var i = 0; i < 10; i++) {
			if (i % 2 == 0) {
				System.out.println(String.format("%d", i)
			}
		}
	}
}
`
	counters, err = getCountersForCode(code, "java")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(9), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(12), counters.Indentations)
	r.Equal(float64(6), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(7), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))
}

func TestCountersForJavaFullSample(t *testing.T) {
	r := require.New(t)

	counters, err := getCountersForCode(test_resources.JavaCode, "java")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(624), counters.Lines)
	r.Equal(float64(448), counters.LinesOfCode)
	r.Equal(float64(104), counters.Keywords)
	r.Equal(float64(3650), counters.Indentations)
	r.Equal(float64(913), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(387), math.Round(counters.IndentationsDiff))
	r.Equal(float64(97), math.Round(counters.IndentationsDiffNormalized))
	r.Equal(float64(23), math.Round(counters.KeywordsComplexity * 100))
	r.Equal(float64(204), math.Round(counters.IndentationsComplexity * 100))
	r.Equal(float64(22), math.Round(counters.IndentationsDiffComplexity * 100))
}
