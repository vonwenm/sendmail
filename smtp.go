package sendmail

import (
	"net/smtp"
)

type SMTP struct {
	Host string
}

func (mailer SMTP) Send(m *Mail) error {
	m.SanitizeHeaders()

	c, err := smtp.Dial(mailer.Host)
	if err != nil {
		return err
	}
	defer c.Quit()

	if err = c.Mail(m.From); err != nil {
		return err
	}

	for _, rcpt := range m.To {
		if err = c.Rcpt(rcpt); err != nil {
			return err
		}
	}

	wc, err := c.Data()
	if err != nil {
		return err
	}
	return m.Encode(wc)
}
