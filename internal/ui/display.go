package ui

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/vocabmaster/vocabmaster/internal/model"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorBold   = "\033[1m"
)

func DisplayWordCard(w *model.Word) {
	width := 50

	fmt.Println(colorGray + strings.Repeat("─", width) + colorReset)
	fmt.Printf("  %s%s%s  %s%s%s\n",
		colorCyan+colorBold, w.Text, colorReset,
		colorGray, w.Pronunciation, colorReset)
	fmt.Printf("  %s[%s]%s  %s[%s]%s\n",
		colorGray, w.PartOfSpeech, colorReset,
		colorGray, w.Difficulty.String(), colorReset)
	fmt.Println()
	fmt.Printf("  %s中文：%s%s\n", colorBold, colorReset, w.ChineseDef)
	fmt.Println()

	if len(w.Examples) > 0 {
		fmt.Printf("  %s例句：%s\n", colorBold, colorReset)
		for _, ex := range w.Examples {
			fmt.Printf("  %s%s%s\n", colorGreen, ex.Sentence, colorReset)
			fmt.Printf("  %s%s%s\n", colorGray, ex.Translation, colorReset)
			fmt.Println()
		}
	}

	fmt.Println(colorGray + strings.Repeat("─", width) + colorReset)
}

func DisplayPrompt(text string) {
	fmt.Printf("\n  %s%s%s\n", colorYellow, text, colorReset)
}

func DisplayReveal(w *model.Word) {
	fmt.Printf("\n  %s答案：%s %s%s%s %s%s%s\n",
		colorBold, colorReset,
		colorCyan+colorBold, w.Text, colorReset,
		colorGray, w.Pronunciation, colorReset)
	fmt.Printf("  %s%s%s\n", colorGreen, w.ChineseDef, colorReset)

	if len(w.Examples) > 0 {
		fmt.Printf("  %s%s%s\n", colorGray, w.Examples[0].Sentence, colorReset)
	}
}

func DisplayGradeOptions() {
	fmt.Printf("\n  评分：%s[0] Again%s  %s[1] Hard%s  %s[2] Good%s  %s[3] Easy%s\n",
		colorRed, colorReset,
		colorYellow, colorReset,
		colorGreen, colorReset,
		colorCyan, colorReset)
}

func DisplaySessionSummary(reviewed, learned, correct, total int) {
	width := 40
	fmt.Println()
	fmt.Println(colorGray + strings.Repeat("─", width) + colorReset)
	fmt.Printf("  %s学习完成！%s\n", colorBold+colorGreen, colorReset)

	if reviewed > 0 {
		fmt.Printf("  复习：%d 词\n", reviewed)
	}
	if learned > 0 {
		fmt.Printf("  新学：%d 词\n", learned)
	}
	if total > 0 {
		rate := float64(correct) / float64(total) * 100
		fmt.Printf("  正确率：%.1f%%\n", rate)
	}
	fmt.Println(colorGray + strings.Repeat("─", width) + colorReset)
}

func DisplayStatsTable(totalLib int, stats *StatsData) {
	fmt.Printf("\n%s📊 学习统计%s\n\n", colorBold, colorReset)

	rows := [][]string{
		{"词库总量", fmt.Sprintf("%d 词", totalLib)},
		{"已学习", fmt.Sprintf("%d 词", stats.Learned)},
		{"已掌握", fmt.Sprintf("%d 词", stats.Mastered)},
		{"待复习", fmt.Sprintf("%d 词", stats.Due)},
		{"总复习次数", fmt.Sprintf("%d 次", stats.TotalReviews)},
		{"正确率", fmt.Sprintf("%.1f%%", stats.Accuracy)},
		{"连续打卡", fmt.Sprintf("%d 天", stats.Streak)},
	}

	maxLabelWidth := 0
	for _, row := range rows {
		w := runewidth.StringWidth(row[0])
		if w > maxLabelWidth {
			maxLabelWidth = w
		}
	}

	for _, row := range rows {
		label := row[0]
		padding := maxLabelWidth - runewidth.StringWidth(label) + 2
		fmt.Printf("  %s%s%s%s%s\n", colorGray, label, strings.Repeat(" ", padding), colorReset, row[1])
	}
	fmt.Println()
}

type StatsData struct {
	Learned      int
	Mastered     int
	Due          int
	TotalReviews int
	Accuracy     float64
	Streak       int
}
