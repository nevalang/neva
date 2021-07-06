package std

// type StdComponent struct {
// }

// func NewPlus() Plus {
// 	in := make([]chan int)
// 	return Plus{
// 		inports: in,
// 		outports: out,
// 	}
// }

// type Plus struct {
// 	inports Ports
// 	outports Ports
// }

// func (p Plus) Ports() struct{ In, Out Ports } {
// 	return struct{ In, Out Ports }{}
// }

// func Array(cc []<-chan interface{}) []interface{} {
// 	ch := make(chan interface{})

// 	for i := range cc {
// 		c := cc[i]
// 		go func() {
// 			ch <- <-c
// 		}()
// 	}

// 	var res []interface{}
// 	for v := range ch {
// 		res = append(res, v)
// 	}
// 	close(ch)

// 	return res
// }
