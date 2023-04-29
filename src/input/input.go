package input

import "github.com/memmaker/ECon/geometry"

type GridInput interface {
	GetMousePos() geometry.Point
	HasMouseMoved() bool

	IsMouseLeft() bool
	IsMouseRight() bool

	IsMenuClose() bool
	IsMenuConfirm() bool
	IsMenuDown() bool
	IsMenuUp() bool

	GetJustPressedKeys() []string
}
