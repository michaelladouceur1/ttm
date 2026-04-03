package context

import (
	"ttm/pkg/config"
	"ttm/pkg/store"

	"github.com/michaelladouceur1/gonfig"
)

type SectionDim struct {
	Width  int
	Height int
}

type SectionDims struct {
	Header   SectionDim
	Sessions SectionDim
	Left     SectionDim
	Middle   SectionDim
	Right    SectionDim
	Footer   SectionDim
}

type TUIContext struct {
	Config     *gonfig.Gonfig[config.Config]
	Store      *store.Store
	TermWidth  int
	TermHeight int
	Dims       SectionDims
}
