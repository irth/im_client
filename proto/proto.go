package proto

type InitMessage struct {
	Name   string
	Secret string
}

type InitMessageResponse struct {
	Success bool
}

type Event string

type SubscribeMessage struct {
	Events []Event
}
