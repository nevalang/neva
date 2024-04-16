package funcs

import (
	"image"
	"image/color"

	"github.com/nevalang/neva/internal/runtime"
)

type rgbaMsg struct {
	r runtime.IntMsg
	g runtime.IntMsg
	b runtime.IntMsg
	a runtime.IntMsg
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
	return c.r.Decode(m["r"]) &&
		c.r.Decode(m["g"]) &&
		c.r.Decode(m["b"]) &&
		c.r.Decode(m["a"])
}
func (c *rgbaMsg) encode() runtime.Msg {
	return runtime.NewMapMsg(map[string]runtime.Msg{
		"r": c.r,
		"g": c.g,
		"b": c.b,
		"a": c.a,
	})
}
func (c *rgbaMsg) color() color.Color {
	return color.RGBA64{R: uint16(c.r.Int()), G: uint16(c.g.Int()), B: uint16(c.b.Int()), A: uint16(c.a.Int())}
}
func (c *rgbaMsg) decodeColor(clr color.Color) {
	r, g, b, a := clr.RGBA()
	c.r = runtime.NewIntMsg(int64(r))
	c.g = runtime.NewIntMsg(int64(g))
	c.b = runtime.NewIntMsg(int64(b))
	c.a = runtime.NewIntMsg(int64(a))
}

type pointMsg struct {
	x int64
	y int64
}

func (p *pointMsg) reset() { *p = pointMsg{} }
func (p *pointMsg) decode(msg runtime.Msg) bool {
	p.reset()
	if msg == nil {
		return false
	}
	m := msg.Map()
	if m == nil {
		return false
	}
	p.x = m["x"].Int()
	p.y = m["y"].Int()
	return true
}
func (p *pointMsg) encode() runtime.Msg {
	return runtime.NewMapMsg(map[string]runtime.Msg{
		"x": runtime.NewIntMsg(p.x),
		"y": runtime.NewIntMsg(p.y),
	})
}
func (p *pointMsg) decodePoint(pt image.Point) {
	p.x = int64(pt.X)
	p.y = int64(pt.Y)
}

type pixelMsg struct {
	point pointMsg
	color rgbaMsg
}

func (p *pixelMsg) reset() {}
func (p *pixelMsg) decode(msg runtime.Msg) bool {
	if msg == nil {
		return false
	}
	m := msg.Map()
	if m == nil {
		return false
	}
	return p.point.decode(m["point"]) && p.color.decode(m["color"])
}
func (p *pixelMsg) encode() runtime.Msg {
	return runtime.NewMapMsg(map[string]runtime.Msg{
		"point": p.point.encode(),
		"color": p.color.encode(),
	})
}

type boundsMsg struct {
	min pointMsg
	max pointMsg
}

func (b *boundsMsg) reset() {}
func (b *boundsMsg) decode(msg runtime.Msg) bool {
	if msg == nil {
		return false
	}
	m := msg.Map()
	if m == nil {
		return false
	}
	return b.min.decode(m["min"]) && b.max.decode(m["max"])
}
func (b *boundsMsg) rect() image.Rectangle {
	return image.Rect(int(b.min.x), int(b.min.y), int(b.max.x), int(b.max.y))
}
func (b *boundsMsg) decodeRect(r image.Rectangle) {
	b.min.decodePoint(r.Min)
	b.max.decodePoint(r.Max)
}
func (b *boundsMsg) encode() runtime.Msg {
	return runtime.NewMapMsg(map[string]runtime.Msg{
		"min": b.min.encode(),
		"max": b.max.encode(),
	})
}

type formatMsg int64

func (f *formatMsg) reset() { *f = 0 }
func (f *formatMsg) decode(msg runtime.Msg) bool {
	f.reset()
	if msg == nil {
		return false
	}
	*f = formatMsg(msg.Int())
	return *f >= 0 && *f <= 2
}
