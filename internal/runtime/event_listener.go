package runtime

import "fmt"

type EmptyListener struct{}

func (l EmptyListener) Send(event Event, msg Msg) Msg {
	// if event.Type != MessageReceivedEvent {
	// 	return msg
	// }
	// fmt.Println(event, msg)
	return msg
}

type Event struct {
	Type            EventType
	MessageSent     *EventMessageSent
	MessagePending  *EventMessagePending
	MessageReceived *EventMessageReceived
}

func (e Event) String() string {
	var s string
	switch e.Type {
	case MessageSentEvent:
		s = e.MessageSent.String()
	case MessagePendingEvent:
		s = e.MessagePending.String()
	case MessageReceivedEvent:
		s = e.MessageReceived.String()
	}

	return fmt.Sprintf("%v: %v", e.Type.String(), s)
}

type EventMessageSent struct {
	SenderPortAddr    PortAddr
	ReceiverPortAddrs map[PortAddr]struct{} // We use map to work with breakpoints
}

func (e EventMessageSent) String() string {
	rr := "{ "
	for r := range e.ReceiverPortAddrs {
		rr += r.String() + ", "
	}
	rr += "}"
	return fmt.Sprintf("%v -> %v", e.SenderPortAddr, rr)
}

type EventMessagePending struct {
	Meta             ConnectionMeta // We can use sender from here and receivers just as a handy metadata
	ReceiverPortAddr PortAddr       // So what we really need is sender and receiver port addrs
}

func (e EventMessagePending) String() string {
	return fmt.Sprintf("%v -> %v", e.Meta.SenderPortAddr, e.ReceiverPortAddr)
}

type EventMessageReceived struct {
	Meta             ConnectionMeta // Same as with pending event
	ReceiverPortAddr PortAddr
}

func (e EventMessageReceived) String() string {
	return fmt.Sprintf("%v -> %v", e.Meta.SenderPortAddr, e.ReceiverPortAddr)
}

type EventType uint8

const (
	MessageSentEvent     EventType = 1 // Message is sent from sender to its receivers
	MessagePendingEvent  EventType = 2 // Message has reached receiver but not yet passed inside
	MessageReceivedEvent EventType = 3 // Message is passed inside receiver
)

func (e EventType) String() string {
	switch e {
	case MessageSentEvent:
		return "Message sent"
	case MessagePendingEvent:
		return "Message pending"
	case MessageReceivedEvent:
		return "Message received"
	}
	return "Unknown Event Type"
}

type EventListener interface {
	Send(Event, Msg) Msg
}
