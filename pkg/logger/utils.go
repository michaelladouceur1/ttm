package logger

import (
	"strings"

	"github.com/charmbracelet/lipgloss/tree"
)

type SummaryTreeItem struct {
	Key   string
	Value string
}

func createSummaryTree(data []SummaryTreeItem, title string) *tree.Tree {
	longestKey := getLongestKeyLength(data)
	children := getTreeChildStrings(data, longestKey)
	return tree.Root("⚙ " + title).
		Child(children...).
		Enumerator(tree.RoundedEnumerator).
		EnumeratorStyle(treeConnStyle).
		RootStyle(headerStyle).
		ItemStyle(textStyle)
}

func getTreeChildStrings(data []SummaryTreeItem, longestKey int) []any {
	children := make([]any, 0, len(data))
	for _, kv := range data {
		paddingLength := longestKey - len(kv.Key) + SeparatorMargin
		childString := kv.Key + treeConnStyle.Render(" "+strings.Repeat(Separator, paddingLength)) + " " + textStyle.Render(kv.Value)
		children = append(children, childString)
	}
	return children
}

func getLongestKeyLength(data []SummaryTreeItem) int {
	maxLen := 0
	for _, kv := range data {
		if len(kv.Key) > maxLen {
			maxLen = len(kv.Key)
		}
	}
	return maxLen
}
