package message

type Message struct {
	subject string
	message string
}

func New(subject string, msg string) Message {
	return Message{
		subject: subject,
		message: msg,
	}
}

func (m Message) GetMessage() string {
	return m.message
}

func (m Message) GetSubject() string {
	return m.subject
}
