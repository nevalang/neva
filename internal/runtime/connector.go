package runtime

type Connector interface {
	Connect()
}

type connector struct {
	pp        []pair
	onSend    func(msg Msg, from PortAddr)
	onReceive func(msg Msg, to PortAddr)
}

func (c connector) Connect() {
	for i := range c.pp {
		go c.connectPair(c.pp[i])
	}
}

func (c connector) connectPair(p pair) {
	for msg := range p.from.ch {
		c.onSend(msg, p.from.addr)

		for _, recv := range p.to {
			select {
			case recv.ch <- msg:
				c.onReceive(msg, recv.addr)
			default:
				go func(to Port, m Msg) {
					to.ch <- m
					c.onReceive(m, to.addr)
				}(recv, msg)
			}
		}
	}
}
