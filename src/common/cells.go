package common

import (
	"encoding/gob"
	"image/color"
)

type Cell struct {
	Foreground color.Color
	Background color.Color
	Char       rune
}

func (c Cell) WithChar(newChar int32) Cell {
	c.Char = newChar
	return c
}

func (c Cell) WithBackgroundColor(rgbColor RGBColor) Cell {
	c.Background = rgbColor
	return c
}

func (c Cell) WithForegroundColor(rgbColor RGBColor) Cell {
	c.Foreground = rgbColor
	return c
}

func init() {
	gob.Register(&RGBColor{})
	gob.Register(&HSVColor{})
}
