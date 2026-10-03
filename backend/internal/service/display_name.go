package service

import (
	"strings"
	"unicode"
)

const (
	maxNameLen   = 12
	fallbackName = "Jogador"
)

// CleanName devolve o primeiro candidato que sobra depois de tirar emojis e
// caracteres especiais. Nome longo demais vira só o primeiro nome.
func CleanName(candidates ...string) string {
	for _, c := range candidates {
		if name := cleanName(c); name != "" {
			return name
		}
	}
	return fallbackName
}

func cleanName(raw string) string {
	words := strings.FieldsFunc(raw, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	if len(words) == 0 {
		return ""
	}
	name := strings.Join(words, " ")
	if len([]rune(name)) <= maxNameLen {
		return name
	}
	first := []rune(words[0])
	if len(first) > maxNameLen {
		first = first[:maxNameLen]
	}
	return string(first)
}
