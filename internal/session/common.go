package session

import (
	"github.com/vocabmaster/vocabmaster/internal/library"
	"github.com/vocabmaster/vocabmaster/internal/model"
	"github.com/vocabmaster/vocabmaster/internal/store"
)

type Config struct {
	Lang     string
	Level    model.DifficultyLevel
	Count    int
	NewWords int
}

type SessionResult struct {
	Reviewed int
	Learned  int
	Correct  int
	Total    int
}

func getAllWordIDs(lib *library.Library) []string {
	words := lib.AllWords()
	ids := make([]string, len(words))
	for i, w := range words {
		ids[i] = w.ID
	}
	return ids
}

func resolveEnrichedLinkedWords(s store.Store, linkedWords []*model.Word) []*model.Word {
	out := make([]*model.Word, len(linkedWords))
	copy(out, linkedWords)
	for i, lw := range out {
		if enriched, err := s.GetEnrichedWord(lw.ID); err == nil && enriched != nil {
			out[i] = enriched
		}
	}
	return out
}
