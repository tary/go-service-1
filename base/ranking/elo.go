package ranking

import (
	"math"
)

const (
	// K is the default K-Factor
	K = 32
	// D is the default deviation
	D = 400
)

// Elo calculates Elo rating changes based on the configured factors.
type Elo struct {
	K int32
	D int32
}

// Outcome is a match result data for a single player.
type Outcome struct {
	Delta  int32
	Rating int32
}

// NewElo instantiates the Elo object with default factors.
// Default K-Factor is 32
// Default deviation is 400
func NewElo() *Elo {
	return &Elo{K, D}
}

// NewEloWithFactors instantiates the Elo object with custom factor values.
func NewEloWithFactors(k, d int32) *Elo {
	return &Elo{k, d}
}

// ExpectedScore gives the expected chance that the first player wins
func (e *Elo) ExpectedScore(ratingA, ratingB int32) float64 {
	return e.ExpectedScoreWithFactors(ratingA, ratingB, e.D)
}

// ExpectedScoreWithFactors overrides default factors and gives the expected chance that the first player wins
func (e *Elo) ExpectedScoreWithFactors(ratingA, ratingB, d int32) float64 {
	return 1 / (1 + math.Pow(10, float64(ratingB-ratingA)/float64(d)))
}

// RatingDelta gives the ratings change for the first player for the given score
func (e *Elo) RatingDelta(ratingA, ratingB int32, score float64) int32 {
	return e.RatingDeltaWithFactors(ratingA, ratingB, score, e.K, e.D)
}

// RatingDeltaWithFactors overrides default factors and gives the ratings change for the first player for the given score
func (e *Elo) RatingDeltaWithFactors(ratingA, ratingB int32, score float64, k, d int32) int32 {
	return int32(float64(k) * (score - e.ExpectedScoreWithFactors(ratingA, ratingB, d)))
}

// Rating gives the new rating for the first player for the given score
func (e *Elo) Rating(ratingA, ratingB int32, score float64) int32 {
	return e.RatingWithFactors(ratingA, ratingB, score, e.K, e.D)
}

// RatingWithFactors overrides default factors and gives the new rating for the first player for the given score
func (e *Elo) RatingWithFactors(ratingA, ratingB int32, score float64, k, d int32) int32 {
	return ratingA + e.RatingDeltaWithFactors(ratingA, ratingB, score, k, d)
}

// Outcome gives an Outcome object for each player for the given score
func (e *Elo) Outcome(ratingA, ratingB int32, score float64) (Outcome, Outcome) {
	return e.OutcomeWithFactors(ratingA, ratingB, score, e.K, e.D)
}

// OutcomeWithFactors overrides default factors and gives an Outcome object for each player for the given score
func (e *Elo) OutcomeWithFactors(ratingA, ratingB int32, score float64, k, d int32) (Outcome, Outcome) {
	delta := e.RatingDeltaWithFactors(ratingA, ratingB, score, k, d)
	return Outcome{delta, ratingA + delta}, Outcome{-delta, ratingB - delta}
}
