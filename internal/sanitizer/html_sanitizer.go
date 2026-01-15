package sanitizer

import (
	"github.com/microcosm-cc/bluemonday"
)

type HTMLSanitizer struct {
	strictPolicy  *bluemonday.Policy
	relaxedPolicy *bluemonday.Policy
}

func NewHTMLSanitizer() *HTMLSanitizer {
	return &HTMLSanitizer{
		strictPolicy:  createStrictPolicy(),
		relaxedPolicy: createRelaxedPolicy(),
	}
}

func createStrictPolicy() *bluemonday.Policy {
	p := bluemonday.StrictPolicy()

	return p
}

func createRelaxedPolicy() *bluemonday.Policy {
	p := bluemonday.NewPolicy()

	p.AllowElements("p", "br", "strong", "em", "u", "h1", "h2", "h3")

	p.AllowAttrs("class").OnElements("p", "div")

	p.SkipElementsContent("script", "style", "iframe", "object", "embed")

	return p
}

func (s *HTMLSanitizer) SanitizeToPlainText(html string) string {
	return s.strictPolicy.Sanitize(html)
}

func (s *HTMLSanitizer) SanitizeToSafeHTML(html string) string {
	return s.relaxedPolicy.Sanitize(html)
}

func (s *HTMLSanitizer) RemoveScriptTags(html string) string {
	p := bluemonday.NewPolicy()
	p.SkipElementsContent("script")
	return p.Sanitize(html)
}
