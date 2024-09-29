package funcs

import (
	"image"
	"image/color"

	"github.com/nevalang/neva/internal/runtime"
)

// TODO can't we use uint8 here?
type rgbaMsg struct {
	r int64
	g int64
	b int64
	a int64
}

func (c *rgbaMsg) decode(msg runtime.Msg) {
	m := msg.Struct()
	c.r = m.Get("r").Int()
	c.g = m.Get("g").Int()
	c.b = m.Get("b").Int()
	c.a = m.Get("a").Int()
}

func (c rgbaMsg) color() color.Color {
	return color.RGBA64{R: uint16(c.r), G: uint16(c.g), B: uint16(c.b), A: uint16(c.a)}
}

type pixelMsg struct {
	x     int64
	y     int64
	color rgbaMsg
}

func (p *pixelMsg) decode(msg runtime.Msg) {
	m := msg.Struct()
	p.x = m.Get("x").Int()
	p.y = m.Get("y").Int()
	p.color.decode(m.Get("color"))
}

type imageMsg struct {
	pixels string
	width  int64
	height int64
}

func (i imageMsg) createImage() image.Image {
	// Use pixels directly if available.
	pix := []uint8(i.pixels)
	if len(pix) == 0 {
		if size := i.width * i.height; size > 0 {
			// Allocate new pixels.
			// One byte for each color flow.
			pix = make([]uint8, 4*size)
		}
	}
	im := &image.RGBA{
		Stride: int(i.width),
		Pix:    pix,
		Rect:   image.Rect(0, 0, int(i.width), int(i.height)),
	}
	return im
}

type pixelStreamMsg struct {
	idx int64
	pixelMsg
	last bool
}

func (i *pixelStreamMsg) decode(msg runtime.Msg) {
	m := msg.Struct()
	i.idx = m.Get("idx").Int()
	i.pixelMsg.decode(m.Get("data"))
	i.last = m.Get("last").Bool()
}
