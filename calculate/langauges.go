package calculate

type Language = string

var languageToExtensions = map[Language][]string{
	"java":       {"java"},
	"csharp":     {"cs", "cshtml"},
	"node":       {"js", "jsx", "ts", "tsx"},
	"python":     {"py", "py3", "py2"},
	"kotlin":     {"kt", "kts", "ktm"},
	"c":          {"c", "h"},
	"cpp":        {"cpp", "cxx", "cc", "hpp", "hh", "txx", "tpp"},
	"objectivec": {"m", "mm"},
	"swift":      {"swift"},
	"ruby":       {"rb"},
	"go":         {"go"},
	"rust":       {"rs"},
	"scala":      {"scala", "sc"},
	"php":        {"php", "phtml", "php3", "php4", "php5", "php7", "phps", "pht", "phar"},
}

var extensionToLanguage = make(map[Language]string)

func init() {
	for language, extensions := range languageToExtensions {
		for _, extension := range extensions {
			extensionToLanguage[extension] = language
		}
	}
}

func tryGetLanguage(ext string) (Language, bool) {
	if len(ext) == 0 {
		return "", false
	}
	ext = ext[1:]
	language, found := extensionToLanguage[ext]
	return language, found
}
