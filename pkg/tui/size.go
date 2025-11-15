package tui

import "ttm/pkg/tui/context"

const (
	headerHeight        = 1
	footerHeight        = 1
	sectionHeightOffset = 0
	sectionWidthOffset  = 6
)

func calculateSectionDims(width, height int) context.SectionDims {

	return context.SectionDims{
		Header: context.SectionDim{
			Width:  width,
			Height: headerHeight,
		},
		Left: context.SectionDim{
			Width:  width / 4,
			Height: height - headerHeight - footerHeight - sectionHeightOffset,
		},
		Middle: context.SectionDim{
			Width:  width/2 - sectionWidthOffset,
			Height: height - headerHeight - footerHeight - sectionHeightOffset,
		},
		Right: context.SectionDim{
			Width:  width / 4,
			Height: height - headerHeight - footerHeight - sectionHeightOffset,
		},
		Footer: context.SectionDim{
			Width:  width,
			Height: footerHeight,
		},
	}
}
