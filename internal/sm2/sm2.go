package sm2

import (
	"math"
	"time"

	"github.com/vocabmaster/vocabmaster/internal/model"
)

func Apply(record model.ReviewRecord, grade model.Grade, now time.Time) model.ReviewRecord {
	g := int(grade)

	var newInterval int
	var newReps int

	// Update EF on every review per standard SM-2
	newEF := record.EaseFactor + (0.1 - float64(5-g)*(0.08+float64(5-g)*0.02))
	if newEF < 1.3 {
		newEF = 1.3
	}

	if g >= 3 { // correct
		switch record.Repetitions {
		case 0:
			newInterval = 1
		case 1:
			newInterval = 6
		default:
			newInterval = int(math.Round(float64(record.Interval) * newEF))
		}
		newReps = record.Repetitions + 1
		record.CorrectCount++
	} else { // incorrect — reset repetitions
		newReps = 0
		newInterval = 1
	}
	record.EaseFactor = newEF

	record.Interval = newInterval
	record.Repetitions = newReps
	record.NextReviewAt = now.AddDate(0, 0, newInterval)
	record.LastReviewAt = now
	record.TotalReviews++

	return record
}
