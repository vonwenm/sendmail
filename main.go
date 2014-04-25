package sendmail

import (
	"fmt"
	"github.com/Merovius/qp"
	"io"
	"strings"
	"time"
)

type Headers map[string]string

type Mail struct {
	From    string
	To      []string
	Subject string
	Headers Headers
	Body    io.Reader
}

type Mailer interface {
	Send(Mail) error
}

func NewMail() *Mail {
	return &Mail{Headers: make(map[string]string)}
}

func (m *Mail) SanitizeHeaders() {
	m.Headers.Overwrite("From", m.From)
	m.Headers.Overwrite("To", strings.Join(m.To, ", "))
	if m.Subject == "" {
		m.Headers.Overwrite("Subject", "")
	} else {
		m.Headers.Overwrite("Subject", string(qp.EncodedWord(m.Subject)))
	}

	m.Headers.Overwrite("MIME-Version", "1.0")
	m.Headers.Set("Content-Type", "text/plain; charset=\"utf-8\"")
	m.Headers.Set("Content-Disposition", "inline")
	m.Headers.Overwrite("Content-Transfer-Encoding", "quoted-printable")
	m.Headers.Set("User-Agent", "go-sendmail/1.0")
	m.Headers.Set("Date", time.Now().Format(time.RFC1123Z))
	// TODO: Message-Id?
}

func (m *Mail) Encode(w io.Writer) (err error) {
	for key, val := range m.Headers {
		_, err = fmt.Fprintf(w, "%s: %s\r\n", key, val)
		if err != nil {
			return err
		}
	}
	io.WriteString(w, "\r\n")

	enc := qp.NewWriter(w)
	_, err = io.Copy(enc, m.Body)
	return err
}

func (h Headers) Set(key, value string) {
	if _, ok := h[key]; ok {
		return
	}
	h[key] = value
}

func (h Headers) Overwrite(key, value string) {
	h[key] = value
}
