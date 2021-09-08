package runtime

import "fmt"

type Connector interface {
	Connect([]connection)
	ConnectOperator(string, IO) error
}

type connector struct {
	ops       map[string]Operator
	onSend    func(msg Msg, from PortAddr)
	onReceive func(msg Msg, from, to PortAddr)
}

func (c connector) Connect(pp []connection) {
	for i := range pp {
		go c.connectPair(pp[i])
	}
}

func (c connector) connectPair(con connection) {
	for msg := range con.from.ch {
		c.onSend(msg, con.from.addr)

		for _, recv := range con.to {
			select {
			case recv.ch <- msg:
				c.onReceive(msg, con.from.addr, recv.addr)
			default:
				go func(to Port, m Msg) {
					to.ch <- m
					c.onReceive(m, con.from.addr, to.addr)
				}(recv, msg)
			}
		}
	}
}

func (c connector) ConnectOperator(name string, io IO) error {
	op, ok := c.ops[name]
	if !ok {
		return fmt.Errorf("ErrUnknownOperator: %s", name)
	}

	if err := op(io); err != nil {
		return err
	}

	return nil
}
