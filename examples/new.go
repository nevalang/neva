package new

type subscription struct {
	from <-chan interface{}
	to   []chan<- interface{}
}

func subscribe(s subscription) {
	for v := range s.from {
		for i := range s.to {
			c := s.to[i]
			go func() { c <- v }()
		}
	}
}

func subscribeAll(ss []subscription) {
	for i := range ss {
		go subscribe(ss[i])
	}
}
