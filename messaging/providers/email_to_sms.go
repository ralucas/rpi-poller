package providers

import (
	"fmt"
	"os/exec"

	"github.com/ralucas/rpi-poller/messaging/message"
)

type EmailToSMS struct{}

const command = "mail -s %s %s"

func (e *EmailToSMS) Send(msg message.Message) error {
	cmd := exec.Command(fmt.Sprintf(command, msg.GetMessage(), msg.GetReceipient()))
	return cmd.Run()
}
