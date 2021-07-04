package runtime

// ConnectAll spawns a goroutine for every receiver in every connection.
func ConnectAll(cc []Conn) {
	for i := range cc {
		go connect(cc[i])
	}
}

func connect(c Conn) {
	for msg := range c.sender {
		for i := range c.receivers {
			r := c.receivers[i]
			go func() { r <- msg }()
		}
	}
}
