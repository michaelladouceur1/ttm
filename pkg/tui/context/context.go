package context

import (
	"ttm/pkg/config"

	"github.com/michaelladouceur1/gonfig"
)

type SectionDims struct {
	Width  int
	Height int
}

type TUIContext struct {
	Config     *gonfig.Gonfig[config.Config]
	TermWidth  int
	TermHeight int
	LeftDims   SectionDims
	MiddleDims SectionDims
	RightDims  SectionDims
	FooterDims SectionDims
}
