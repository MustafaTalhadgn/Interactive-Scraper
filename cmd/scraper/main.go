package main

import (
	"errors"
	"io"
	"net/http"
	"time"
)

type Record struct {
	SourceName string
	SourceUrl  string
	Content    string
	CreatedAt  time.Time
}

func FetchRecord(sourceName string, sourceUrl string) (Record, error) {

	record := Record{
		SourceName: sourceName,
		SourceUrl:  sourceUrl,
		Content:    "",
		CreatedAt:  time.Now(),
	}

	resp, err := http.Get(record.SourceUrl)
	if err != nil {
		return Record{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New("Sayfa 200 vermedi")
		return Record{}, err

	}
	
	rawContent, err := io.ReadAll(resp.Body)

	if err != nil {
		return Record{}, err
	}

	content := string(rawContent)

	record.Content = content

	return record, nil

}
