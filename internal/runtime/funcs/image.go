package funcs

import (
	"image"
	"image/color"

	"github.com/nevalang/neva/internal/runtime/messages"
)

// TODO can't we use uint8 here?
//
//nolint:recvcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
type rgbaMsg struct {
	r int64
	g int64
	b int64
	a int64
}

func (c *rgbaMsg) decode(msg messages.Msg) {
	m := msg.Struct()
	c.r = messages.StructGetField(m, "r").Int()
	c.g = messages.StructGetField(m, "g").Int()
	c.b = messages.StructGetField(m, "b").Int()
	c.a = messages.StructGetField(m, "a").Int()
}

func (c rgbaMsg) color() color.Color {
	return color.RGBA64{
		R: clampUint16(c.r),
		G: clampUint16(c.g),
		B: clampUint16(c.b),
		A: clampUint16(c.a),
	}
}

type pixelMsg struct {
	x     int64
	y     int64
	color rgbaMsg
}

func (p *pixelMsg) decode(msg messages.Msg) {
	m := msg.Struct()
	p.x = messages.StructGetField(m, "x").Int()
	p.y = messages.StructGetField(m, "y").Int()
	p.color.decode(messages.StructGetField(m, "color"))
}

type imageMsg struct {
	pixels []byte
	width  int64
	height int64
}

func (i imageMsg) createImage() image.Image {
	// Use pixels directly if available.
	pix := i.pixels
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

func clampUint16(value int64) uint16 {
	if value < 0 {
		return 0
	}
	const maxUint16 = int64(^uint16(0))
	if value > maxUint16 {
		return ^uint16(0)
	}
	return uint16(value)
}
