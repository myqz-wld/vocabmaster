package model

type Language string

const (
	LangEnglish  Language = "en"
	LangJapanese Language = "ja"
)

type DifficultyLevel int

const (
	Beginner     DifficultyLevel = 1
	Intermediate DifficultyLevel = 2
	Advanced     DifficultyLevel = 3
)

func (d DifficultyLevel) String() string {
	switch d {
	case Beginner:
		return "初级"
	case Intermediate:
		return "中级"
	case Advanced:
		return "高级"
	default:
		return "未知"
	}
}

type Example struct {
	Sentence    string `json:"sentence"`
	Translation string `json:"translation"`
}

type Word struct {
	ID            string          `json:"id"`
	Language      Language        `json:"language"`
	Text          string          `json:"text"`
	Pronunciation string          `json:"pronunciation"`
	ChineseDef    string          `json:"chinese_def"`
	Difficulty    DifficultyLevel `json:"difficulty"`
	PartOfSpeech  string          `json:"part_of_speech"`
	Examples      []Example       `json:"examples"`
	LinkedWordIDs []string        `json:"linked_word_ids"`
	Tags          []string        `json:"tags,omitempty"`
}

type WordLibraryFile struct {
	Version    string `json:"version"`
	Language   string `json:"language"`
	Words      []Word `json:"words"`
}
