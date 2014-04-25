package sendmail

import (
	"net/smtp"
)

// SMTP uses net/smtp to send the Mail
type SMTP struct {
	// Host is the SMTP-Server to connect to, for example "localhost:25"
	Host string
}

// Send implements the Sender interface. It uses both (*Mail).SanitizeHeaders
// and (*Mail).Encode
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
