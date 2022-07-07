package helpers

import (
	"fmt"
	"time"
)

// RoundTime возвращает форматированную строку
// разницу в годах и месяцах между двумя датами

func RoundTime(from, to time.Time) string {
	if from.Location() != to.Location() {
		to = to.In(from.Location())
	}
	if from.After(to) {
		from, to = to, from
	}
	y1, M1, d1 := from.Date()
	y2, M2, d2 := to.Date()

	year := int(y2 - y1)
	month := int(M2 - M1)
	day := int(d2 - d1)

	if day < 0 {
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	var yearSuf string

	sh := year % 10

	if (year > 0 && year < 5) ||
		(year > 20 && sh > 0 && sh < 5) {
		yearSuf = "г"
	} else {
		yearSuf = "л"
	}
	if month == 0 {
		return fmt.Sprintf("%d %s", year, yearSuf)
	}
	if year == 0 {
		return fmt.Sprintf("%d мес", month)
	}
	return fmt.Sprintf("%d %s %d мес", year, yearSuf, month)
}
