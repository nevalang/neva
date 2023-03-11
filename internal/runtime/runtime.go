package runtime

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

/* --- PROGRAM --- */

type (
	Program struct {
		StartPortAddr PortAddr
		Ports         Ports
		Connections   []Connection
		Routines      Routines
	}

	PortAddr struct {
		Path, Name string
		Idx        uint8
	}

	Ports map[PortAddr]chan Msg

	Connection struct {
		Sender    ConnectionSide
		Receivers []ConnectionSide
	}

	ConnectionSide struct {
		Port chan Msg
		Meta ConnectionSideMeta
	}

	ConnectionSideMeta struct {
		PortAddr  PortAddr
		Selectors []Selector
	}

	Selector struct {
		RecField string // "" means use ArrIdx
		ArrIdx   int
	}

	Routines struct {
		Giver     []GiverRoutine
		Component []ComponentRoutine
	}

	GiverRoutine struct {
		OutPort chan Msg
		Msg     Msg
	}

	ComponentRoutine struct {
		Ref ComponentRef
		IO  IO
	}

	ComponentRef struct {
		Pkg, Name string
	}

	IO struct {
		In, Out IOPorts
	}
)

/* --- PORTS --- */

type IOPorts map[string][]chan Msg

var (
	ErrSinglePortCount = errors.New("number of ports found by name not equals to one")
	ErrArrPortNotFound = errors.New("number of ports found by name equals to zero")
)

func (i IOPorts) Port(name string) (chan Msg, error) {
	if len(i[name]) != 1 {
		return nil, fmt.Errorf("%w: %v", ErrSinglePortCount, len(i[name]))
	}
	return i[name][0], nil
}

func (i IOPorts) ArrPort(name string) ([]chan Msg, error) {
	if len(i[name]) == 0 {
		return nil, ErrArrPortNotFound
	}
	return i[name], nil
}

/* --- MESSAGES --- */

type Msg interface {
	fmt.Stringer
	Type() Type
	Bool() bool
	Int() int64
	Float() float64
	Str() string
	Vec() []Msg
	Map() map[string]Msg
}

type Type uint8

const (
	BoolMsgType Type = iota
	IntMsgType
	FloatMsgType
	StrMsgType
	VecMsgType
	MapMsgType
)

// Empty

type emptyMsg struct{}

func (emptyMsg) Bool() bool          { return false }
func (emptyMsg) Int() int64          { return 0 }
func (emptyMsg) Float() float64      { return 0 }
func (emptyMsg) Str() string         { return "" }
func (emptyMsg) Vec() []Msg          { return []Msg{} }
func (emptyMsg) Map() map[string]Msg { return map[string]Msg{} }

// Int

type IntMsg struct {
	emptyMsg
	v int64
}

func (msg IntMsg) Int() int64     { return msg.v }
func (msg IntMsg) Type() Type     { return IntMsgType }
func (msg IntMsg) String() string { return strconv.Itoa(int(msg.v)) }

func NewIntMsg(n int64) IntMsg {
	return IntMsg{
		emptyMsg: emptyMsg{},
		v:        n,
	}
}

// Str

type StrMsg struct {
	emptyMsg
	v string
}

func (msg StrMsg) Str() string    { return msg.v }
func (msg StrMsg) Type() Type     { return StrMsgType }
func (msg StrMsg) String() string { return strconv.Quote(msg.v) }

func NewStrMsg(s string) StrMsg {
	return StrMsg{
		emptyMsg: emptyMsg{},
		v:        s,
	}
}

// Bool

type BoolMsg struct {
	emptyMsg
	v bool
}

func (msg BoolMsg) Bool() bool     { return msg.v }
func (msg BoolMsg) Type() Type     { return BoolMsgType }
func (msg BoolMsg) String() string { return fmt.Sprint(msg.v) }

func NewBoolMsg(b bool) BoolMsg {
	return BoolMsg{
		emptyMsg: emptyMsg{},
		v:        b,
	}
}

/* --- Map --- */

type MapMsg struct {
	emptyMsg
	v map[string]Msg
}

func (msg MapMsg) Map() map[string]Msg { return msg.v }
func (msg MapMsg) Type() Type          { return MapMsgType }
func (msg MapMsg) String() string {
	b := &strings.Builder{}
	b.WriteString("{")
	c := 0
	for k, el := range msg.v {
		c++
		if c < len(msg.v) {
			fmt.Fprintf(b, " %s: %s, ", k, el.String())
			continue
		}
		fmt.Fprintf(b, "%s: %s ", k, el.String())
	}
	b.WriteString("}")
	return b.String()
}

func NewMapMsg(v map[string]Msg) MapMsg {
	return MapMsg{
		emptyMsg: emptyMsg{},
		v:        v,
	}
}

// Vec

type VecMsg struct {
	emptyMsg
	v []Msg
}

func (msg VecMsg) Vec() []Msg { return msg.v }
func (msg VecMsg) Type() Type { return VecMsgType }
func (msg VecMsg) String() string {
	b := &strings.Builder{}
	b.WriteString("[")
	c := 0
	for _, el := range msg.v {
		c++
		if c < len(msg.v) {
			fmt.Fprintf(b, "%s, ", el.String())
			continue
		}
		fmt.Fprint(b, el.String())
	}
	b.WriteString("]")
	return b.String()
}

func NewVecMsg(v []Msg) VecMsg {
	return VecMsg{
		emptyMsg: emptyMsg{},
		v:        v,
	}
}

// Eq (NotEq, Less, Greater, Max, Min?)

func Eq(a, b Msg) bool {
	if a.Type() != b.Type() {
		return false
	}
	switch a.Type() {
	case BoolMsgType:
		return a.Bool() == b.Bool()
	case IntMsgType:
		return a.Int() == b.Int()
	case StrMsgType:
		return a.Str() == b.Str()
	case VecMsgType:
		l1 := a.Vec()
		l2 := b.Vec()
		if len(l1) != len(l2) {
			return false
		}
		for i := range l1 {
			if !Eq(l1[i], l2[i]) {
				return false
			}
		}
	case MapMsgType:
		s1 := a.Map()
		s2 := a.Map()
		if len(s1) != len(s2) {
			return false
		}
		for k := range s1 {
			if !Eq(s1[k], s2[k]) {
				return false
			}
		}
	}
	return false
}

/* --- RUNTIME --- */

type Runtime struct {
	connector     Connector
	routineRunner RoutineRunner
}

func NewRuntime(
	connector Connector,
	routineRunner RoutineRunner,
) Runtime {
	return Runtime{
		connector:     connector,
		routineRunner: routineRunner,
	}
}

type (
	Connector interface {
		Connect(context.Context, []Connection) error
	}
	RoutineRunner interface {
		Run(context.Context, Routines) error
	}
)

var (
	ErrStartPortNotFound = errors.New("start port not found")
	ErrConnector         = errors.New("connector")
	ErrRoutineRunner     = errors.New("routine runner")
)

func (r Runtime) Run(ctx context.Context, prog Program) error {
	startPort, ok := prog.Ports[prog.StartPortAddr]
	if !ok {
		return fmt.Errorf("%w: %v", ErrStartPortNotFound, prog.StartPortAddr)
	}

	g, gctx := WithContext(ctx)
	g.Go(func() error {
		if err := r.connector.Connect(gctx, prog.Connections); err != nil {
			return fmt.Errorf("%w: %v", ErrConnector, err)
		}
		return nil
	})
	g.Go(func() error {
		if err := r.routineRunner.Run(gctx, prog.Routines); err != nil {
			return fmt.Errorf("%w: %v", ErrRoutineRunner, err)
		}
		return nil
	})

	startPort <- nil
	return g.Wait()
}

/*  --- ROUTINE-RUNNER --- */

type RoutineRunnerImlp struct {
	giver     GiverRunner
	component ComponentRunner
}

func NewRoutineRunner(giver GiverRunner, component ComponentRunner) RoutineRunner {
	return RoutineRunnerImlp{
		giver:     giver,
		component: component,
	}
}

type (
	GiverRunner interface {
		Run(context.Context, []GiverRoutine) error
	}
	ComponentRunner interface {
		Run(context.Context, []ComponentRoutine) error
	}
)

var (
	ErrComponent = errors.New("component")
	ErrGiver     = errors.New("giver")
)

func (e RoutineRunnerImlp) Run(ctx context.Context, routines Routines) error {
	g, gctx := WithContext(ctx)

	g.Go(func() error {
		if err := e.giver.Run(gctx, routines.Giver); err != nil {
			return errors.Join(ErrGiver, err)
		}
		return nil
	})

	g.Go(func() error {
		if err := e.component.Run(gctx, routines.Component); err != nil {
			return errors.Join(ErrComponent, err)
		}
		return nil
	})

	return g.Wait()
}

/* --- GIVER-RUNNER --- */

type GiverRunnerImlp struct{}

func (e GiverRunnerImlp) Run(ctx context.Context, givers []GiverRoutine) error {
	wg := sync.WaitGroup{}
	wg.Add(len(givers))

	for i := range givers {
		giver := givers[i]
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					giver.OutPort <- giver.Msg
				}
			}
		}()
	}

	wg.Wait()

	return ctx.Err()
}

/* --- COMPONENT-RUNNER --- */

var (
	ErrRepo          = errors.New("repo")
	ErrComponentFunc = errors.New("operator func")
)

type ComponentRunnerImpl struct {
	repo map[ComponentRef]func(context.Context, IO) error
}

func NewComponentRunner(repo map[ComponentRef]func(context.Context, IO) error) ComponentRunnerImpl {
	return ComponentRunnerImpl{
		repo: repo,
	}
}

func (c ComponentRunnerImpl) Run(ctx context.Context, components []ComponentRoutine) error {
	g, gctx := WithContext(ctx)

	for i := range components {
		component := components[i]

		f, ok := c.repo[component.Ref]
		if !ok {
			return fmt.Errorf("%w: %v", ErrRepo, component.Ref)
		}

		g.Go(func() error {
			if err := f(gctx, component.IO); err != nil {
				return fmt.Errorf("%w: %v", errors.Join(ErrComponentFunc, err), component.Ref)
			}
			return nil
		})
	}

	return g.Wait()
}

/* --- CONNECTOR --- */

var (
	ErrBroadcast         = errors.New("broadcast")
	ErrDistribute        = errors.New("distribute")
	ErrSelectorSending   = errors.New("selector after sending")
	ErrSelectorReceiving = errors.New("selector before receiving")
)

type ConnectorImlp struct {
	interceptor Interceptor
}

func NewConnector(interceptor Interceptor) Connector {
	return ConnectorImlp{
		interceptor: interceptor,
	}
}

type Interceptor interface {
	AfterSending(from ConnectionSideMeta, msg Msg) Msg
	BeforeReceiving(from, to ConnectionSideMeta, msg Msg) Msg
	AfterReceiving(from, to ConnectionSideMeta, msg Msg)
}

func (c ConnectorImlp) Connect(ctx context.Context, net []Connection) error {
	g, gctx := WithContext(ctx)

	for i := range net {
		conn := net[i]

		g.Go(func() error {
			if err := c.broadcast(gctx, conn); err != nil {
				return fmt.Errorf("%w: %v", errors.Join(ErrBroadcast, err), conn)
			}
			return nil
		})
	}

	return g.Wait()
}

func (c ConnectorImlp) broadcast(ctx context.Context, conn Connection) error {
	var err error
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-conn.Sender.Port:
			msg, err = c.applySelector(msg, conn.Sender.Meta.Selectors)
			if err != nil {
				return fmt.Errorf("%w: %v: %v", errors.Join(ErrSelectorSending, err), conn.Sender.Meta, msg)
			}

			msg = c.interceptor.AfterSending(conn.Sender.Meta, msg)

			if err := c.distribute(ctx, msg, conn.Sender.Meta, conn.Receivers); err != nil {
				return fmt.Errorf("%w: %v", errors.Join(ErrDistribute, err), msg)
			}
		}
	}
}

func (c ConnectorImlp) distribute(
	ctx context.Context,
	msg Msg,
	senderMeta ConnectionSideMeta,
	q []ConnectionSide,
) error {
	i := 0
	processedMessages := make(map[PortAddr]Msg, len(q)) // intercepted and selected

	for len(q) > 0 {
		recv := q[i]

		if _, ok := processedMessages[recv.Meta.PortAddr]; !ok {
			msg4 := c.interceptor.BeforeReceiving(senderMeta, recv.Meta, msg)
			msg5, err := c.applySelector(msg4, recv.Meta.Selectors)
			if err != nil {
				return fmt.Errorf("%w: %v", errors.Join(ErrSelectorReceiving, err), recv.Meta)
			}
			processedMessages[recv.Meta.PortAddr] = msg5
		}

		msg5 := processedMessages[recv.Meta.PortAddr]

		select {
		case <-ctx.Done():
			return ctx.Err()
		case recv.Port <- msg5:
			c.interceptor.AfterReceiving(senderMeta, recv.Meta, msg5)
			q = append(q[:i], q[i+1:]...) // remove cur from q
		default: // cur is too busy to receive, it's buf is full
			if i < len(q) {
				i++ // so let's go ask next while it's busy and then return
			}
		}

		if i == len(q) { // end of q, last el was processed (maybe it was busy)
			i = 0 // start over
		}
	}

	return nil
}

func (c ConnectorImlp) applySelector(msg Msg, selectors []Selector) (Msg, error) {
	if len(selectors) == 0 {
		return msg, nil
	}

	first := selectors[0]
	if first.RecField != "" {
		// msg = msg.Rec()[first.RecField]
	} else {
		// msg = msg.Arr()[first.ArrIdx]
	}

	return c.applySelector(
		msg,
		selectors[1:],
	)
}

/* ---  INTERCEPTOR ---*/

type InterceptorImlp struct{}

func (i InterceptorImlp) AfterSending(from ConnectionSideMeta, msg Msg) Msg {
	fmt.Println("after sending", from, msg)
	return msg
}
func (i InterceptorImlp) BeforeReceiving(from, to ConnectionSideMeta, msg Msg) Msg {
	fmt.Println("before receiving", from, to, msg)
	return msg
}
func (i InterceptorImlp) AfterReceiving(from, to ConnectionSideMeta, msg Msg) {
	fmt.Println("after receiving", from, to, msg)
}

/* --- ERRGROUP copy of golang.org/x/sync/errgroup --- */

type token struct{}

type Group struct {
	cancel  func()
	wg      sync.WaitGroup
	sem     chan token
	errOnce sync.Once
	err     error
}

func (g *Group) done() {
	if g.sem != nil {
		<-g.sem
	}
	g.wg.Done()
}

func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{cancel: cancel}, ctx
}

func (g *Group) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}

func (g *Group) Go(f func() error) {
	if g.sem != nil {
		g.sem <- token{}
	}
	g.wg.Add(1)
	go func() {
		defer g.done()
		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	}()
}
