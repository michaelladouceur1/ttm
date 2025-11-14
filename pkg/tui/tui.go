package tui

import (
	"os"
	"ttm/pkg/config"
	"ttm/pkg/tui/context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/michaelladouceur1/gonfig"
	"golang.org/x/term"
)

type TuiPages string

const (
	RootPage  TuiPages = "root"
	AddPage   TuiPages = "add"
	StartPage TuiPages = "start"
)

type Pages map[TuiPages]tea.Model

var pages = Pages{
	// RootPage:  NewRootModel(),
	AddPage:   NewAddModel(),
	StartPage: NewStartModel(),
}

type NavItem struct {
	title, page string
}

func (i NavItem) Title() string       { return i.title }
func (i NavItem) Description() string { return "" }
func (i NavItem) FilterValue() string { return i.title }

type TUI struct {
	program *tea.Program
	ctx     *context.TUIContext
	root    RootModel
	pages   Pages
}

var tui *TUI

func NewTUI(config *gonfig.Gonfig[config.Config]) *TUI {

	tui = &TUI{
		pages: pages,
	}

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 100
		height = 30
	}

	leftDims, middleDims, rightDims, footerDims := calculateSectionDims(width, height)

	ctx := &context.TUIContext{
		Config:     config,
		TermWidth:  width,
		TermHeight: height,
		LeftDims:   leftDims,
		MiddleDims: middleDims,
		RightDims:  rightDims,
		FooterDims: footerDims,
	}

	tui.ctx = ctx
	tui.root = NewRootModel(ctx)

	tui.program = tea.NewProgram(tui.root, tea.WithAltScreen())

	return tui
}

func (t *TUI) Run() error {
	if _, err := t.program.Run(); err != nil {
		return err
	}
	return nil
}

func (t *TUI) switchPage(page TuiPages) (tea.Model, tea.Cmd) {
	var model tea.Model
	if m, ok := t.pages[page]; ok {
		model = m
	} else {
		model = pages[RootPage]
	}
	return model.Update(tea.KeyMsg{})
}

func defaultUpdate(m tea.Model, message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return nil, nil
}
