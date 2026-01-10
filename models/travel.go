package models

import (
	"fmt"
	"strings"
	"time"
)

type TravelPreferences struct {
	DestinationType     string    `json:"destination_type"`
	Countries           []string  `json:"countries"`
	BudgetPerPerson     string    `json:"budget_per_person"`
	TravelersCount      string    `json:"travelers_count"`
	TravelDates         string    `json:"travel_dates"`
	Duration            string    `json:"duration"`
	AccommodationType   string    `json:"accommodation_type"`
	SpecialRequirements string    `json:"special_requirements"`
	CreatedAt           time.Time `json:"created_at"`
}

type UserInfo struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// Функция для экранирования Markdown символов
func escapeMarkdown(text string) string {
	// Экранируем специальные символы Markdown
	replacements := []struct {
		old string
		new string
	}{
		{"_", "\\_"},
		{"*", "\\*"},
		{"[", "\\["},
		{"]", "\\]"},
		{"(", "\\("},
		{")", "\\)"},
		{"~", "\\~"},
		{"`", "\\`"},
		{">", "\\>"},
		{"#", "\\#"},
		{"+", "\\+"},
		{"-", "\\-"},
		{"=", "\\="},
		{"|", "\\|"},
		{"{", "\\{"},
		{"}", "\\}"},
		{".", "\\."},
		{"!", "\\!"},
	}

	result := text
	for _, r := range replacements {
		result = strings.ReplaceAll(result, r.old, r.new)
	}

	return result
}

func (tp *TravelPreferences) ToFormattedString(userInfo UserInfo) string {
	var builder strings.Builder

	builder.WriteString("🗺 *Новая заявка от клиента\\!*\n")

	// Экранируем информацию о пользователе
	firstName := escapeMarkdown(userInfo.FirstName)
	lastName := escapeMarkdown(userInfo.LastName)
	username := ""
	if userInfo.Username != "" {
		username = escapeMarkdown(userInfo.Username)
	}

	builder.WriteString(fmt.Sprintf("👤 *Клиент:* %s %s\n", firstName, lastName))

	if username != "" {
		builder.WriteString(fmt.Sprintf("📱 *Username:* @%s\n", username))
	}

	builder.WriteString(fmt.Sprintf("🆔 *ID:* %d\n", userInfo.ID))
	builder.WriteString("*===============================*\n")

	// Экранируем все поля
	writeField(&builder, "Тип отдыха", escapeMarkdown(tp.DestinationType))
	writeField(&builder, "Страны/Направления", escapeMarkdown(strings.Join(tp.Countries, ", ")))
	writeField(&builder, "Бюджет на человека", escapeMarkdown(tp.BudgetPerPerson))
	writeField(&builder, "Количество путешественников", escapeMarkdown(tp.TravelersCount))
	writeField(&builder, "Даты/Период поездки", escapeMarkdown(tp.TravelDates))
	writeField(&builder, "Продолжительность", escapeMarkdown(tp.Duration))
	writeField(&builder, "Тип проживания", escapeMarkdown(tp.AccommodationType))

	specialReqs := tp.SpecialRequirements
	if specialReqs == "" {
		specialReqs = "Нет"
	}
	writeField(&builder, "Особые пожелания", escapeMarkdown(specialReqs))

	builder.WriteString("*===============================*\n")
	builder.WriteString(fmt.Sprintf("*Время подачи заявки:* %s\n",
		tp.CreatedAt.Format("02\\.01\\.2006 15:04"))) // Экранируем точки в дате

	return builder.String()
}

func writeField(builder *strings.Builder, name, value string) {
	if value == "" || value == "Нет" || value == "Нет особых" {
		value = "Не указано"
	}
	builder.WriteString(fmt.Sprintf("*%s:* %s\n", name, value))
}
