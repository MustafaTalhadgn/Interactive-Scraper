package parser

import "time"

type ParsedData struct {
	Title   string
	Content string
	Author  string
	Date    time.Time

	URL       string
	Language  string
	WordCount int
	Links     []string
	ImageURLs []string

	RawHTML string
}

type ParserConfig struct {
	ExtractLinks  bool
	ExtractImages bool
	SaveRawHTML   bool
	MaxContentLen int
}

func DefaultParserConfig() *ParserConfig {
	return &ParserConfig{
		ExtractLinks:  true,
		ExtractImages: false,
		SaveRawHTML:   false,
		MaxContentLen: 50000,
	}
}
