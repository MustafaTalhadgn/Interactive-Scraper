package normalizer

import (
	"unicode"

	"golang.org/x/text/unicode/norm"
)

type UnicodeNormalizer struct {
	form norm.Form
}

func NewUnicodeNormalizer() *UnicodeNormalizer {
	return &UnicodeNormalizer{
		form: norm.NFC,
	}
}

func (u *UnicodeNormalizer) Normalize(text string) string {
	return u.form.String(text)
}

func (u *UnicodeNormalizer) RemoveAccents(text string) string {

	decomposed := norm.NFD.String(text)
	var result []rune
	for _, r := range decomposed {
		if !unicode.Is(unicode.Mn, r) { // Mn = Nonspacing Mark (accents)
			result = append(result, r)
		}
	}
	return norm.NFC.String(string(result))
}

func (u *UnicodeNormalizer) NormalizeWhitespace(text string) string {
	var result []rune
	for _, r := range text {
		if unicode.IsSpace(r) {
			result = append(result, ' ') // Convert all whitespace to space
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}
