package gridmap

import (
	"github.com/memmaker/ECon/common"
	"github.com/memmaker/ECon/geometry"
)

type MapCell struct {
	Icon            rune
	BackgroundColor common.RGBColor
	ForegroundColor common.RGBColor
	IsOpaque        bool
}
type Actor struct {
	Pos  geometry.Point
	Icon rune
}

type LightSource struct {
	Pos          geometry.Point
	Radius       int
	Color        common.RGBColor
	MaxIntensity float64
}

type GridMap struct {
	cells  []MapCell
	actors map[geometry.Point]*Actor
	width  int
	height int
}

func NewMap(width, height int) *GridMap {
	return &GridMap{
		width:  width,
		height: height,
		cells:  make([]MapCell, width*height),
		actors: make(map[geometry.Point]*Actor),
	}
}

func (m *GridMap) GetCell(p geometry.Point) MapCell {
	return m.cells[p.X+p.Y*m.width]
}

func (m *GridMap) SetCell(p geometry.Point, cell MapCell) {
	m.cells[p.X+p.Y*m.width] = cell
}

func (m *GridMap) GetActor(p geometry.Point) *Actor {
	return m.actors[p]
}

func (m *GridMap) AddActor(actor *Actor) {
	m.actors[actor.Pos] = actor
}

func (m *GridMap) RemoveActor(actor *Actor) {
	delete(m.actors, actor.Pos)
}

func (m *GridMap) MoveActor(actor *Actor, newPos geometry.Point) {
	delete(m.actors, actor.Pos)
	actor.Pos = newPos
	m.actors[newPos] = actor
}

func (m *GridMap) Fill(mapCell MapCell) {
	for i := range m.cells {
		m.cells[i] = mapCell
	}
}

func (m *GridMap) Iterate(f func(p geometry.Point, cell MapCell)) {
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			p := geometry.Point{X: x, Y: y}
			f(p, m.cells[p.X+p.Y*m.width])
		}
	}
}
func (m *GridMap) IsTransparent(p geometry.Point) bool {
	return !m.GetCell(p).IsOpaque
}

func (m *GridMap) Contains(dest geometry.Point) bool {
	return dest.X >= 0 && dest.X < m.width && dest.Y >= 0 && dest.Y < m.height
}
