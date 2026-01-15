package parser

import (
	"strings"
	"time"
)

var dateFormats = []string{
	// ISO 8601
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02T15:04:05Z",
	"2006-01-02 15:04:05",
	"2006-01-02",

	// US formats
	"01/02/2006 15:04:05",
	"01/02/2006",
	"Jan 02, 2006",
	"January 02, 2006",

	// EU formats
	"02/01/2006 15:04:05",
	"02/01/2006",
	"02.01.2006",
	"02-01-2006",

	// RFC formats
	time.RFC3339,
	time.RFC1123,
	time.RFC822,
}

func ParseDate(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return time.Time{}, nil
	}

	for _, format := range dateFormats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, nil
}

func ExtractDateFromText(text string) (time.Time, error) {
	text = strings.TrimPrefix(text, "Posted on")
	text = strings.TrimPrefix(text, "Date:")
	text = strings.TrimPrefix(text, "Published:")
	text = strings.TrimSpace(text)
	return ParseDate(text)
}
