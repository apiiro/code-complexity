package calculate

import (
	"code-complexity/options"
	"code-complexity/test_resources"
	"io/fs"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getFileCount(basePath string, includes []string, excludes []string) (float64, error) {
	opts := &options.Options{
		CodePath:        basePath,
		IncludePatterns: includes,
		ExcludePatterns: excludes,
		VerboseLogging:  true,
	}
	summary, err := Complexity(opts)
	if err != nil {
		return 0, err
	}
	totalNumberOfFiles := float64(0)
	for _, counters := range summary.CountersByLanguage {
		totalNumberOfFiles += counters.NumberOfFiles
	}
	return totalNumberOfFiles, err
}

func mkdir(path string) {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		panic(err)
	}
}

func touch(filePath string) {
	err := os.WriteFile(filePath, []byte{}, 0777)
	if err != nil {
		panic(err)
	}
}

func TestIncludeExcludePatterns(t *testing.T) {
	r := assert.New(t)

	basePath, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := os.RemoveAll(basePath)
		if err != nil {
			panic(err)
		}
	}()

	mkdir(filepath.Join(basePath, "src"))
	mkdir(filepath.Join(basePath, "src", "nested"))
	touch(filepath.Join(basePath, "root.java"))
	touch(filepath.Join(basePath, "a.js"))
	touch(filepath.Join(basePath, "b.js"))
	touch(filepath.Join(basePath, "src", "svc.java"))
	touch(filepath.Join(basePath, "src", "api.js"))
	touch(filepath.Join(basePath, "src", "nested", "util.js"))

	filesCount, err := getFileCount(
		basePath,
		[]string{
			"**/*.java",
			"a.js",
		},
		[]string{},
	)
	r.Nil(err)
	r.Equal(float64(3), filesCount)

	filesCount, err = getFileCount(
		basePath,
		[]string{},
		[]string{
			"**/*.java",
		},
	)
	r.Nil(err)
	r.Equal(float64(4), filesCount)

	filesCount, err = getFileCount(
		basePath,
		[]string{},
		[]string{
			"*/*.java",
			"a.js",
		},
	)
	r.Nil(err)
	r.Equal(float64(3), filesCount)

	filesCount, err = getFileCount(
		basePath,
		[]string{
			"**/*.js",
		},
		[]string{
			"**/nested/**",
		},
	)
	r.Nil(err)
	r.Equal(float64(3), filesCount)
}

func TestEncodings(t *testing.T) {
	r := require.New(t)

	wdPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	sourcePath := filepath.Join(wdPath, "..", "test_resources", "encoding")

	basePath, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := os.RemoveAll(basePath)
		if err != nil {
			panic(err)
		}
	}()

	err = copy.Copy(sourcePath, basePath)
	if err != nil {
		panic(err)
	}
	err = filepath.Walk(basePath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		return os.Rename(path, strings.Replace(path, ".txt", ".go", 1))
	})
	if err != nil {
		panic(err)
	}

	opts := &options.Options{
		CodePath:         basePath,
		IncludePatterns:  []string{},
		ExcludePatterns:  []string{},
		VerboseLogging:   true,
		MaxFileSizeBytes: 1024 * 1024,
	}
	summary, err := Complexity(opts)
	r.Nil(err)
	r.Equal(float64(3), summary.CountersByLanguage["go"].NumberOfFiles)
	r.Equal(float64(5*3), summary.CountersByLanguage["go"].Total.LinesOfCode)
}

func inRange(r *assert.Assertions, value float64, min int, max int) {
	r.GreaterOrEqual(value, float64(min))
	r.LessOrEqual(value, float64(max))
}

func TestDogFood(t *testing.T) {
	r := assert.New(t)

	wdPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	basePath := filepath.Join(wdPath, "..")

	opts := &options.Options{
		CodePath:        basePath,
		IncludePatterns: []string{},
		ExcludePatterns: []string{
			"test_resources/**",
			".git/**",
			".idea/**",
		},
		VerboseLogging:   true,
		MaxFileSizeBytes: 1024 * 1024,
	}
	summary, err := Complexity(opts)
	r.Nil(err)

	r.Len(summary.CountersByLanguage, 1)

	r.Equal(float64(9), summary.CountersByLanguage["go"].NumberOfFiles)

	total := summary.CountersByLanguage["go"].Total
	inRange(r, total.Lines, 2000, 4000)
	inRange(r, total.LinesOfCode, 2000, 4000)
	inRange(r, total.Keywords, 200, 400)
	inRange(r, total.Indentations, 2500, 3500)
	inRange(r, total.IndentationsNormalized, 2500, 3500)
	inRange(r, total.IndentationsDiff, 400, 600)
	inRange(r, total.IndentationsDiffNormalized, 400, 600)
	inRange(r, total.IndentationsComplexity, 10, 12)
	inRange(r, total.IndentationsDiffComplexity*100, 150, 250)
	inRange(r, total.KeywordsComplexity*100, 200, 250)

	average := summary.CountersByLanguage["go"].Average
	inRange(r, average.Lines, 300, 400)
	inRange(r, average.LinesOfCode, 250, 300)
	inRange(r, average.Keywords, 25, 35)
	inRange(r, average.Indentations, 300, 350)
	inRange(r, average.IndentationsNormalized, 300, 350)
	inRange(r, average.IndentationsDiff, 50, 60)
	inRange(r, average.IndentationsDiffNormalized, 50, 60)
	inRange(r, average.IndentationsComplexity, 1, 2)
	inRange(r, average.IndentationsDiffComplexity*100, 20, 30)
	inRange(r, average.KeywordsComplexity*100, 20, 30)
}

func getCountersForCode(code string, language Language) (*CodeCounters, error) {
	ctx := newContext()
	return ctx.getCountersForCode(code, language)
}

func TestCountersForEmptyInput(t *testing.T) {
	r := assert.New(t)

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
	r := assert.New(t)

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
int x= 3; /*
multiline comment
*/ int y = 4;
`
	counters, err = getCountersForCode(code, "java")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(2), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=java
	code = `
const path = "dir/*.ext";
int x = 1;
int y = 2;
`
	counters, err = getCountersForCode(code, "java")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(3), counters.LinesOfCode)
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
	r.Equal(float64(16), counters.Indentations)
	r.Equal(float64(16), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(4), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))
}

func TestCountersForJavaFullSample(t *testing.T) {
	r := assert.New(t)

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
	r.Equal(float64(23), math.Round(counters.KeywordsComplexity*100))
	r.Equal(float64(204), math.Round(counters.IndentationsComplexity*100))
	r.Equal(float64(22), math.Round(counters.IndentationsDiffComplexity*100))
}

func TestCountersForCSharp(t *testing.T) {
	r := assert.New(t)

	// language=cs
	code := `
// comment
int x = 3;
/* another comment */
`
	counters, err := getCountersForCode(code, "csharp")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=cs
	code = `
/*
multiline comment
*/
int x = 3;
`
	counters, err = getCountersForCode(code, "csharp")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=cs
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
	counters, err = getCountersForCode(code, "csharp")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(8), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(6), counters.Indentations)
	r.Equal(float64(6), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(4), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	// language=cs
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
	counters, err = getCountersForCode(code, "csharp")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(8), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(12), counters.Indentations)
	r.Equal(float64(6), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(7), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	// language=cs
	code = `
public class Classic {
	public Classic() {
		for (var i = 0; i < 10; i++) {
			if (i % 2 == 0) {
				Console.WriteLine($"{i}")
			}
		}
	}
}
`
	counters, err = getCountersForCode(code, "csharp")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(9), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(16), counters.Indentations)
	r.Equal(float64(16), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(4), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))
}

func TestCountersForCSharpFullSample(t *testing.T) {
	r := assert.New(t)

	counters, err := getCountersForCode(test_resources.CSharpCode, "csharp")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(775), counters.Lines)
	r.Equal(float64(584), counters.LinesOfCode)
	r.Equal(float64(122), counters.Keywords)
	r.Equal(float64(8236), counters.Indentations)
	r.Equal(float64(2059), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(484), math.Round(counters.IndentationsDiff))
	r.Equal(float64(121), math.Round(counters.IndentationsDiffNormalized))
	r.Equal(float64(21), math.Round(counters.KeywordsComplexity*100))
	r.Equal(float64(353), math.Round(counters.IndentationsComplexity*100))
	r.Equal(float64(21), math.Round(counters.IndentationsDiffComplexity*100))
}

func TestCountersForNode(t *testing.T) {
	r := assert.New(t)

	// language=js
	code := `
// comment
const x = 3;
/* another comment */
`
	counters, err := getCountersForCode(code, "node")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=js
	code = `
/*
multiline comment
*/
let x = 3
`
	counters, err = getCountersForCode(code, "node")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=js
	code = "const a = `${b}`"
	counters, err = getCountersForCode(code, "node")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=js
	code = `
if (x > 3) {
	x = 4;
}
else {
	if (x === 7) {
			x = 8
	}
}
`
	counters, err = getCountersForCode(code, "node")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(8), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(6), counters.Indentations)
	r.Equal(float64(6), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(4), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	// language=js
	code = `
if (x > 3) {
  x = 4;
}
else {
  if (x === 7) {
     x = 8
   }
}
`
	counters, err = getCountersForCode(code, "node")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(8), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(12), counters.Indentations)
	r.Equal(float64(6), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(7), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	// language=js
	code = `
export class Classic {
	constructor() {
		for (var i = 0; i < 10; i++) {
			if (i % 2 === 0) {
				console.log("${i}")
			}
		}
	}
}
`
	counters, err = getCountersForCode(code, "node")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(9), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(16), counters.Indentations)
	r.Equal(float64(16), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(4), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	// language=ts
	code = `
interface User {
  name: string;
  id: number;
}
 
class UserAccount {
  name: string;
  id: number;
 
  constructor(name: string, id: number) {
    this.name = name;
    this.id = id;
  }
}
 
const user: User = new UserAccount("Murphy", 1);
`
	counters, err = getCountersForCode(code, "node")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(13), counters.LinesOfCode)
	r.Equal(float64(2), counters.Keywords)
	r.Equal(float64(20), counters.Indentations)
	r.Equal(float64(10), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(6), math.Round(counters.IndentationsDiff))
	r.Equal(float64(3), math.Round(counters.IndentationsDiffNormalized))
}

func TestCountersForNodeFullSample(t *testing.T) {
	r := assert.New(t)

	counters, err := getCountersForCode(test_resources.NodeCode, "node")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(414), counters.Lines)
	r.Equal(float64(287), counters.LinesOfCode)
	r.Equal(float64(112), counters.Keywords)
	r.Equal(float64(711), counters.Indentations)
	r.Equal(float64(711), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(60), math.Round(counters.IndentationsDiff))
	r.Equal(float64(60), math.Round(counters.IndentationsDiffNormalized))
	r.Equal(float64(39), math.Round(counters.KeywordsComplexity*100))
	r.Equal(float64(248), math.Round(counters.IndentationsComplexity*100))
	r.Equal(float64(21), math.Round(counters.IndentationsDiffComplexity*100))
}

func TestCountersForPython(t *testing.T) {
	r := assert.New(t)

	// language=py
	code := `
// comment
x = 3
/* another comment */
`
	counters, err := getCountersForCode(code, "python")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=py
	code = `
/*
multiline comment
*/
global x
`
	counters, err = getCountersForCode(code, "python")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=py
	code = `
"""
multiline comment
"""
global x
`
	counters, err = getCountersForCode(code, "python")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=py
	code = `
var x = """
some long text
"""
var x = 3 """
some long comment
"""
`
	counters, err = getCountersForCode(code, "python")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(2), counters.LinesOfCode) // this isn't completely true -- should count the assigned string as content too, but cest la vi
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=py
	code = `
if x > 3:
	x = 4
elif (x == 7):
			x = 8
else:
	x = x if x else x
`
	counters, err = getCountersForCode(code, "python")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(6), counters.LinesOfCode)
	r.Equal(float64(4), counters.Keywords)
	r.Equal(float64(5), counters.Indentations)
	r.Equal(float64(5), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(5), math.Round(counters.IndentationsDiff))
	r.Equal(float64(5), math.Round(counters.IndentationsDiffNormalized))

	// language=py
	code = `
if x > 3:
  x = 4
elif (x === 7):
     x = 8
`
	counters, err = getCountersForCode(code, "python")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(4), counters.LinesOfCode)
	r.Equal(float64(2), counters.Keywords)
	r.Equal(float64(7), counters.Indentations)
	r.Equal(float64(4), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(7), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	// language=py
	code = `
class Animal(models.Model):
    name = models.CharField(max_length=150)
    latin_name = models.CharField(max_length=150)
    count = models.IntegerField()
    weight = models.FloatField()

    # use a non-default name for the default manager
    specimens = models.Manager()

    def __str__(self):
        return self.name
`
	counters, err = getCountersForCode(code, "python")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(8), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(32), counters.Indentations)
	r.Equal(float64(8), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(8), math.Round(counters.IndentationsDiff))
	r.Equal(float64(2), math.Round(counters.IndentationsDiffNormalized))
}

func TestCountersForPythonFullSample(t *testing.T) {
	r := assert.New(t)

	counters, err := getCountersForCode(test_resources.PythonCode, "python")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(240), counters.Lines)
	r.Equal(float64(146), counters.LinesOfCode)
	r.Equal(float64(79), counters.Keywords)
	r.Equal(float64(1104), counters.Indentations)
	r.Equal(float64(276), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(224), math.Round(counters.IndentationsDiff))
	r.Equal(float64(56), math.Round(counters.IndentationsDiffNormalized))
	r.Equal(float64(54), math.Round(counters.KeywordsComplexity*100))
	r.Equal(float64(189), math.Round(counters.IndentationsComplexity*100))
	r.Equal(float64(38), math.Round(counters.IndentationsDiffComplexity*100))
}

func TestCountersForKotlin(t *testing.T) {
	r := assert.New(t)

	// language=kt
	code := `
// comment
internal const val x = 3
/* another comment */
`
	counters, err := getCountersForCode(code, "kotlin")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=kt
	code = `
/*
multiline comment
*/
val x = "kt"
`
	counters, err = getCountersForCode(code, "kotlin")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=kt
	code = `
if (x > 3) {
	x = 4
}
else {
	if (x > 7) {
			x = 8
	}
}
`
	counters, err = getCountersForCode(code, "kotlin")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(8), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(6), counters.Indentations)
	r.Equal(float64(6), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(4), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	// language=kt
	code = `
if (x > 3) {
  x = 4
}
else {
  if (x > 7) {
     x = 8
   }
}
`
	counters, err = getCountersForCode(code, "kotlin")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(8), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(12), counters.Indentations)
	r.Equal(float64(6), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(7), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	// language=kt
	code = `
class Pet {
    constructor(owner: Person) {
        owner.pets.add(this) // adds this pet to the list of its owner's pets
		println("owner: $owner")
    }
}
`
	counters, err = getCountersForCode(code, "kotlin")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(6), counters.LinesOfCode)
	r.Equal(float64(1), counters.Keywords)
	r.Equal(float64(18), counters.Indentations)
	r.Equal(float64(9), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(10), math.Round(counters.IndentationsDiff))
	r.Equal(float64(5), math.Round(counters.IndentationsDiffNormalized))
}

func TestCountersForKotlinFullSample(t *testing.T) {
	r := assert.New(t)

	counters, err := getCountersForCode(test_resources.KotlinCode, "kotlin")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(183), counters.Lines)
	r.Equal(float64(125), counters.LinesOfCode)
	r.Equal(float64(62), counters.Keywords)
	r.Equal(float64(592), counters.Indentations)
	r.Equal(float64(148), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(128), math.Round(counters.IndentationsDiff))
	r.Equal(float64(32), math.Round(counters.IndentationsDiffNormalized))
	r.Equal(float64(50), math.Round(counters.KeywordsComplexity*100))
	r.Equal(float64(118), math.Round(counters.IndentationsComplexity*100))
	r.Equal(float64(26), math.Round(counters.IndentationsDiffComplexity*100))
}

func TestCountersForScala(t *testing.T) {
	r := assert.New(t)

	code := `
// comment
val x = 3
/* another comment */
`
	counters, err := getCountersForCode(code, "scala")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	code = `
/*
multiline comment
*/
val x = "kt"
`
	counters, err = getCountersForCode(code, "scala")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	code = `
if (x > 3) {
	x = 4
}
else {
	if (x > 7) {
			x = 8
	}
}
`
	counters, err = getCountersForCode(code, "scala")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(8), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(6), counters.Indentations)
	r.Equal(float64(6), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(4), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	code = `
if (x > 3) {
  x = 4
}
else {
  if (x > 7) {
     x = 8
   }
}
`
	counters, err = getCountersForCode(code, "scala")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(8), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(12), counters.Indentations)
	r.Equal(float64(6), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(7), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	code = `
class Point(var x: Int, var y: Int) {

  def move(dx: Int, dy: Int): Unit = {
    x = x + dx
    y = y + dy
  }

  override def toString: String =
    s"($x, $y)"
}
`
	counters, err = getCountersForCode(code, "scala")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(8), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(18), counters.Indentations)
	r.Equal(float64(9), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(6), math.Round(counters.IndentationsDiff))
	r.Equal(float64(3), math.Round(counters.IndentationsDiffNormalized))
}

func TestCountersForScalaFullSample(t *testing.T) {
	r := assert.New(t)

	counters, err := getCountersForCode(test_resources.ScalaCode, "scala")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(638), counters.Lines)
	r.Equal(float64(500), counters.LinesOfCode)
	r.Equal(float64(199), counters.Keywords)
	r.Equal(float64(2374), counters.Indentations)
	r.Equal(float64(1187), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(327), math.Round(counters.IndentationsDiff))
	r.Equal(float64(164), math.Round(counters.IndentationsDiffNormalized))
	r.Equal(float64(40), math.Round(counters.KeywordsComplexity*100))
	r.Equal(float64(237), math.Round(counters.IndentationsComplexity*100))
	r.Equal(float64(33), math.Round(counters.IndentationsDiffComplexity*100))
}

func TestCountersFoC(t *testing.T) {
	r := assert.New(t)

	// language=c
	code := `
// comment
int x = 3;
/* another comment */
`
	counters, err := getCountersForCode(code, "c")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=c
	code = `
/*
multiline comment
*/
char str[] = "str";
`
	counters, err = getCountersForCode(code, "c")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=c
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
	counters, err = getCountersForCode(code, "c")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(8), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(6), counters.Indentations)
	r.Equal(float64(6), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(4), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	// language=c
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
	counters, err = getCountersForCode(code, "c")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(8), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(12), counters.Indentations)
	r.Equal(float64(6), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(7), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	// language=c
	code = `
struct v {
   union { // anonymous union
      struct { int i, j; }; // anonymous structure
      struct { long k, l; } w;
   };
   int m;
} v1;
`
	counters, err = getCountersForCode(code, "c")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(7), counters.LinesOfCode)
	r.Equal(float64(4), counters.Keywords)
	r.Equal(float64(21), counters.Indentations)
	r.Equal(float64(7), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(6), math.Round(counters.IndentationsDiff))
	r.Equal(float64(2), math.Round(counters.IndentationsDiffNormalized))
}

func TestCountersForCFullSample(t *testing.T) {
	r := assert.New(t)

	counters, err := getCountersForCode(test_resources.CCode, "c")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(693), counters.Lines)
	r.Equal(float64(566), counters.LinesOfCode)
	r.Equal(float64(187), counters.Keywords)
	r.Equal(float64(922), counters.Indentations)
	r.Equal(float64(922), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(228), math.Round(counters.IndentationsDiff))
	r.Equal(float64(228), math.Round(counters.IndentationsDiffNormalized))
	r.Equal(float64(33), math.Round(counters.KeywordsComplexity*100))
	r.Equal(float64(163), math.Round(counters.IndentationsComplexity*100))
	r.Equal(float64(40), math.Round(counters.IndentationsDiffComplexity*100))
}

func TestCountersFoCpp(t *testing.T) {
	r := assert.New(t)

	// language=cpp
	code := `
// comment
int x = 3;
/* another comment */
`
	counters, err := getCountersForCode(code, "cpp")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=cpp
	code = `
/*
* multiline comment
*/
char str[] = "str";
`
	counters, err = getCountersForCode(code, "cpp")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=cpp
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
	counters, err = getCountersForCode(code, "cpp")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(8), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(6), counters.Indentations)
	r.Equal(float64(6), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(4), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	// language=cpp
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
	counters, err = getCountersForCode(code, "cpp")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(8), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(12), counters.Indentations)
	r.Equal(float64(6), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(7), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	// language=cpp
	code = `
#include <iostream>
using namespace std;
  
template <typename T>
class Array {
private:
    T* ptr;
    int size;
  
public:
    Array(T arr[], int s);
    void print();
};
`
	counters, err = getCountersForCode(code, "cpp")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(10), counters.LinesOfCode)
	r.Equal(float64(4), counters.Keywords)
	r.Equal(float64(16), counters.Indentations)
	r.Equal(float64(4), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(8), math.Round(counters.IndentationsDiff))
	r.Equal(float64(2), math.Round(counters.IndentationsDiffNormalized))
}

func TestCountersForCppFullSample(t *testing.T) {
	r := assert.New(t)

	counters, err := getCountersForCode(test_resources.CppCode, "cpp")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(369), counters.Lines)
	r.Equal(float64(240), counters.LinesOfCode)
	r.Equal(float64(48), counters.Keywords)
	r.Equal(float64(1336), counters.Indentations)
	r.Equal(float64(1336), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(344), math.Round(counters.IndentationsDiff))
	r.Equal(float64(344), math.Round(counters.IndentationsDiffNormalized))
	r.Equal(float64(20), math.Round(counters.KeywordsComplexity*100))
	r.Equal(float64(557), math.Round(counters.IndentationsComplexity*100))
	r.Equal(float64(143), math.Round(counters.IndentationsDiffComplexity*100))
}

func TestCountersForObjectivec(t *testing.T) {
	r := assert.New(t)

	// language=mm
	code := `
// comment
int x = 3;
/* another comment */
`
	counters, err := getCountersForCode(code, "objectivec")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=mm
	code = `
/*
 multiline comment
*/
char str[] = "str";
`
	counters, err = getCountersForCode(code, "objectivec")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=mm
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
	counters, err = getCountersForCode(code, "objectivec")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(8), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(6), counters.Indentations)
	r.Equal(float64(6), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(4), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	// language=mm
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
	counters, err = getCountersForCode(code, "objectivec")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(8), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(12), counters.Indentations)
	r.Equal(float64(6), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(7), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	// language=mm
	code = `
@interface XYZPerson : NSObject
- (void)sayHello;
@end

#import "XYZPerson.h"
 
@implementation XYZPerson
- (void)sayHello {
    NSLog( @"Hello, World!" );
}
@end
`
	counters, err = getCountersForCode(code, "objectivec")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(8), counters.LinesOfCode)
	r.Equal(float64(4), counters.Keywords)
	r.Equal(float64(4), counters.Indentations)
	r.Equal(float64(1), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(4), math.Round(counters.IndentationsDiff))
	r.Equal(float64(1), math.Round(counters.IndentationsDiffNormalized))
}

func TestCountersForOjbectivecFullSample(t *testing.T) {
	r := assert.New(t)

	counters, err := getCountersForCode(test_resources.ObjectivecCode, "cpp")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1404), counters.Lines)
	r.Equal(float64(1105), counters.LinesOfCode)
	r.Equal(float64(120), counters.Keywords)
	r.Equal(float64(1758), counters.Indentations)
	r.Equal(float64(1758), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(207), math.Round(counters.IndentationsDiff))
	r.Equal(float64(207), math.Round(counters.IndentationsDiffNormalized))
	r.Equal(float64(11), math.Round(counters.KeywordsComplexity*100))
	r.Equal(float64(159), math.Round(counters.IndentationsComplexity*100))
	r.Equal(float64(19), math.Round(counters.IndentationsDiffComplexity*100))
}

func TestCountersForSwift(t *testing.T) {
	r := assert.New(t)

	// language=swift
	code := `
// comment
var x: Int = 17
/* another comment */
`
	counters, err := getCountersForCode(code, "swift")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=swift
	code = `
/*
 multiline comment
*/
var x: Int = 17
`
	counters, err = getCountersForCode(code, "swift")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=swift
	code = `
if (x > 3) {
	x = 4
}
else {
	if (x > 7) {
			x = 8
	}
}
`
	counters, err = getCountersForCode(code, "swift")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(8), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(6), counters.Indentations)
	r.Equal(float64(6), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(4), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	// language=swift
	code = `
if (x > 3) {
  x = 4
}
else {
  if (x > 7) {
     x = 8
   }
}
`
	counters, err = getCountersForCode(code, "swift")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(8), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(12), counters.Indentations)
	r.Equal(float64(6), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(7), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	// language=swift
	code = `
public class Person {
    private var _id: Int = 0
    private var _lastName: String = ""

    public init(id: Int, lastName: String) {
        self.id = id
        self.lastName = lastName
    }

    public var id: Int {
        get {
            return self._id;
        }
        set {
            if newValue < 0 || newValue > 1000 {
                // Swift setter cannot throw error.
                fatalError("invalid value for id")
            } else {
                self._id = newValue
            }
        }
    }
}
`
	counters, err = getCountersForCode(code, "swift")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(20), counters.LinesOfCode)
	r.Equal(float64(7), counters.Keywords)
	r.Equal(float64(152), counters.Indentations)
	r.Equal(float64(38), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(28), math.Round(counters.IndentationsDiff))
	r.Equal(float64(7), math.Round(counters.IndentationsDiffNormalized))
}

func TestCountersForSwiftFullSample(t *testing.T) {
	r := assert.New(t)

	counters, err := getCountersForCode(test_resources.SwiftCode, "swift")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(415), counters.Lines)
	r.Equal(float64(253), counters.LinesOfCode)
	r.Equal(float64(91), counters.Keywords)
	r.Equal(float64(2016), counters.Indentations)
	r.Equal(float64(504), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(280), math.Round(counters.IndentationsDiff))
	r.Equal(float64(70), math.Round(counters.IndentationsDiffNormalized))
	r.Equal(float64(36), math.Round(counters.KeywordsComplexity*100))
	r.Equal(float64(199), math.Round(counters.IndentationsComplexity*100))
	r.Equal(float64(28), math.Round(counters.IndentationsDiffComplexity*100))
}

func TestCountersForGo(t *testing.T) {
	r := assert.New(t)

	// language=swift
	code := `
// comment
x := 3
/* another comment */
`
	counters, err := getCountersForCode(code, "go")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	code = `
/*
 multiline comment
*/
x := "x"
`
	counters, err = getCountersForCode(code, "go")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	code = "var x:=`x`"
	counters, err = getCountersForCode(code, "go")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	code = `
if x > 3 {
	x = 4
} else {
	if x > 7 {
			x = 8
	}
}
`
	counters, err = getCountersForCode(code, "go")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(7), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(6), counters.Indentations)
	r.Equal(float64(6), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(4), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	code = `
if x > 3 {
  x = 4
} else {
  if x > 7 {
     x = 8
   }
}
`
	counters, err = getCountersForCode(code, "go")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(7), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(12), counters.Indentations)
	r.Equal(float64(6), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(7), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	code = `
type Person struct {
  Name string
}
func (p *Person) Talk() {
  fmt.Println("Hi, my name is", p.Name)
}
`
	counters, err = getCountersForCode(code, "go")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(6), counters.LinesOfCode)
	r.Equal(float64(2), counters.Keywords)
	r.Equal(float64(4), counters.Indentations)
	r.Equal(float64(2), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(4), math.Round(counters.IndentationsDiff))
	r.Equal(float64(2), math.Round(counters.IndentationsDiffNormalized))
}

func TestCountersForGoFullSample(t *testing.T) {
	r := assert.New(t)

	counters, err := getCountersForCode(test_resources.GoCode, "go")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(499), counters.Lines)
	r.Equal(float64(388), counters.LinesOfCode)
	r.Equal(float64(205), counters.Keywords)
	r.Equal(float64(671), counters.Indentations)
	r.Equal(float64(671), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(108), math.Round(counters.IndentationsDiff))
	r.Equal(float64(108), math.Round(counters.IndentationsDiffNormalized))
	r.Equal(float64(53), math.Round(counters.KeywordsComplexity*100))
	r.Equal(float64(173), math.Round(counters.IndentationsComplexity*100))
	r.Equal(float64(28), math.Round(counters.IndentationsDiffComplexity*100))
}

func TestCountersFoRust(t *testing.T) {
	r := assert.New(t)

	// language=rs
	code := `
// comment
let x: u32 = 4;
/* another comment */
`
	counters, err := getCountersForCode(code, "rust")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=rs
	code = `
/*
multiline comment
*/
let mut hello = String::from("Hello, ");
`
	counters, err = getCountersForCode(code, "rust")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=rs
	code = `
if x > 3 {
	x = 4;
}
else {
	if x > 7 {
			x = 8;
	}
}
`
	counters, err = getCountersForCode(code, "rust")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(8), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(6), counters.Indentations)
	r.Equal(float64(6), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(4), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	// language=rs
	code = `
if x > 3 {
  x = 4;
}
else {
  if x > 7 {
     x = 8;
   }
}
`
	counters, err = getCountersForCode(code, "rust")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(8), counters.LinesOfCode)
	r.Equal(float64(3), counters.Keywords)
	r.Equal(float64(12), counters.Indentations)
	r.Equal(float64(6), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(7), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	// language=rs
	code = `
struct User {
    username: &str,
    email: &str,
    sign_in_count: u64,
    active: bool,
}

fn main() {
    let user1 = User {
        email: "email",
        username: "username",
        active: true,
        sign_in_count: 1,
    };
}
`
	counters, err = getCountersForCode(code, "rust")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(14), counters.LinesOfCode)
	r.Equal(float64(2), counters.Keywords)
	r.Equal(float64(56), counters.Indentations)
	r.Equal(float64(14), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(12), math.Round(counters.IndentationsDiff))
	r.Equal(float64(3), math.Round(counters.IndentationsDiffNormalized))
}

func TestCountersForRustFullSample(t *testing.T) {
	r := assert.New(t)

	counters, err := getCountersForCode(test_resources.RustCode, "rust")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(202), counters.Lines)
	r.Equal(float64(143), counters.LinesOfCode)
	r.Equal(float64(34), counters.Keywords)
	r.Equal(float64(936), counters.Indentations)
	r.Equal(float64(234), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(144), math.Round(counters.IndentationsDiff))
	r.Equal(float64(36), math.Round(counters.IndentationsDiffNormalized))
	r.Equal(float64(24), math.Round(counters.KeywordsComplexity*100))
	r.Equal(float64(164), math.Round(counters.IndentationsComplexity*100))
	r.Equal(float64(25), math.Round(counters.IndentationsDiffComplexity*100))
}

func TestCountersFoRuby(t *testing.T) {
	r := assert.New(t)

	// language=rb
	code := `
// comment
x = 4
/* another comment */
`
	counters, err := getCountersForCode(code, "ruby")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=rb
	code = `
/*
multiline comment
*/
x = "x"
`
	counters, err = getCountersForCode(code, "ruby")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(1), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=rb
	code = `
#!/usr/bin/env ruby

=begin
Every body mentioned this way
to have multiline comments.
=end

puts "Hello world!"

<<-DOC
Also, you could create a docstring.
which...
DOC

puts "Hello world!"
`
	counters, err = getCountersForCode(code, "ruby")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(2), counters.LinesOfCode)
	r.Equal(float64(0), counters.Keywords)
	r.Equal(float64(0), counters.Indentations)
	r.Equal(float64(0), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(0), math.Round(counters.IndentationsDiff))
	r.Equal(float64(0), math.Round(counters.IndentationsDiffNormalized))

	// language=rb
	code = `
if x > 3 {
	x = 4;
}
elsif x > 7 {
			x = 8;
	}
}
`
	counters, err = getCountersForCode(code, "ruby")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(7), counters.LinesOfCode)
	r.Equal(float64(2), counters.Keywords)
	r.Equal(float64(5), counters.Indentations)
	r.Equal(float64(5), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(4), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	// language=rb
	code = `
if x > 3 {
  x = 4;
}
elsif x > 7 {
     x = 8;
   }
}
`
	counters, err = getCountersForCode(code, "ruby")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(7), counters.LinesOfCode)
	r.Equal(float64(2), counters.Keywords)
	r.Equal(float64(10), counters.Indentations)
	r.Equal(float64(5), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(7), math.Round(counters.IndentationsDiff))
	r.Equal(float64(4), math.Round(counters.IndentationsDiffNormalized))

	// language=rb
	code = `
class Customer
   @@no_of_customers = 0
   def initialize(id, name, addr)
      @cust_id = id
      @cust_name = name
      @cust_addr = addr
   end
end
`
	counters, err = getCountersForCode(code, "ruby")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(8), counters.LinesOfCode)
	r.Equal(float64(4), counters.Keywords)
	r.Equal(float64(27), counters.Indentations)
	r.Equal(float64(9), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(6), math.Round(counters.IndentationsDiff))
	r.Equal(float64(2), math.Round(counters.IndentationsDiffNormalized))
}

func TestCountersForRubyFullSample(t *testing.T) {
	r := assert.New(t)

	counters, err := getCountersForCode(test_resources.RubyCode, "ruby")
	r.Nil(err)
	r.NotNil(counters)

	r.Equal(float64(465), counters.Lines)
	r.Equal(float64(242), counters.LinesOfCode)
	r.Equal(float64(149), counters.Keywords)
	r.Equal(float64(1922), counters.Indentations)
	r.Equal(float64(961), math.Round(counters.IndentationsNormalized))
	r.Equal(float64(145), math.Round(counters.IndentationsDiff))
	r.Equal(float64(73), math.Round(counters.IndentationsDiffNormalized))
	r.Equal(float64(62), math.Round(counters.KeywordsComplexity*100))
	r.Equal(float64(397), math.Round(counters.IndentationsComplexity*100))
	r.Equal(float64(30), math.Round(counters.IndentationsDiffComplexity*100))
}
