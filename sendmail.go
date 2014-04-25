package sendmail

import (
	"io"
	"os/exec"
)

var (
	SendmailDefaultPath = "/usr/sbin/sendmail"
	SendmailDefaultArgs = []string{
		"-t",
	}
)

type Sendmail struct {
	// Path to the sendmail binary. If emtpy, SendmailDefaultPath is used
	Path string

	// Args to pass to sendmail. If nil, SendmailDefaultArgs is used
	Args []string
}

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
