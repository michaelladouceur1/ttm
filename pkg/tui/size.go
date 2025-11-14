package tui

import "ttm/pkg/tui/context"

const (
	footerHeight        = 2
	sectionHeightOffset = 4
	sectionWidthOffset  = 6
)

func calculateSectionDims(width, height int) (leftDims, middleDims, rightDims, footerDims context.SectionDims) {
	leftDims = context.SectionDims{
		Width:  width / 4,
		Height: height - sectionHeightOffset,
	}
	middleDims = context.SectionDims{
		Width:  width/2 - sectionWidthOffset,
		Height: height - sectionHeightOffset,
	}
	rightDims = context.SectionDims{
		Width:  width / 4,
		Height: height - sectionHeightOffset,
	}
	footerDims = context.SectionDims{
		Width:  width,
		Height: footerHeight,
	}
	return
}
