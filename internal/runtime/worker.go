package runtime

type Worker struct {
	relations []ChanRelation
}

func (m Worker) Run() {
	for i := range m.relations {
		go m.connect(m.relations[i])
	}
}

func (m Worker) connect(c ChanRelation) {
	for msg := range c.Sender {
		for i := range c.Receivers {
			r := c.Receivers[i]
			go func() { r <- msg }()
		}
	}
}
