package runtime

// Connection represents one-to-many relation between sender and recievers.
type Connection struct {
	sender    <-chan Message   // outport
	recievers []chan<- Message // inports
}

// ConnectAll spawns routine for every reciever in every connection.
func ConnectAll(cc []Connection) {
	for i := range cc {
		go connect(cc[i])
	}
}

func connect(c Connection) {
	for msg := range c.sender {
		for i := range c.recievers {
			r := c.recievers[i]
			go func() { r <- msg }()
		}
	}
}
