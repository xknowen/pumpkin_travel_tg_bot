package utils

import (
	"regexp"
	"strings"
)

func ValidateNotEmpty(text string) bool {
	return strings.TrimSpace(text) != ""
}

func ValidateBudget(text string) bool {
	// Проверяем, содержит ли текст цифры
	hasDigits := regexp.MustCompile(`\d`).MatchString(text)
	return hasDigits && ValidateNotEmpty(text)
}

func ValidateCountries(countriesStr string) []string {
	// Разделяем строку на страны, убираем пустые элементы
	countries := strings.Split(countriesStr, ",")
	var cleanCountries []string

	for _, country := range countries {
		trimmed := strings.TrimSpace(country)
		if trimmed != "" {
			cleanCountries = append(cleanCountries, trimmed)
		}
	}

	return cleanCountries
}
