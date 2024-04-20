package _ssh

import (
	"fmt"
	"github.com/MayMistery/maygit/_ssh"
	"github.com/MayMistery/maygit/cmd"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type SSHTerminal struct {
	Session *ssh.Session
	exitMsg string
	stdout  io.Reader
	stdin   io.Writer
	stderr  io.Reader
}

func SSHSession(cfg cmd.Config) {
	client, err := _ssh.EstablishSSHConnection(cfg)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer client.Close()

	err = newSSHSession(client)
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func (t *SSHTerminal) updateTerminalSize() {
	go func() {
		// SIGWINCH is sent to the process when the window size of the terminal has
		// changed.
		sigwinchCh := make(chan os.Signal, 1)

		// TODO in windows undefined: syscall.SIGWINCH

		signal.Notify(sigwinchCh, syscall.SIGWINCH)

		fd := int(os.Stdin.Fd())
		termWidth, termHeight, err := terminal.GetSize(fd)
		if err != nil {
			fmt.Println(err)
		}

		for {
			select {
			// The client updated the size of the local PTY. This change needs to occur
			// on the server side PTY as well.
			case sigwinch := <-sigwinchCh:
				if sigwinch == nil {
					return
				}
				currTermWidth, currTermHeight, err := terminal.GetSize(fd)
				if err != nil {
					fmt.Printf("Unable to get windows size: %s.", err)
					continue
				}
				// Terminal size has not changed, don't do anything.
				if currTermHeight == termHeight && currTermWidth == termWidth {
					continue
				}

				err = t.Session.WindowChange(currTermHeight, currTermWidth)
				if err != nil {
					fmt.Printf("Unable to send window-change reqest: %s.", err)
					continue
				}

				termWidth, termHeight = currTermWidth, currTermHeight

			}
		}
	}()
}

func (t *SSHTerminal) interactiveSession() error {

	defer func() {
		var err error
		if t.exitMsg == "" {
			_, err = fmt.Fprintln(os.Stdout, "the connection was closed on the remote side on ", time.Now().Format(time.RFC822))
		} else {
			_, err = fmt.Fprintln(os.Stdout, t.exitMsg)
		}
		if err != nil {
			fmt.Printf("Fprintln error: %v\n", err)
		}
	}()

	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer func(fd int, oldState *terminal.State) {
		err := terminal.Restore(fd, oldState)
		if err != nil {
			log.Fatalf("Restore Fatal failed : %v", err)
		}
	}(fd, state)

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		return err
	}

	termType := os.Getenv("TERM")
	if termType == "" {
		termType = "xterm-256color"
	}

	err = t.Session.RequestPty(termType, termHeight, termWidth, ssh.TerminalModes{})
	if err != nil {
		return err
	}

	// TODO to handle error in windows
	t.updateTerminalSize()

	t.stdin, err = t.Session.StdinPipe()
	if err != nil {
		return err
	}
	t.stdout, err = t.Session.StdoutPipe()
	if err != nil {
		return err
	}
	t.stderr, err = t.Session.StderrPipe()

	// TODO to handle error
	go io.Copy(os.Stderr, t.stderr)
	go io.Copy(os.Stdout, t.stdout)
	go func() {
		buf := make([]byte, 128)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				fmt.Println(err)
				return
			}
			if n > 0 {
				_, err = t.stdin.Write(buf[:n])
				if err != nil {
					fmt.Println(err)
					t.exitMsg = err.Error()
					return
				}
			}
		}
	}()

	err = t.Session.Shell()
	if err != nil {
		return err
	}
	err = t.Session.Wait()
	if err != nil {
		return err
	}
	return nil
}

func newSSHSession(client *ssh.Client) error {

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	s := SSHTerminal{
		Session: session,
	}

	return s.interactiveSession()
}
