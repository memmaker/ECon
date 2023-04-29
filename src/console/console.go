package console

import (
	"embed"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tinne26/etxt"

	"github.com/memmaker/ECon/common"
	"github.com/memmaker/ECon/geometry"
)

type CellInterface interface {
	Set(p geometry.Point, cell common.Cell)
	At(p geometry.Point) common.Cell
	Size() geometry.Point
	Flush()
	ClearScreen()
	Fill(bbox geometry.Rect, cell common.Cell)
}
type GridConfig struct {
	TileWidth  int
	TileHeight int
	GridWidth  int
	GridHeight int
}

type Console struct {
	// basics
	screenDPIScale float64
	TileWidth      int
	TileHeight     int
	// fonts
	txtRenderer *etxt.Renderer
	currentFont *etxt.Font
	// delta drawing
	drawGrid            geometry.Grid
	deltaGrid           geometry.Grid
	frameIsDirty        bool
	clearBeforeNextDraw bool

	nextFrame Frame
}

func (c *Console) Size() geometry.Point {
	return c.drawGrid.Size()
}
func (c *Console) Set(p geometry.Point, cell common.Cell) {
	//println("set", p.String(), cell.Char)
	c.drawGrid.Set(p, cell)
}

func (c *Console) At(p geometry.Point) common.Cell {
	return c.drawGrid.At(p)
}

func NewConsole(config GridConfig) *Console {
	return &Console{
		TileWidth:   config.TileWidth,
		TileHeight:  config.TileHeight,
		txtRenderer: NewTextRenderer(),
		drawGrid:    geometry.NewGrid(config.GridWidth, config.GridHeight),
	}
}

func (c *Console) SetFont(font *etxt.Font) {
	c.txtRenderer.SetFont(font)
	c.currentFont = font
}
func (c *Console) ClearScreen() {
	c.drawGrid.Fill(common.Cell{Char: ' ', Foreground: common.White, Background: common.Black})
	//c.nextFrame = Frame{}
	c.clearBeforeNextDraw = true
	c.frameIsDirty = true
}
func (c *Console) Draw(screen *ebiten.Image) {
	if !c.frameIsDirty || c.currentFont == nil {
		return
	}
	if c.clearBeforeNextDraw {
		c.clearBeforeNextDraw = false
		screen.Fill(common.Black)
	}
	//println("drawing")
	scale := c.screenDPIScale
	tilewidth := int(math.Ceil(float64(c.TileWidth) * scale))
	tileheight := int(math.Ceil(float64(c.TileHeight) * scale))

	c.drawGrid.Iter(func(p geometry.Point, cell common.Cell) {
		vector.DrawFilledRect(screen, float32((p.X)*tilewidth), float32((p.Y)*tileheight), float32(tilewidth), float32(tileheight), cell.Background)
	})
	/*
		for _, cellAt := range c.nextFrame.Cells {
			vector.DrawFilledRect(screen, float32((cellAt.P.X)*tilewidth), float32((cellAt.P.Y)*tileheight), float32(tilewidth), float32(tileheight), cellAt.Cell.Background)
		}*/

	c.txtRenderer.SetTarget(screen)
	c.txtRenderer.SetSizePx(int(math.Ceil(float64(c.TileHeight) * scale)))
	c.drawGrid.Iter(func(p geometry.Point, cell common.Cell) {
		glyph := string(cell.Char)
		runes, _ := etxt.GetMissingRunes(c.currentFont, glyph)
		if len(runes) > 0 {
			glyph = " "
		}
		xPos := (p.X) * tilewidth
		yPos := (p.Y) * tileheight
		c.txtRenderer.SetColor(cell.Foreground)
		c.txtRenderer.Draw(glyph, xPos, yPos)
	})
	c.frameIsDirty = false
}

func (c *Console) SetScale(scale float64) {
	if scale != c.screenDPIScale {
		c.screenDPIScale = scale
		c.ClearScreen()
	}
}

// Flush will compute the delta between the last frame and the current frame and
// record the frame if recording is enabled. Call it after all your drawing is
// done. And before you call Draw()
func (c *Console) Flush() {
	c.computeAndRecordNextFrame()
}

func (c *Console) Fill(rect geometry.Rect, cell common.Cell) {
	c.drawGrid.Slice(rect).Fill(cell)
}

// computeAndRecordNextFrame will compute the delta between the last frame and the current frame and trigger a new tick for the recording.
func (c *Console) computeAndRecordNextFrame() {
	forceRedrawOfAllCells := c.clearBeforeNextDraw
	frame := c.computeFrame(c.drawGrid, forceRedrawOfAllCells)
	if len(frame.Cells) > 0 {
		c.frameIsDirty = true
	}
}

func (c *Console) computeFrame(gd geometry.Grid, exposed bool) Frame {
	if gd.Ug == nil || gd.Rg.Empty() && !exposed {
		return Frame{}
	}
	if c.deltaGrid.Ug == nil {
		c.deltaGrid = geometry.NewGrid(gd.Ug.Width, gd.Ug.Height)
	} else if c.deltaGrid.Ug.Width != gd.Ug.Width || c.deltaGrid.Ug.Height != gd.Ug.Height {
		c.deltaGrid = c.deltaGrid.Resize(gd.Ug.Width, gd.Ug.Height)
	}
	c.nextFrame.Cells = c.nextFrame.Cells[:0]
	if exposed {
		return c.refresh(gd)
	}
	w := gd.Ug.Width
	cells := gd.Ug.Cells
	pcells := c.deltaGrid.Ug.Cells // previous cells
	yimax := gd.Rg.Max.Y * w
	for y, yi := 0, gd.Rg.Min.Y*w; yi < yimax; y, yi = y+1, yi+w {
		ximax := yi + gd.Rg.Max.X
		for x, xi := 0, yi+gd.Rg.Min.X; xi < ximax; x, xi = x+1, xi+1 {
			cellAt := cells[xi]
			if cellAt == pcells[xi] {
				continue
			}
			pcells[xi] = cellAt
			p := geometry.Point{X: x, Y: y}
			cdraw := FrameCell{Cell: cellAt, P: p}
			c.nextFrame.Cells = append(c.nextFrame.Cells, cdraw)
		}
	}
	return c.nextFrame
}

func (c *Console) refresh(gd geometry.Grid) Frame {
	gd.Rg.Min = geometry.Point{}
	gd.Rg.Max = gd.Rg.Min.Add(geometry.Point{X: gd.Ug.Width, Y: gd.Ug.Height})
	c.deltaGrid.Copy(gd)
	it := gd.Iterator()
	for it.Next() {
		cdraw := FrameCell{Cell: it.Cell(), P: it.P()}
		c.nextFrame.Cells = append(c.nextFrame.Cells, cdraw)
	}
	return c.nextFrame
}

func (c *Console) LoadEmbeddedFont(fontDir string, fontName string, fs embed.FS) {
	fontLib := etxt.NewFontLibrary()
	_, _, err := fontLib.ParseEmbedDirFonts(fontDir, fs)
	if err != nil {
		log.Fatalf("Error while loading EmbeddedData: %s", err.Error())
	}

	// check that we have the EmbeddedData we want
	// (shown for completeness, you don't need this in most cases)
	expectedFonts := []string{fontName}
	for _, expectedFontName := range expectedFonts {
		if !fontLib.HasFont(expectedFontName) {
			log.Fatal("missing expectedFontName: " + expectedFontName)
		}
	}
	loadedFont := fontLib.GetFont(expectedFonts[0])
	c.SetFont(loadedFont)

}

func NewTextRenderer() *etxt.Renderer {
	txtRenderer := etxt.NewStdRenderer()
	glyphsCache := etxt.NewDefaultCache(10 * 1024 * 1024) // 10MB
	txtRenderer.SetCacheHandler(glyphsCache.NewHandler())
	txtRenderer.SetAlign(etxt.Top, etxt.Left)
	whiteColor := common.White
	txtRenderer.SetColor(whiteColor)
	return txtRenderer
}
