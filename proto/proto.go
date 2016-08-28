package proto

type InitMessage struct {
	Name   string
	Secret string
}

type AuthMessage struct {
	Success bool
}
