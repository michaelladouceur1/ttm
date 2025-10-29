package logger

import (
	"fmt"
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
	children := []any{}
	for _, kv := range data {
		paddingLength := longestKey - len(kv.Key) + SeparatorMargin
		childString := fmt.Sprintf(kv.Key) + fmt.Sprintf(treeConnStyle.Render(" "+strings.Repeat(Separator, paddingLength))+" ") + fmt.Sprintf(textStyle.Render(kv.Value))
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
