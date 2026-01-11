package models

import (
	"fmt"
	"strings"
	"time"
)

type TravelRequest struct {
	Destination      string    `json:"destination"`
	DepartureCity    string    `json:"departure_city"`
	TravelDates      string    `json:"travel_dates"`
	Duration         string    `json:"duration"`
	Travelers        string    `json:"travelers"`
	ChildAge         string    `json:"child_age"`
	Budget           string    `json:"budget"`
	VacationType     string    `json:"vacation_type"`
	HotelLevel       string    `json:"hotel_level"`
	MealPlan         string    `json:"meal_plan"`
	ImportantFactors string    `json:"important_factors"`
	CreatedAt        time.Time `json:"created_at"`
}

type UserInfo struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

func escapeHTML(text string) string {
	replacements := []struct {
		old string
		new string
	}{
		{"&", "&amp;"},
		{"<", "&lt;"},
		{">", "&gt;"},
		{"\"", "&quot;"},
		{"'", "&#39;"},
	}

	result := text
	for _, r := range replacements {
		result = strings.ReplaceAll(result, r.old, r.new)
	}

	return result
}

func (tr *TravelRequest) ToFormattedString(userInfo UserInfo) string {
	var builder strings.Builder

	builder.WriteString("<b>🌴 НОВАЯ ЗАЯВКА НА ПОДБОР ТУРА</b>\n\n")

	// Информация о клиенте (для менеджера)
	builder.WriteString("<b>👤 Клиент:</b> ")
	if userInfo.FirstName != "" || userInfo.LastName != "" {
		builder.WriteString(escapeHTML(userInfo.FirstName + " " + userInfo.LastName))
	}
	if userInfo.Username != "" {
		builder.WriteString(fmt.Sprintf("\n<b>📱 @:</b> %s", escapeHTML(userInfo.Username)))
	}
	builder.WriteString(fmt.Sprintf("\n<b>🆔 ID:</b> %d\n", userInfo.ID))

	builder.WriteString("\n<b>═══════════════════════════════════</b>\n\n")

	// Данные заявки
	writeFieldHTML(&builder, "1️⃣ Куда планируете поездку?", tr.Destination)
	writeFieldHTML(&builder, "2️⃣ Город вылета", tr.DepartureCity)
	writeFieldHTML(&builder, "3️⃣ Даты поездки", tr.TravelDates)
	writeFieldHTML(&builder, "4️⃣ Длительность отдыха", tr.Duration)
	writeFieldHTML(&builder, "5️⃣ Количество туристов", tr.Travelers)

	if tr.ChildAge != "" && tr.ChildAge != "Нет детей" {
		writeFieldHTML(&builder, "   Возраст ребенка", tr.ChildAge)
	}

	writeFieldHTML(&builder, "6️⃣ Бюджет на всех", tr.Budget)
	writeFieldHTML(&builder, "7️⃣ Тип отдыха", tr.VacationType)
	writeFieldHTML(&builder, "8️⃣ Уровень отеля", tr.HotelLevel)
	writeFieldHTML(&builder, "9️⃣ Тип питания", tr.MealPlan)
	writeFieldHTML(&builder, "🔟 Принципиально важно", tr.ImportantFactors)

	builder.WriteString("\n<b>═══════════════════════════════════</b>\n")
	builder.WriteString(fmt.Sprintf("<b>📅 Заявка создана:</b> %s\n",
		tr.CreatedAt.Format("02.01.2006 в 15:04")))

	return builder.String()
}

func (tr *TravelRequest) ToClientPreview() string {
	var builder strings.Builder

	builder.WriteString("<b>🌴 ВАША ЗАЯВКА НА ПОДБОР ТУРА</b>\n\n")

	builder.WriteString("<b>═══════════════════════════════════</b>\n\n")

	writeFieldHTML(&builder, "1️⃣ Куда планируете поездку?", tr.Destination)
	writeFieldHTML(&builder, "2️⃣ Город вылета", tr.DepartureCity)
	writeFieldHTML(&builder, "3️⃣ Даты поездки", tr.TravelDates)
	writeFieldHTML(&builder, "4️⃣ Длительность отдыха", tr.Duration)
	writeFieldHTML(&builder, "5️⃣ Количество туристов", tr.Travelers)

	if tr.ChildAge != "" && tr.ChildAge != "Нет детей" {
		writeFieldHTML(&builder, "   Возраст ребенка", tr.ChildAge)
	}

	writeFieldHTML(&builder, "6️⃣ Бюджет на всех", tr.Budget)
	writeFieldHTML(&builder, "7️⃣ Тип отдыха", tr.VacationType)
	writeFieldHTML(&builder, "8️⃣ Уровень отеля", tr.HotelLevel)
	writeFieldHTML(&builder, "9️⃣ Тип питания", tr.MealPlan)
	writeFieldHTML(&builder, "🔟 Принципиально важно", tr.ImportantFactors)

	builder.WriteString("\n<b>═══════════════════════════════════</b>\n")
	builder.WriteString(fmt.Sprintf("<b>📅 Заявка создана:</b> %s\n",
		tr.CreatedAt.Format("02.01.2006 в 15:04")))

	return builder.String()
}

func writeFieldHTML(builder *strings.Builder, name, value string) {
	if value == "" {
		value = "Не указано"
	}
	builder.WriteString(fmt.Sprintf("<b>%s</b>\n%s\n\n",
		escapeHTML(name),
		escapeHTML(value)))
}
