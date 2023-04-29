package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/memmaker/ECon/geometry"
)

type InputState struct {
	MousePos     geometry.Point
	LastMousePos geometry.Point
}

func (i InputState) IsMenuClose() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEscape)
}

func (i InputState) IsMenuConfirm() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) || i.IsMouseLeft()
}

func (i InputState) IsMenuDown() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS)
}

func (i InputState) IsMenuUp() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW)
}

func (i InputState) GetMousePos() geometry.Point {
	return i.MousePos
}
func (i InputState) HasMouseMoved() bool {
	return i.MousePos != i.LastMousePos
}
func (i InputState) IsMouseLeft() bool {
	return inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
}
func (i InputState) IsMouseRight() bool {
	return inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)
}

func (i InputState) GetJustPressedKeys() []string {
	buffer := make([]ebiten.Key, 0)
	buffer = inpututil.AppendJustPressedKeys(buffer)
	keys := make([]string, len(buffer))
	for index, key := range buffer {
		keys[index] = key.String()
	}
	return keys
}

func NewInput() *InputState {
	return &InputState{}
}
