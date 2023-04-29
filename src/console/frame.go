package console

import (
	"github.com/memmaker/ECon/common"
	"github.com/memmaker/ECon/geometry"
)

type Frame struct {
	Cells       []FrameCell // cells that changed from previous squareDeltaFrame
	IsHalfWidth bool
}
type FrameCell struct {
	Cell common.Cell    // cell content and styling
	P    geometry.Point // absolute position in the whole squarePreviousGrid
}
