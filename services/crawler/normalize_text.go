package crawler

import (
	"strings"
	"unicode"
)

func BuildStopWordSet(words []string) map[string]struct{} {
	set := make(map[string]struct{}, len(words))
	for _, w := range words {
		set[w] = struct{}{}
	}
	return set
}

func TokenizeAndRemoveStopWords(text string, stopWords map[string]struct{}) []string {
	text = strings.ToLower(text)

	tokens := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	result := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if _, isStopWord := stopWords[token]; !isStopWord {
			result = append(result, token)
		}
	}

	return result
}
