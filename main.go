// package sendmail sends MIME-Mails over smtp or the OSes MTA
//
// Please note, that this package is experimental, it's API might still change
// unannounced.
package sendmail

import (
	"fmt"
	"github.com/Merovius/qp"
	"io"
	"strings"
	"time"
)

type Headers map[string]string

// Mail represents a single E-Mail
type Mail struct {
	// From is the sender of the email
	From string

	// To is a slice of recipients
	To []string

	// Subject is the subject of the Mail as an UTF-8 string
	Subject string

	// Headers are the headers
	Headers Headers

	// Body provides the actual body of the mail. It has to be UTF-8 encoded,
	// or you must set the Content-Type Header
	Body io.Reader
}

// Sender is an interface with a Send method, that dispatches a single Mail
type Sender interface {
	// Send the given E-Mail via this sender. Send might use
	// (*Mail).SanitizeHeaders and (*Mail).Encode as helpers
	Send(*Mail) error
}

// NewMail returns a new Mail with Headers initialized to an empty map
func NewMail() *Mail {
	return &Mail{Headers: make(map[string]string)}
}

// SanitizeHeaders sets some default-headers and encodes the Subject. It should
// not be called manually, but should be called by the used Sender
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

// Encode writes the encoded Mail, including the headers to the provided
// Writer. It should not be called manually, but should be called by the used
// Sender
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

// Set will set the Header with the given key to the given Value, if not
// already present
func (h Headers) Set(key, value string) {
	if _, ok := h[key]; ok {
		return
	}
	h[key] = value
}

// Overwrite will set the Hedaer with the given key to the given Value,
// overwriting it, if present
func (h Headers) Overwrite(key, value string) {
	h[key] = value
}
