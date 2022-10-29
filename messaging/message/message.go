package message

type Message struct {
	message string
}

func New(msg string) Message {
	return Message{
		message: msg,
	}
}

func (m Message) GetMessage() string {
	return m.message
}
