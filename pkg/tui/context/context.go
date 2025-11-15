package context

import (
	"ttm/pkg/config"

	"github.com/michaelladouceur1/gonfig"
)

type SectionDim struct {
	Width  int
	Height int
}

type SectionDims struct {
	Header SectionDim
	Left   SectionDim
	Middle SectionDim
	Right  SectionDim
	Footer SectionDim
}

type TUIContext struct {
	Config     *gonfig.Gonfig[config.Config]
	TermWidth  int
	TermHeight int
	Dims       SectionDims
}
