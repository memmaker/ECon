package common

import (
	"image/color"
	"math"
)

// Cell represents a cell in the console
var White = RGBColor{R: 1.0, G: 1.0, B: 1.0}
var Black = RGBColor{R: 0.0, G: 0.0, B: 0.0}
var Red = RGBColor{R: 1.0, G: 0.0, B: 0.0}
var Green = RGBColor{R: 0.0, G: 1.0, B: 0.0}
var Blue = RGBColor{R: 0.0, G: 0.0, B: 1.0}

// HSVColor represents a RGBA color in the console
type HSVColor struct {
	H float64 // [0, 1]
	S float64 // [0, 1]
	V float64 // [0, 1]
}

type RGBColor struct {
	R float64
	G float64
	B float64
}

func (R RGBColor) RGBA() (r, g, b, a uint32) {
	return R.ExposureToneMapping()
}
func (R RGBColor) ExposureToneMapping() (r, g, b, a uint32) {
	exposure := 1.0
	//lightness := R.Lightness()
	scale := float64(0xffff) // * math.Sqrt(lightness)
	// vec3 mapped = hdrColor / (hdrColor + vec3(1.0));
	//gamma := 2.2
	mappedR := 1.0 - math.Exp(-(R.R * exposure))
	r = uint32(mappedR * scale)
	mappedG := 1.0 - math.Exp(-(R.G * exposure))
	g = uint32(mappedG * scale)
	mappedB := 1.0 - math.Exp(-(R.B * exposure))
	b = uint32(mappedB * scale)
	a = uint32(0xffff)
	return r, g, b, a
}
func (R RGBColor) ReinhardToneMapping() (r, g, b, a uint32) {
	lightness := R.Lightness()
	scale := 0xffff * math.Sqrt(lightness)
	// vec3 mapped = hdrColor / (hdrColor + vec3(1.0));
	//gamma := 2.2
	mappedR := (R.R * scale) / (R.R + 1.0)
	//mappedR = math.Pow(mappedR, 1.0/gamma)
	r = uint32(mappedR)
	mappedG := (R.G * scale) / (R.G + 1.0)
	//mappedG = math.Pow(mappedG, 1.0/gamma)
	g = uint32(mappedG)
	mappedB := (R.B * scale) / (R.B + 1.0)
	//mappedB = math.Pow(mappedB, 1.0/gamma)
	b = uint32(mappedB)
	a = uint32(0xffff)
	return r, g, b, a
}
func (R RGBColor) LightnessScaledToneMapping() (r, g, b, a uint32) {
	lightness := R.Lightness()
	scale := 0xffff * math.Sqrt(lightness)
	r = uint32(R.R * scale)
	g = uint32(R.G * scale)
	b = uint32(R.B * scale)
	a = uint32(0xffff)
	return r, g, b, a
}

// Luma() is gamma-compressed
func (R RGBColor) Luma() float64 {
	return 0.2126*R.R + 0.7152*R.G + 0.0722*R.B
}

// Luminance() is not gamma-compressed
func (R RGBColor) Luminance() float64 {
	return 0.2126*degamma(R.R) + 0.7152*degamma(R.G) + 0.0722*degamma(R.B)
}

func (R RGBColor) Lightness() float64 {
	y := R.Luma()
	var result float64
	if y <= (216.0 / 24389) {
		result = y * (24389.0 / 27)
	} else {
		result = math.Pow(y, 1/3.0)*116.0 - 16.0
	}
	result /= 100.0
	return result
}
func degamma(channelValue float64) float64 {
	// Send this function a decimal sRGB gamma encoded color value
	// between 0.0 and 1.0, and it returns a linearized value.
	if channelValue <= 0.04045 {
		return channelValue / 12.92
	} else {
		return math.Pow((channelValue+0.055)/1.055, 2.4)
	}
}

func (R RGBColor) WithClampTo(intensity float64) RGBColor {
	return RGBColor{
		R: Clamp(R.R, 0, intensity),
		G: Clamp(R.G, 0, intensity),
		B: Clamp(R.B, 0, intensity),
	}
}

func (R RGBColor) Darken(currentLightness float64, newLightness float64) RGBColor {
	diff := currentLightness - newLightness
	return RGBColor{
		R: R.R - diff,
		G: R.G - diff,
		B: R.B - diff,
	}
}

func (R RGBColor) ToHSV() HSVColor {
	return NewHSVColorFromRGB(R.R, R.G, R.B)
}

func NewHSVColor(h, s, v float64) HSVColor {
	return HSVColor{h, s, v}
}

// NewHSVColorFromRGB creates a new color from R,G,B values
func NewHSVColorFromRGB(r, g, b float64) HSVColor {
	return NewHSVColor(RGBtoHSV(r, g, b))
}

func HSLColor(h, s, l float64) HSVColor {
	return NewHSVColor(HSLtoHSV(h, s, l))
}

func HSVtoHSL(h float64, s float64, v float64) (float64, float64, float64) {
	// both hsv and hsl values are in [0, 1]
	l := (2 - s) * v / 2
	if l != 0 {
		if l == 1 {
			s = 0
		} else if l < 0.5 {
			s = s * v / (l * 2)
		} else {
			s = s * v / (2 - l*2)
		}
	}

	return h, s, l
}

func HSLtoHSV(hslH float64, hslS float64, hslL float64) (float64, float64, float64) {
	// both hsv and hsl values are in [0, 1]
	var hsvH, hsvS, hsvV float64
	hsvH = hslH
	hsvV = hslL + hslS*math.Min(hslL, 1-hslL)
	if hsvV == 0 {
		hsvS = 0
	} else {
		hsvS = 2 * (1 - hslL/hsvV)
	}
	return hsvH, hsvS, hsvV
}
func RGBtoHSV(fR float64, fG float64, fB float64) (h, s, v float64) {
	max := math.Max(math.Max(fR, fG), fB)
	min := math.Min(math.Min(fR, fG), fB)
	d := max - min
	s, v = 0, max
	if max > 0 {
		s = d / max
	}
	if max == min {
		// Achromatic.
		h = 0
	} else {
		// Chromatic.
		switch max {
		case fR:
			h = (fG - fB) / d
			if fG < fB {
				h += 6
			}
		case fG:
			h = (fB-fR)/d + 2
		case fB:
			h = (fR-fG)/d + 4
		}
		h /= 6
	}
	return
}

// RGBA returns the color values as uint32s
func (c HSVColor) RGBA() (r, g, b, a uint32) {
	cr, cg, cb := HSVtoRGB(c.H, c.S, c.V)
	return uint32(cr * 0xFFFF), uint32(cg * 0xFFFF), uint32(cb * 0xFFFF), uint32(0xFFFF)
}

func AlphaBlend(new, curr color.Color) color.Color {
	nr, ng, nb, na := new.RGBA()
	if na == 0xFFFF {
		return new
	}
	if na == 0 {
		return curr
	}
	cr, cg, cb, ca := curr.RGBA()
	if ca == 0 {
		return new
	}

	return color.RGBA64{
		R: uint16((nr*0xFFFF + cr*(0xFFFF-na)) / 0xFFFF),
		G: uint16((ng*0xFFFF + cg*(0xFFFF-na)) / 0xFFFF),
		B: uint16((nb*0xFFFF + cb*(0xFFFF-na)) / 0xFFFF),
		A: uint16((na*0xFFFF + ca*(0xFFFF-na)) / 0xFFFF),
	}
}

// P returns a pointer to the color
func (c HSVColor) P() *HSVColor {
	return &c
}

func (c HSVColor) Lighten(scale float64) HSVColor {
	// scale must be between 0 and 1
	oldV := c.V
	interval := 1.0 - oldV
	newV := math.Min(1.0, oldV+interval*scale)
	return NewHSVColor(c.H, c.S, newV)
}

func (c HSVColor) WithV(value float64) HSVColor {
	return NewHSVColor(c.H, c.S, Clamp(value, 0, 1))
}

func (c HSVColor) ToRGBColor() RGBColor {
	r, g, b := HSVtoRGB(c.H, c.S, c.V)
	return RGBColor{r, g, b}
}

func HSVtoRGB(h float64, s float64, v float64) (float64, float64, float64) {
	hThreeSixty := h * 360.0
	Hp := hThreeSixty / 60.0
	c := v * s
	x := c * (1.0 - math.Abs(math.Mod(Hp, 2.0)-1.0))

	m := v - c
	r, g, b := 0.0, 0.0, 0.0

	switch {
	case 0.0 <= Hp && Hp < 1.0:
		r = c
		g = x
	case 1.0 <= Hp && Hp < 2.0:
		r = x
		g = c
	case 2.0 <= Hp && Hp < 3.0:
		g = c
		b = x
	case 3.0 <= Hp && Hp < 4.0:
		g = x
		b = c
	case 4.0 <= Hp && Hp < 5.0:
		r = x
		b = c
	case 5.0 <= Hp && Hp < 6.0:
		r = c
		b = x
	}
	return m + r, m + g, m + b
}

func (c HSVColor) BlendRGB(rgbColor RGBColor, value float64) HSVColor {

	solidR, solidG, solidB := HSVtoRGB(c.H, c.S, c.V)
	lightR := rgbColor.R
	lightG := rgbColor.G
	lightB := rgbColor.B

	mixedR := Clamp(lightR+solidR, 0, 1)
	mixedG := Clamp(lightG+solidG, 0, 1)
	mixedB := Clamp(lightB+solidB, 0, 1)

	return NewHSVColorFromRGB(mixedR, mixedG, mixedB)
}

func NewHSVColorFromRGBA(r float64, g float64, b float64, a float64) HSVColor {
	h, s, v := RGBAtoHSV(r, g, b, a)
	return HSVColor{h, s, v}
}

func RGBAtoHSV(r float64, g float64, b float64, a float64) (float64, float64, float64) {
	max := math.Max(math.Max(r, g), b)
	min := math.Min(math.Min(r, g), b)
	delta := max - min
	h := 0.0
	s := 0.0
	v := max

	if max != 0 {
		s = delta / max
	}

	if s != 0 {
		if r == max {
			h = (g - b) / delta
		} else if g == max {
			h = 2 + (b-r)/delta
		} else {
			h = 4 + (r-g)/delta
		}
		h *= 60
		if h < 0 {
			h += 360
		}
	}

	return h, s, v
}

func (c HSVColor) WithH(h float64) HSVColor {
	return NewHSVColor(h, c.S, c.V)
}

func (c HSVColor) LerpH(h float64, ratio float64) HSVColor {
	return NewHSVColor(c.H+(h-c.H)*ratio, c.S, c.V)
}

func (c HSVColor) WithS(saturation float64) HSVColor {
	return NewHSVColor(c.H, saturation, c.V)
}

func Clamp(f float64, min float64, max float64) float64 {
	if f < min {
		return min
	}
	if f > max {
		return max
	}
	return f
}

type Style struct {
	Fg HSVColor // foreground color
	Bg HSVColor // background color
}

// WithFg returns a derived style with a new foreground color.
func (st Style) WithFg(cl HSVColor) Style {
	st.Fg = cl
	return st
}

// WithBg returns a derived style with a new background color.
func (st Style) WithBg(cl HSVColor) Style {
	st.Bg = cl
	return st
}
