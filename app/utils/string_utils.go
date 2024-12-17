package utils

import (
	"math/rand"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func SnakeCase(content string) string {
	return strings.ReplaceAll(strings.ToLower(content), " ", "_")
}

func CamelCase(content string) string {
	words := []string{}
	l := 0
	for r := 0; r < len(content); r++ {
		if content[r] == '_' || content[r] == ' ' {
			if l < r {
				words = append(words, content[l:r])
			}
			l = r
		}
	}
	ans := ""
	c := cases.Title(language.English)
	for _, word := range words {
		ans += c.String(word)
	}
	return ans
}

func SnakeToCamelCase(content string) string {
	return CamelCase(strings.ReplaceAll(content, "_", " "))
}

func CamelToSnakeCase(content string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(content, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func GenerateRandomString(n int, useLower bool, useUpper bool, useSpecialChar bool) string {
	lowerLetters := "abcdefghijklmnopqrstuvwxyz"
	upperLetters := strings.ToUpper(lowerLetters)
	specialChars := "!@#$%&"
	letters := ""
	if useLower {
		letters += lowerLetters
	}
	if useUpper {
		letters += upperLetters
	}
	if useSpecialChar {
		letters += specialChars
	}
	if letters == "" {
		letters = lowerLetters
	}

	b := make([]rune, n)
	for i := range b {
		b[i] = rune(letters[rand.Intn(len(letters))])
	}
	return string(b)
}
