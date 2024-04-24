package funcs

import (
	"image"
	"image/color"

	"github.com/nevalang/neva/internal/runtime"
)

func decodeInt(i *int64, m runtime.Msg) bool {
	if m.Type() != runtime.IntMsgType {
		return false
	}
	*i = m.Int()
	return true
}

func decodeStr(s *string, m runtime.Msg) bool {
	if m.Type() != runtime.StrMsgType {
		return false
	}
	*s = m.Str()
	return true
}

type rgbaMsg struct {
	r int64
	g int64
	b int64
	a int64
}

func (c *rgbaMsg) reset() { *c = rgbaMsg{} }
func (c *rgbaMsg) decode(msg runtime.Msg) bool {
	c.reset()
	if msg == nil {
		return false
	}
	m := msg.Map()
	if m == nil {
		return false
	}
	return decodeInt(&c.r, m["r"]) &&
		decodeInt(&c.g, m["g"]) &&
		decodeInt(&c.b, m["b"]) &&
		decodeInt(&c.a, m["a"])
}
func (c rgbaMsg) encode() runtime.Msg {
	return runtime.NewMapMsg(map[string]runtime.Msg{
		"r": runtime.NewIntMsg(c.r),
		"g": runtime.NewIntMsg(c.g),
		"b": runtime.NewIntMsg(c.b),
		"a": runtime.NewIntMsg(c.a),
	})
}
func (c rgbaMsg) color() color.Color {
	return color.RGBA64{R: uint16(c.r), G: uint16(c.g), B: uint16(c.b), A: uint16(c.a)}
}
func (c *rgbaMsg) decodeColor(clr color.Color) {
	r, g, b, a := clr.RGBA()
	c.r = int64(r)
	c.g = int64(g)
	c.b = int64(b)
	c.a = int64(a)
}

type pixelMsg struct {
	x     int64
	y     int64
	color rgbaMsg
}

func (p *pixelMsg) reset() { *p = pixelMsg{} }
func (p *pixelMsg) decode(msg runtime.Msg) bool {
	if msg == nil {
		return false
	}
	m := msg.Map()
	if m == nil {
		return false
	}
	return decodeInt(&p.x, m["x"]) && decodeInt(&p.y, m["y"]) && p.color.decode(m["color"])
}
func (p pixelMsg) encode() runtime.Msg {
	return runtime.NewMapMsg(map[string]runtime.Msg{
		"x":     runtime.NewIntMsg(p.x),
		"y":     runtime.NewIntMsg(p.y),
		"color": p.color.encode(),
	})
}

type imageMsg struct {
	pixels string
	width  int64
	height int64
}

func (i *imageMsg) reset() { *i = imageMsg{} }
func (i *imageMsg) decode(msg runtime.Msg) bool {
	if i.reset(); msg.Type() != runtime.MapMsgType {
		return false
	}
	m := msg.Map()
	return decodeStr(&i.pixels, m["pixels"]) &&
		decodeInt(&i.width, m["width"]) &&
		decodeInt(&i.height, m["height"])
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
			// One byte for each color component.
			pix = make([]uint8, 4*size)
		}
	}
	im := &image.RGBA{
		Stride: 1,
		Pix:    pix,
		Rect:   image.Rect(0, 0, int(i.width), int(i.height)),
	}
	return im
}
func (i *imageMsg) decodeImage(img *image.RGBA) bool {
	i.pixels = string(img.Pix)
	i.width = int64(img.Rect.Dx())
	i.height = int64(img.Rect.Dy())
	return true
}
