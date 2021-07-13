package runtime

func ConnectAll(cc []ChanRel) {
	for i := range cc {
		go connect(cc[i])
	}
}

func connect(c ChanRel) {
	for msg := range c.Sender {
		for i := range c.Receivers {
			r := c.Receivers[i]
			go func() { r <- msg }()
		}
	}
}
