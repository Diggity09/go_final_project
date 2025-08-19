package api

import (
	"testing"
	"time"
)

func TestNextDate(t *testing.T) {
	// Тестовая дата: 26.01.2024
	now, _ := time.Parse("20060102", "20240126")

	tests := []struct {
		name     string
		dstart   string
		repeat   string
		expected string
		hasError bool
	}{
		// Тесты для ежегодного повторения
		{
			name:     "Yearly repeat - leap year to non-leap year",
			dstart:   "20240229",
			repeat:   "y",
			expected: "20250301",
			hasError: false,
		},
		{
			name:     "Yearly repeat - normal date",
			dstart:   "20240115",
			repeat:   "y",
			expected: "20250115",
			hasError: false,
		},

		// Тесты для повторения по дням
		{
			name:     "Daily repeat - 7 days",
			dstart:   "20240113",
			repeat:   "d 7",
			expected: "20240127",
			hasError: false,
		},
		{
			name:     "Daily repeat - 1 day",
			dstart:   "20240126",
			repeat:   "d 1",
			expected: "20240127",
			hasError: false,
		},

		// Тесты для повторения по месяцам
		{
			name:     "Monthly repeat - specific days",
			dstart:   "20240116",
			repeat:   "m 16,5",
			expected: "20240205",
			hasError: false,
		},
		{
			name:     "Monthly repeat - last day and 18th",
			dstart:   "20240201",
			repeat:   "m -1,18",
			expected: "20240218",
			hasError: false,
		},

		// Тесты ошибок
		{
			name:     "Empty repeat",
			dstart:   "20240126",
			repeat:   "",
			expected: "",
			hasError: true,
		},
		{
			name:     "Invalid date format",
			dstart:   "20240232",
			repeat:   "d 1",
			expected: "",
			hasError: true,
		},
		{
			name:     "Invalid repeat format",
			dstart:   "20240126",
			repeat:   "l 34",
			expected: "",
			hasError: true,
		},
		{
			name:     "Daily without interval",
			dstart:   "20240126",
			repeat:   "d",
			expected: "",
			hasError: true,
		},
		{
			name:     "Daily with too large interval",
			dstart:   "20240126",
			repeat:   "d 405",
			expected: "20250405", // Не ошибка, просто большой интервал
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NextDate(now, tt.dstart, tt.repeat)

			if tt.hasError {
				if err == nil {
					t.Errorf("Ожидалась ошибка, но её не было")
				}
			} else {
				if err != nil {
					t.Errorf("Неожиданная ошибка: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Ожидался результат %s, получен %s", tt.expected, result)
				}
			}
		})
	}
}

func TestWeeklyRepeat(t *testing.T) {
	// Понедельник, 22.01.2024
	now, _ := time.Parse("20060102", "20240122")

	tests := []struct {
		name     string
		dstart   string
		repeat   string
		expected string
		hasError bool
	}{
		{
			name:     "Weekly repeat - Sunday",
			dstart:   "20240121",
			repeat:   "w 7",
			expected: "20240128",
			hasError: false,
		},
		{
			name:     "Weekly repeat - Monday, Thursday, Friday",
			dstart:   "20240121",
			repeat:   "w 1,4,5",
			expected: "20240125",
			hasError: false,
		},
		{
			name:     "Weekly repeat - invalid day",
			dstart:   "20240121",
			repeat:   "w 8",
			expected: "",
			hasError: true,
		},
		{
			name:     "Weekly repeat - no days specified",
			dstart:   "20240121",
			repeat:   "w",
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NextDate(now, tt.dstart, tt.repeat)

			if tt.hasError {
				if err == nil {
					t.Errorf("Ожидалась ошибка, но её не было")
				}
			} else {
				if err != nil {
					t.Errorf("Неожиданная ошибка: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Ожидался результат %s, получен %s", tt.expected, result)
				}
			}
		})
	}
}

func TestCheckDate(t *testing.T) {
	tests := []struct {
		name     string
		task     Task
		expected string
		hasError bool
	}{
		{
			name:     "Empty date should use today",
			task:     Task{Title: "Test", Date: ""},
			expected: time.Now().Format("20060102"),
			hasError: false,
		},
		{
			name:     "Valid future date",
			task:     Task{Title: "Test", Date: "20301225"},
			expected: "20301225",
			hasError: false,
		},
		{
			name:     "Past date without repeat should use today",
			task:     Task{Title: "Test", Date: "20200101"},
			expected: time.Now().Format("20060102"),
			hasError: false,
		},
		{
			name:     "Invalid date format",
			task:     Task{Title: "Test", Date: "2024-12-25"},
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkDate(&tt.task)

			if tt.hasError {
				if err == nil {
					t.Errorf("Ожидалась ошибка, но её не было")
				}
			} else {
				if err != nil {
					t.Errorf("Неожиданная ошибка: %v", err)
				}
				if tt.task.Date != tt.expected {
					t.Errorf("Ожидалась дата %s, получена %s", tt.expected, tt.task.Date)
				}
			}
		})
	}
}
