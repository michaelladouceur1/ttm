package logger

import (
	"fmt"
	"os"
	"strings"
	"ttm/pkg/models"

	"github.com/charmbracelet/x/term"
)

const (
	ClassicStyleName = "classic"
	CompactStyleName = "compact"
	MinimalStyleName = "minimal"
)

type SummaryItem struct {
	Key   string
	Value string
}

// Style controls how logger output is rendered.
type Style interface {
	RenderMessage(message string) string
	RenderError(message string) string
	RenderSummary(title string, items []SummaryItem) string
	RenderTasks(tasks []models.Task) string
}

type Logger struct {
	style Style
}

func New(style Style) *Logger {
	return &Logger{style: style}
}

func NewStyle(name string) (Style, error) {
	width, height, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		width = 80
		height = 24
	}
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "", ClassicStyleName:
		return ClassicStyle{width, height}, nil
	case CompactStyleName:
		return CompactStyle{}, nil
	case MinimalStyleName:
		return MinimalStyle{width, height}, nil
	default:
		return nil, fmt.Errorf("unsupported logging theme %q; choose classic or compact", name)
	}
}

var defaultLogger = New(ClassicStyle{})

func Configure(styleName string) error {
	style, err := NewStyle(styleName)
	if err != nil {
		return err
	}
	defaultLogger = New(style)
	return nil
}

func (l *Logger) LogMessage(strs ...string) {
	fmt.Println(l.style.RenderMessage(strings.Join(strs, "")))
}

func (l *Logger) LogError(strs ...any) {
	fmt.Println(l.style.RenderError(fmt.Sprint(strs...)))
}

func LogMessage(strs ...string) {
	defaultLogger.LogMessage(strs...)
}

func LogError(strs ...any) {
	defaultLogger.LogError(strs...)
}
