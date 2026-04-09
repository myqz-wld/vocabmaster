package library

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/vocabmaster/vocabmaster/data"
	"github.com/vocabmaster/vocabmaster/internal/model"
)

type Library struct {
	words   map[string]*model.Word
	byLang  map[model.Language][]*model.Word
}

func NewLibrary() (*Library, error) {
	lib := &Library{
		words:  make(map[string]*model.Word),
		byLang: make(map[model.Language][]*model.Word),
	}

	for _, name := range []string{"english.json", "japanese.json"} {
		raw, err := data.DataFS.ReadFile(name)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", name, err)
		}

		var file model.WordLibraryFile
		if err := json.Unmarshal(raw, &file); err != nil {
			return nil, fmt.Errorf("parse %s: %w", name, err)
		}

		for i := range file.Words {
			w := &file.Words[i]
			lib.words[w.ID] = w
			lib.byLang[w.Language] = append(lib.byLang[w.Language], w)
		}
	}

	return lib, nil
}

func (l *Library) GetWord(id string) *model.Word {
	return l.words[id]
}

func (l *Library) AllWords() []*model.Word {
	result := make([]*model.Word, 0, len(l.words))
	for _, w := range l.words {
		result = append(result, w)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result
}

func (l *Library) FilterByLang(lang string) []*model.Word {
	if lang == "" || lang == "all" {
		return l.AllWords()
	}
	return l.byLang[model.Language(lang)]
}

func (l *Library) FilterByDifficulty(words []*model.Word, level model.DifficultyLevel) []*model.Word {
	if level == 0 {
		return words
	}
	var result []*model.Word
	for _, w := range words {
		if w.Difficulty == level {
			result = append(result, w)
		}
	}
	return result
}

func (l *Library) Search(query string) []*model.Word {
	query = strings.ToLower(query)
	var result []*model.Word
	for _, w := range l.words {
		if strings.Contains(strings.ToLower(w.Text), query) ||
			strings.Contains(strings.ToLower(w.ChineseDef), query) ||
			strings.Contains(strings.ToLower(w.Pronunciation), query) {
			result = append(result, w)
		}
	}
	return result
}


func (l *Library) TotalCount() int {
	return len(l.words)
}

func (l *Library) CountByLang(lang string) int {
	if lang == "" || lang == "all" {
		return len(l.words)
	}
	return len(l.byLang[model.Language(lang)])
}

func (l *Library) GetUnlearnedWords(lang string, level model.DifficultyLevel, learnedIDs map[string]bool) []*model.Word {
	words := l.FilterByLang(lang)
	words = l.FilterByDifficulty(words, level)

	var result []*model.Word
	for _, w := range words {
		if !learnedIDs[w.ID] {
			result = append(result, w)
		}
	}
	return result
}
