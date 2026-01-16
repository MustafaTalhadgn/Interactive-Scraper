package scorer

import (
	"fmt"
	"log/slog"
	"strings"

	"InteractiveScraper/internal/extractor"
)

type Scorer struct {
	weights    *Weights
	thresholds *CriticalityThresholds
	rules      []ScoringRule
	logger     *slog.Logger
}

type Config struct {
	Weights    *Weights
	Thresholds *CriticalityThresholds
}

func DefaultConfig() *Config {
	return &Config{
		Weights:    DefaultWeights(),
		Thresholds: DefaultThresholds(),
	}
}

func NewScorer(config *Config, logger *slog.Logger) *Scorer {
	if config == nil {
		config = DefaultConfig()
	}

	scorer := &Scorer{
		weights:    config.Weights,
		thresholds: config.Thresholds,
		logger:     logger,
	}

	scorer.rules = []ScoringRule{
		&KeywordRule{weights: scorer.weights},
		&CryptoRule{weights: scorer.weights},
		&CVERule{weights: scorer.weights},
		&ContactRule{weights: scorer.weights},
		&NetworkRule{weights: scorer.weights},
		&ContentRule{weights: scorer.weights},
	}

	return scorer
}

type ScoreResult struct {
	Score            int
	Criticality      Criticality
	CriticalityColor string
	Breakdown        *ScoreBreakdown
}

type ScoreBreakdown struct {
	KeywordScore  int
	CryptoScore   int
	CVEScore      int
	ContactScore  int
	NetworkScore  int
	ContentScore  int
	TotalRawScore int
}

func (s *Scorer) Score(features *extractor.Features, content string) *ScoreResult {
	breakdown := &ScoreBreakdown{}

	for _, rule := range s.rules {
		ruleScore := rule.Apply(features, content)

		switch rule.(type) {
		case *KeywordRule:
			breakdown.KeywordScore = ruleScore
		case *CryptoRule:
			breakdown.CryptoScore = ruleScore
		case *CVERule:
			breakdown.CVEScore = ruleScore
		case *ContactRule:
			breakdown.ContactScore = ruleScore
		case *NetworkRule:
			breakdown.NetworkScore = ruleScore
		case *ContentRule:
			breakdown.ContentScore = ruleScore
		}
	}

	rawScore := breakdown.KeywordScore +
		breakdown.CryptoScore +
		breakdown.CVEScore +
		breakdown.ContactScore +
		breakdown.NetworkScore +
		breakdown.ContentScore

	breakdown.TotalRawScore = rawScore

	finalScore := rawScore
	if finalScore > 100 {
		finalScore = 100
	}
	if finalScore < 0 {
		finalScore = 0
	}

	criticality := ClassifyCriticality(finalScore, s.thresholds)

	s.logger.Info("scoring completed",
		slog.Int("score", finalScore),
		slog.String("criticality", string(criticality)),
		slog.Int("raw_score", rawScore),
		slog.Int("keyword_score", breakdown.KeywordScore),
		slog.Int("crypto_score", breakdown.CryptoScore),
		slog.Int("cve_score", breakdown.CVEScore),
	)

	return &ScoreResult{
		Score:            finalScore,
		Criticality:      criticality,
		CriticalityColor: GetCriticalityColor(criticality),
		Breakdown:        breakdown,
	}
}

func (s *Scorer) ScoreWithExplanation(features *extractor.Features, content string) (*ScoreResult, string) {
	result := s.Score(features, content)

	explanation := s.generateExplanation(result, features)

	return result, explanation
}

func (s *Scorer) generateExplanation(result *ScoreResult, features *extractor.Features) string {
	var parts []string

	if result.Breakdown.KeywordScore > 0 {
		parts = append(parts,
			sprintf("Kritik anahtar kelimeler tespit edildi (+%d)", result.Breakdown.KeywordScore))
	}

	if result.Breakdown.CryptoScore > 0 {
		parts = append(parts,
			sprintf("Kripto para adresleri bulundu (+%d)", result.Breakdown.CryptoScore))
	}

	if result.Breakdown.CVEScore > 0 {
		parts = append(parts,
			sprintf("%d CVE tanımlayıcısı(s) (+%d)", len(features.CVEs), result.Breakdown.CVEScore))
	}

	if result.Breakdown.NetworkScore > 0 {
		parts = append(parts,
			sprintf("Ağ göstergeleri (.onion, IP'ler) (+%d)", result.Breakdown.NetworkScore))
	}

	if len(parts) == 0 {
		return "Düşük tehdit göstergeleri tespit edildi"
	}

	return strings.Join(parts, "; ")
}

func sprintf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
