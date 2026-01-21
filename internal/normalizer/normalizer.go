package normalizer

import "log/slog"

type Normalizer struct {
	unicode   *UnicodeNormalizer
	stopwords *StopwordManager
	textUtils *TextUtils
	config    *Config
	logger    *slog.Logger
}

type Config struct {
	Lowercase          bool
	NormalizeUnicode   bool
	RemoveAccents      bool
	CollapseWhitespace bool
	RemoveStopwords    bool
	RemovePunctuation  bool
	MaxLength          int
}

func DefaultConfig() *Config {
	return &Config{
		Lowercase:          true,
		NormalizeUnicode:   true,
		RemoveAccents:      false,
		CollapseWhitespace: true,
		RemoveStopwords:    false,
		RemovePunctuation:  false,
		MaxLength:          1000,
	}
}

func NewNormalizer(config *Config, logger *slog.Logger) *Normalizer {
	if config == nil {
		config = DefaultConfig()
	}

	return &Normalizer{
		unicode:   NewUnicodeNormalizer(),
		stopwords: NewStopwordManager(config.RemoveStopwords),
		textUtils: NewTextUtils(),
		config:    config,
		logger:    logger,
	}
}

type NormalizedText struct {
	Text      string
	WordCount int
	Original  string
}

func (n *Normalizer) Normalize(text string) *NormalizedText {
	original := text

	if n.config.NormalizeUnicode {
		text = n.unicode.Normalize(text)
	}

	if n.config.RemoveAccents {
		text = n.unicode.RemoveAccents(text)
	}

	if n.config.Lowercase {
		text = n.textUtils.ToLower(text)
	}

	if n.config.CollapseWhitespace {
		text = n.textUtils.RemoveExtraWhitespace(text)
	}

	if n.config.RemovePunctuation {
		text = n.textUtils.RemovePunctuation(text)
	}

	if n.config.RemoveStopwords {
		text = n.stopwords.RemoveStopwords(text)
	}

	if n.config.MaxLength > 0 && len(text) > n.config.MaxLength {
		text = n.textUtils.TrimToWordBoundary(text, n.config.MaxLength)
		n.logger.Debug("kısaltma uygulandı",
			slog.Int("Orijinal uzunluk", len(original)),
			slog.Int("Kısaltılmış uzunluk", len(text)),
		)
	}

	return &NormalizedText{
		Text:      text,
		WordCount: n.textUtils.WordCount(text),
		Original:  original,
	}
}

func (n *Normalizer) NormalizeForScoring(text string) string {

	scoringConfig := &Config{
		Lowercase:          true,
		NormalizeUnicode:   true,
		RemoveAccents:      true,
		CollapseWhitespace: true,
		RemoveStopwords:    true,
		RemovePunctuation:  true,
		MaxLength:          n.config.MaxLength,
	}

	scoringNormalizer := NewNormalizer(scoringConfig, n.logger)

	result := scoringNormalizer.Normalize(text)
	return result.Text
}

func (n *Normalizer) NormalizeForDisplay(text string) string {
	displayConfig := &Config{
		Lowercase:          false,
		NormalizeUnicode:   true,
		RemoveAccents:      false,
		CollapseWhitespace: true,
		RemoveStopwords:    false,
		RemovePunctuation:  false,
		MaxLength:          n.config.MaxLength,
	}

	displayNormalizer := NewNormalizer(displayConfig, n.logger)

	result := displayNormalizer.Normalize(text)
	return result.Text
}

func (n *Normalizer) NormalizeKeywords(text string) []string {

	normalized := n.NormalizeForScoring(text)

	words := n.textUtils.SplitIntoWords(normalized)

	var keywords []string
	for _, word := range words {
		if len(word) >= 3 {
			keywords = append(keywords, word)
		}
	}

	return keywords
}
