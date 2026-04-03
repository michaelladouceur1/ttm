package tui

import "ttm/pkg/tui/context"

const (
	headerHeight        = 1
	footerHeight        = 1
	sectionHeightOffset = 6
	sectionWidthOffset  = 6
)

func calculateSectionDims(width, height int) context.SectionDims {

	return context.SectionDims{
		Header: context.SectionDim{
			Width:  width - 2,
			Height: headerHeight,
		},
		Sessions: context.SectionDim{
			Width:  width - 2,
			Height: 5,
		},
		Left: context.SectionDim{
			Width:  width / 4,
			Height: height - headerHeight - footerHeight - sectionHeightOffset - 5,
		},
		Middle: context.SectionDim{
			Width:  width/2 - sectionWidthOffset,
			Height: height - headerHeight - footerHeight - sectionHeightOffset - 5,
		},
		Right: context.SectionDim{
			Width:  width / 4,
			Height: height - headerHeight - footerHeight - sectionHeightOffset - 5,
		},
		Footer: context.SectionDim{
			Width:  width - 2,
			Height: footerHeight,
		},
	}
}
