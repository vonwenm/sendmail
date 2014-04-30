package sendmail

import (
	"bytes"
	"log"
	"strings"
)

// Dummy is a non-functional Sender, used for debugging
type Dummy struct {
	// If not nil, write the mail to this logger
	Log *log.Logger
}

func (mailer Dummy) Send(m *Mail) error {
	if mailer.Log == nil {
		return nil
	}

	m.SanitizeHeaders()

	mailer.Log.Println("Sending mail from", m.From)
	mailer.Log.Println("Recipients:", strings.Join(m.To, ", "))
	mailer.Log.Println("=== BODY START ===")
	defer mailer.Log.Println("=== BODY END ===")

	body := new(bytes.Buffer)
	if err := m.Encode(body); err != nil {
		return err
	}

	mailer.Log.Println(body)
	return nil
}
