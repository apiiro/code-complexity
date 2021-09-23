package calculate

import (
	"code-complexity/test_resources"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

// test_resources patterns, encoding

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
let x = 3;
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
			x = 8;
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
     x = 8;
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
	x = 4;
}
else {
	if (x > 7) {
			x = 8;
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
  x = 4;
}
else {
  if (x > 7) {
     x = 8;
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
	r.Equal(float64(63), counters.Keywords)
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
	x = 4;
}
else {
	if (x > 7) {
			x = 8;
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
  x = 4;
}
else {
  if (x > 7) {
     x = 8;
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

// c
// cpp
// objc
// swift
// rb
// go
// rs
// tf
