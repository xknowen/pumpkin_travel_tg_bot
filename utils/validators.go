package utils

import (
	"regexp"
	"strings"
)

func ValidateNotEmpty(text string) bool {
	return strings.TrimSpace(text) != ""
}

func ValidateBudget(text string) bool {
	hasDigits := regexp.MustCompile(`\d`).MatchString(text)
	hasFlexible := strings.Contains(strings.ToLower(text), "без строг") ||
		strings.Contains(strings.ToLower(text), "не имеет") ||
		strings.Contains(strings.ToLower(text), "не принципиал")

	return (hasDigits || hasFlexible) && ValidateNotEmpty(text)
}

func ValidateCountries(countriesStr string) []string {
	separators := []string{"/", ",", "и", "или"}

	result := countriesStr
	for _, sep := range separators {
		result = strings.ReplaceAll(result, sep, ",")
	}

	countries := strings.Split(result, ",")
	var cleanCountries []string

	for _, country := range countries {
		trimmed := strings.TrimSpace(country)
		if trimmed != "" && trimmed != "пока не определились" {
			cleanCountries = append(cleanCountries, trimmed)
		}
	}

	return cleanCountries
}
