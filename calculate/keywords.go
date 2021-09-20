package calculate

import "strings"

func countKeywords(line string, language Language) float64 {
	keywords, found := languageToKeywords[language]
	if !found {
		return 0
	}
	keywordsCount := float64(0)
	tokens := strings.Fields(line)
	tokensSet := make(map[string]bool, len(tokens))
	for _, token := range strings.Fields(line) {
		if strings.HasPrefix(token, "@") {
			if _, found := languagesWithAtSignPrefix[language]; found {
				keywordsCount++
			}
		}
		token = strings.TrimRight(strings.TrimRight(token, ";"), "{")
		tokensSet[token] = true
	}
	for _, keyword := range keywords {
		if _, found := tokensSet[keyword]; found {
			keywordsCount++
		}
	}
	return keywordsCount
}
