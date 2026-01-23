package parser

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type Parser struct {
	config    *ParserConfig
	strategy  *SelectorStrategy
	converter *md.Converter
}

func NewParser(config *ParserConfig) *Parser {
	if config == nil {
		config = DefaultParserConfig()
	}

	return &Parser{
		config:    config,
		strategy:  DefaultStrategy(),
		converter: md.NewConverter("", true, nil),
	}
}

func (p *Parser) Parse(htmlBytes []byte, url string) (*ParsedData, error) {

	doc, err := html.Parse(bytes.NewReader(htmlBytes))
	if err != nil {
		return nil, fmt.Errorf("html parse hatalı: %w", err)
	}
	gqDoc := goquery.NewDocumentFromNode(doc)

	gqDoc.Find("script, style, link, meta, noscript, iframe, svg").Remove()

	data := &ParsedData{
		URL: url,
	}

	data.Title = p.extractTitle(gqDoc)

	data.Content = p.extractContent(gqDoc)

	data.Author = p.extractAuthor(gqDoc)
	data.Date = p.extractDate(gqDoc)

	data.WordCount = WordCount(data.Content)

	if p.config.ExtractLinks {
		data.Links = p.extractLinks(gqDoc)
	}

	if p.config.ExtractImages {
		data.ImageURLs = p.extractImages(gqDoc)
	}

	if p.config.SaveRawHTML {
		data.RawHTML = string(htmlBytes)
	}

	if len(data.Content) > p.config.MaxContentLen {
		data.Content = TruncateText(data.Content, p.config.MaxContentLen)
	}

	return data, nil
}

func (p *Parser) extractContent(doc *goquery.Document) string {
	var contentParts []string
	found := false

	for _, selector := range p.strategy.ContentSelectors {
		selection := doc.Find(selector)

		if selection.Length() > 0 {
			found = true

			selection.Each(func(i int, s *goquery.Selection) {

				markdown := p.converter.Convert(s)

				if strings.TrimSpace(markdown) != "" {
					contentParts = append(contentParts, markdown)
				} else {

					text := CleanText(s.Text())
					if text != "" {
						contentParts = append(contentParts, text)
					}
				}
			})

			break
		}
	}

	if !found {
		bodySel := doc.Find("body")
		md := p.converter.Convert(bodySel)
		if md != "" {
			return md
		}
		return CleanText(bodySel.Text())
	}

	return strings.Join(contentParts, "\n\n---\n\n")
}

func (p *Parser) extractTitle(doc *goquery.Document) string {
	for _, selector := range p.strategy.TitleSelectors {
		if text := doc.Find(selector).First().Text(); text != "" {
			return CleanText(text)
		}
	}
	return "Untitled"
}

func (p *Parser) extractAuthor(doc *goquery.Document) string {
	for _, selector := range p.strategy.AuthorSelectors {
		if text := doc.Find(selector).First().Text(); text != "" {
			return CleanText(text)
		}
	}
	return "Unknown"
}

func (p *Parser) extractDate(doc *goquery.Document) time.Time {
	for _, selector := range p.strategy.DateSelectors {
		elem := doc.Find(selector).First()

		if datetime, exists := elem.Attr("datetime"); exists {
			if date, err := ParseDate(datetime); err == nil && !date.IsZero() {
				return date
			}
		}

		if text := elem.Text(); text != "" {
			if date, err := ExtractDateFromText(text); err == nil && !date.IsZero() {
				return date
			}
		}
	}
	return time.Time{}
}

func (p *Parser) extractLinks(doc *goquery.Document) []string {
	var links []string
	seen := make(map[string]bool)

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			href = strings.TrimSpace(href)
			if href != "" && href != "#" && !seen[href] {
				links = append(links, href)
				seen[href] = true
			}
		}
	})

	return links
}

func (p *Parser) extractImages(doc *goquery.Document) []string {
	var images []string
	seen := make(map[string]bool)

	doc.Find("img[src]").Each(func(i int, s *goquery.Selection) {
		if src, exists := s.Attr("src"); exists {
			src = strings.TrimSpace(src)
			if src != "" && !seen[src] {
				images = append(images, src)
				seen[src] = true
			}
		}
	})

	return images
}
