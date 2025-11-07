package context

import (
	"ttm/pkg/config"

	"github.com/michaelladouceur1/gonfig"
)

type TUIContext struct {
	Config *gonfig.Gonfig[config.Config]
	Width  int
	Height int
}
