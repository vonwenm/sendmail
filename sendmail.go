package sendmail

import (
	"io"
	"os/exec"
)

var (
	// The default path where to search for the sendmail-binary
	SendmailDefaultPath = "/usr/sbin/sendmail"
	// The default arguments to pass to sendmail
	SendmailDefaultArgs = []string{
		"-t",
	}
)

type Sendmail struct {
	// Path to the sendmail binary. If emtpy, SendmailDefaultPath is used
	Path string

	// Args to pass to sendmail. There might be format strings implemented for
	// this in the future, so to provide forward compatibility, you should
	// avoid '%'-characters. If nil, SendmailDefaultArgs is used
	Args []string
}

// Send implements the Sender interface. It uses both (*Mail).SanitizeHeaders()
// and (*Mail).Encode
func (mailer Sendmail) Send(m *Mail) error {
	m.SanitizeHeaders()

	if mailer.Path == "" {
		mailer.Path = "/usr/sbin/sendmail"
	}

	if mailer.Args == nil {
		mailer.Args = SendmailDefaultArgs
	}

	cmd := exec.Command(mailer.Path, mailer.Args...)
	w, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	if err = m.Encode(w); err != nil {
		return err
	}

	_, err = io.WriteString(w, ".\r\n")
	return err
}
