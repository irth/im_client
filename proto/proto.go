package proto

type Message struct {
	Type             string
	InitMessage      *InitMessage      `json:",omitempty"`
	SubscribeMessage *SubscribeMessage `json:",omitempty"`
}

type Messageable interface {
	ToMessage() Message
}

type InitMessage struct {
	Name   string
	Secret string
}

func (i *InitMessage) ToMessage() Message {
	return Message{
		Type:        "InitMessage",
		InitMessage: i,
	}
}

type InitMessageResponse struct {
	Success bool
}

type Event string

type SubscribeMessage struct {
	Events []Event
}

func (s *SubscribeMessage) ToMessage() Message {
	return Message{
		Type:             "SubscribeMessage",
		SubscribeMessage: s,
	}
}
