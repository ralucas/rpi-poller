package message

type Message struct {
	subject    string
	message    string
	receipient string
}

func New(subject string, msg string, receipient string) Message {
	return Message{
		subject:    subject,
		message:    msg,
		receipient: receipient,
	}
}

func (m Message) GetMessage() string {
	return m.message
}

func (m Message) GetSubject() string {
	return m.subject
}

func (m Message) GetReceipient() string {
	return m.receipient
}
