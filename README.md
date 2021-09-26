# code-complexity

Tool to estimate code complexity with the intention of providing baseline metrics for full static code analysis.

The algorithm is inspired by [terryyin/lizard](https://github.com/terryyin/lizard) and [thoughtbot/complexity](https://github.com/thoughtbot/complexity).

![icon](code-complexity.png)

```
NAME:
   complexity - 1.0.0 - Estimate source code complexity

USAGE:
   complexity        [optional flags]

OPTIONS:
   --dir value, -d value      path to directory containing directory path, defaults to current directory
   --config value, -c value   include/exclude patterns config file (default: "unset")
   --out value, -o value      output file, or empty to print to stdout
   --include value, -i value  patterns of file paths to include, comma delimited, may contain any glob pattern
   --exclude value, -e value  patterns of file paths to exclude, comma delimited, may contain any glob pattern
   --verbose, --vv            verbose logging (default: false)
   --max-size value           maximal file size, in MB (default: 6)
   --help, -h                 show help (default: false)
   --version, -v              print the version (default: false)
```

## Output

Per supported [programming language](#languages), the tool will plot the number of source files, and following metrics, in both `total` and `average` sections:

* Lines of Code (`lines_of_code`) - Number of lines that don't contain whitespace or comments.
* Keywords Complexity (`keywords_complexity`) - Number of keywords per line of code. Keyword is a rough estimation of control statements that are defined per language, see [languageToKeywords](calculate/keywords.go).
* Indentations Complexity (`indentations_complexity`) - Normalized number of indentations per line of code.
* Indentations Diff Complexity (`indentations_diff_complexity`) - Normalized number of positive indentations diff per line of code.

Output example:

```json
{
  "counters_by_language": {
    "go": {
      "number_of_files": 9,
      "total": {
        "lines_of_code": 2374,
        "keywords_complexity": 2.039620976028679,
        "indentations_complexity": 11.930908025104817,
        "indentations_diff_complexity": 1.9046008903365483
      },
      "average": {
        "lines_of_code": 263.77777777777777,
        "keywords_complexity": 0.22662455289207545,
        "indentations_complexity": 1.3256564472338686,
        "indentations_diff_complexity": 0.21162232114850538
      }
    }
  }
}
```

## Examples

```bash
complexity # will run on current directory with default configs
complexity -d "path/to/src" -o "output.json"
complexity -d "proj/src" -o "proj/output.json" -c "proj/.config.json"
complexity -d "proj/src" -o "proj/output.json" -c "proj/.config.json" -i 'src/**,**.js,**.ts' -e 'test/**'
```

## Install

```bash
curl -s https://raw.githubusercontent.com/apiiro/code-complexity/main/install.sh | sudo bash
# or for a specific version:
curl -s https://raw.githubusercontent.com/apiiro/code-complexity/main/install.sh | sudo bash -s 1.4
```

If that doesn't work, try:
```bash
curl -s https://raw.githubusercontent.com/apiiro/code-complexity/main/install.sh -o install.sh
sudo bash install.sh
```

then run:

```bash
complexity -h
```

## Test and Build

```bash
# run tests:
make test
# run benchmark
make benchmark
# build binaries and run whole ci flow
make
```

### Languages

Following languages are currently supported:

* Java
* C#
* Node (Javascript/Typescript backend)
* Python
* Kotlin
* C
* C++
* Objective-C
* Swift
* Ruby
* Go
* Rust
* Scala

### Credits

<div>Icons made by <a href="https://www.freepik.com" title="Freepik">Freepik</a> from <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a></div>
