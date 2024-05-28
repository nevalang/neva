package funcs

import (
	"image"
	"image/color"

	"github.com/nevalang/neva/internal/runtime"
)

type rgbaMsg struct {
	r int64
	g int64
	b int64
	a int64
}

func (c *rgbaMsg) decode(msg runtime.Msg) {
	m := msg.Map()
	c.r = m["r"].Int()
	c.g = m["g"].Int()
	c.b = m["b"].Int()
	c.a = m["a"].Int()
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
	m := msg.Map()
	p.x = m["x"].Int()
	p.y = m["y"].Int()
	p.color.decode(m["color"])
}

type imageMsg struct {
	pixels string
	width  int64
	height int64
}

func (i *imageMsg) reset() { *i = imageMsg{} }
func (i *imageMsg) decode(msg runtime.Msg) {
	i.reset()
	m := msg.Map()
	i.pixels = m["pixels"].Str()
	i.width = m["width"].Int()
	i.height = m["height"].Int()
}
func (i imageMsg) encode() runtime.Msg {
	return runtime.NewMapMsg(map[string]runtime.Msg{
		"pixels": runtime.NewStrMsg(i.pixels),
		"width":  runtime.NewIntMsg(i.width),
		"height": runtime.NewIntMsg(i.height),
	})
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
func (i *imageMsg) decodeImage(img *image.RGBA) {
	i.pixels = string(img.Pix)
	i.width = int64(img.Rect.Dx())
	i.height = int64(img.Rect.Dy())
}

type pixelStreamMsg struct {
	idx int64
	pixelMsg
	last bool
}

func (i *pixelStreamMsg) decode(msg runtime.Msg) {
	m := msg.Map()
	i.idx = m["idx"].Int()
	i.pixelMsg.decode(m["data"])
	i.last = m["last"].Bool()
}
