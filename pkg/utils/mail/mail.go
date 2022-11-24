package mail

import (
	"context"
	"fmt"
	"github.com/mailgun/mailgun-go/v4"
	"time"
)

type Mail struct {
	Subject   string
	Template  string
	Variable  map[string]string
	Recipient string
}

type Client interface {
	SendMail(ctx context.Context, mail *Mail) error
}

type clientImpl struct {
	mg         mailgun.Mailgun
	sender     string
	senderName string
}

func NewClient(domain, apiKey string, sender string, senderName string) Client {
	return &clientImpl{
		mg:         mailgun.NewMailgun(domain, apiKey),
		sender:     sender,
		senderName: senderName,
	}
}

func (c *clientImpl) SendMail(ctx context.Context, mail *Mail) error {
	sender := fmt.Sprintf("%s <%s>", c.senderName, c.sender)
	m := c.mg.NewMessage(sender, mail.Subject, "", mail.Recipient)

	m.SetTemplate(mail.Template)
	for key, value := range mail.Variable {
		err := m.AddVariable(key, value)
		if err != nil {
			return err
		}
	}

	ct, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, _, err := c.mg.Send(ct, m)
	return err
}
