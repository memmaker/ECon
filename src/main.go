package main

import (
	"embed"
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/memmaker/ECon/common"
	"github.com/memmaker/ECon/console"
	"github.com/memmaker/ECon/game"
	"github.com/memmaker/ECon/geometry"
	"github.com/memmaker/ECon/input"
)

//go:embed embedded
var embeddedFS embed.FS

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

type Game struct {
	// Config
	Config console.GridConfig
	// Input
	Input *InputState
	// Console
	Console *console.Console
	// Model
	Model          *game.Model
	deviceDPIScale float64
}

func (g *Game) GetInput() input.GridInput {
	return g.Input
}

func (g *Game) Update() error {
	g.pollInput()
	g.Model.Update(g)       // This is our model's update() call
	g.Model.Draw(g.Console) // This is our model's draw() call
	g.Console.Flush()
	return nil
}

func (g *Game) pollInput() {
	// mouse
	g.Input.LastMousePos = g.Input.MousePos
	mx, my := ebiten.CursorPosition()
	xMouse := float64(mx) / (float64(g.Config.TileWidth) * g.deviceDPIScale)
	yMouse := float64(my) / (float64(g.Config.TileHeight) * g.deviceDPIScale)
	xMouse = common.Clamp(xMouse, 0, float64(g.Config.GridWidth-1))
	yMouse = common.Clamp(yMouse, 0, float64(g.Config.GridHeight-1))
	g.Input.MousePos = geometry.Point{X: int(xMouse), Y: int(yMouse)}
}

// This is the draw() function of ebitengine. It will draw the console to the screen.
func (g *Game) Draw(screen *ebiten.Image) {
	g.Console.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	panic("should use layoutf")
}

func (g *Game) LayoutF(outsideWidth, outsideHeight float64) (screenWidth, screenHeight float64) {
	scale := ebiten.DeviceScaleFactor()
	g.deviceDPIScale = scale
	g.Console.SetScale(scale)
	return float64(g.Config.GridWidth*g.Config.TileWidth) * scale, float64(g.Config.GridHeight*g.Config.TileHeight) * scale
}

func (g *Game) Init() {
	g.deviceDPIScale = ebiten.DeviceScaleFactor()
	g.Model.Init(g)
}

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}

	gameTitle := "E-Console"
	fontDirectory := "embedded/font"
	fontName := "Square"

	config := console.GridConfig{
		TileWidth:  20,
		TileHeight: 20,
		GridWidth:  64,
		GridHeight: 36,
	}

	con := console.NewConsole(config)
	con.LoadEmbeddedFont(fontDirectory, fontName, embeddedFS)
	consoleGame := &Game{
		Config:  config,
		Console: con,
		Input:   NewInput(),
		Model:   game.NewModel(config),
	}
	ebiten.SetWindowTitle(gameTitle)
	ebiten.SetWindowSize(int(float64(config.GridWidth*config.TileWidth)), int(float64(config.GridHeight*config.TileHeight)))
	ebiten.SetScreenClearedEveryFrame(false)
	consoleGame.Init()
	if err := ebiten.RunGameWithOptions(consoleGame, &ebiten.RunGameOptions{
		GraphicsLibrary: ebiten.GraphicsLibraryOpenGL,
	}); err != nil {
		log.Fatal(err)
	}
}
