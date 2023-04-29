package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/memmaker/ECon/common"
	"github.com/memmaker/ECon/console"
	"github.com/memmaker/ECon/geometry"
	"github.com/memmaker/ECon/gridmap"
)

type Model struct {
	oldMousePos geometry.Point
	playerPos   geometry.Point
	config      console.GridConfig
	gridMap     *gridmap.GridMap
	player      *gridmap.Actor
	clearScreen bool
}

func NewModel(config console.GridConfig) *Model {
	model := &Model{
		config:  config,
		gridMap: gridmap.NewMap(config.GridWidth, config.GridHeight),
	}
	return model
}

// Update is called every frame
func (m *Model) Update(engine console.Engine) {
	userInput := engine.GetInput()

	newMousePos := userInput.GetMousePos()
	if newMousePos != m.oldMousePos {
		m.gridMap.MoveActor(m.player, newMousePos)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		m.clearScreen = true
	}
	if userInput.IsMouseLeft() {
		m.PlaceWall(newMousePos)
		//m.PlaceLight(newMousePos)
	}
	m.oldMousePos = newMousePos
}

// Draw is called every frame
func (m *Model) Draw(con console.CellInterface) {
	if m.clearScreen {
		con.ClearScreen()
		m.clearScreen = false
	}

	m.drawMap(con)
}

func (m *Model) drawMap(con console.CellInterface) {
	m.gridMap.Iterate(func(p geometry.Point, cell gridmap.MapCell) {
		drawCell := m.drawCell(p, cell)
		con.Set(p, drawCell)
	})
}

func (m *Model) drawCell(p geometry.Point, cell gridmap.MapCell) common.Cell {
	drawRune := cell.Icon
	actorAt := m.gridMap.GetActor(p)
	if actorAt != nil {
		drawRune = actorAt.Icon
	}

	return common.Cell{Char: drawRune, Foreground: cell.ForegroundColor, Background: cell.BackgroundColor}
}

var groundCell = gridmap.MapCell{Icon: '.', ForegroundColor: common.RGBColor{R: 0.8, G: 0.8, B: 0.8}, BackgroundColor: common.RGBColor{R: 97 / 255.0, G: 158 / 255.0, B: 1.0}}
var wallCell = gridmap.MapCell{Icon: '#', IsOpaque: true, ForegroundColor: common.RGBColor{R: 0.8, G: 0.8, B: 0.8}, BackgroundColor: common.RGBColor{R: 0.9, G: 0.9, B: 0.9}}

func (m *Model) Init(engine console.Engine) {

	m.gridMap.Fill(groundCell)
	playerSpawn := geometry.Point{X: 10, Y: 10}
	m.player = &gridmap.Actor{
		Icon: '@',
		Pos:  playerSpawn,
	}
	m.gridMap.AddActor(m.player)
}

func (m *Model) PlaceWall(pos geometry.Point) {
	dest := pos.Add(geometry.Point{Y: 1})
	m.gridMap.SetCell(dest, wallCell)
}
