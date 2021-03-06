package options

type Config struct {
	IncludePatterns []string `json:"include_patterns"`
	ExcludePatterns []string `json:"exclude_patterns"`
}

var defaultConfig = &Config{
	IncludePatterns: []string{},
	ExcludePatterns: []string{
		"**/bin",
		"**/obj",
		"**/venv",
		"**/node_modules",
		"**/.idea",
		"**/.git",
		"**/site-packages",
		"**/vendor",
		"**/test_resources",
		"**/tests/**",
		"**/testing/**",
		"**/resources/**",
		"**/testdata/**",
		"**/simulation/**",
		"**/simulator/**",
		"**/automation/**",
		"**/*test_resources.*",
		"**/*tests.*",
		"**/*spec.*",
	},
}
