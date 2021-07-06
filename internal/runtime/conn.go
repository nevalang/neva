package runtime

// ConnectAll spawns a goroutine for every receiver in every connection.
func ConnectAll(cc []Conn) {
	for i := range cc {
		go connect(cc[i])
	}
}

func connect(c Conn) {
	for msg := range c.Sender {
		for i := range c.Receivers {
			r := c.Receivers[i]
			go func() { r <- msg }()
		}
	}
}
