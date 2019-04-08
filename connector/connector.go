package connector

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
)

type Connector struct {
	Host        string `json:"host"`
	User        string `json:"user"`
	Port        string `json:"port"`
	Password    []byte `json:"password"`
	Key         string `json:"ssh_key"`
	Explanation string `json:"explanation"`
	Position    int    `json:"-"`
}

func (h *Connector) Clone() *Connector {
	cp := *h
	return &cp
}

func (h Connector) SshCommand() string {
	user := ""
	port := ""

	if h.User != "" {
		user = h.User + "@"
	}

	if h.Port != "" {
		port = "-p " + h.Port
	}
	return fmt.Sprintf("ssh %s%s %s", user, h.Host, port)
}

func (h Connector) SshConnection(password string) {
	SshhData.SetTopPosition(&h)
	SshhData.Save()
	println(fmt.Sprintf("\033[07m\033[34m%s\033[0m", h.SshCommand()))
	println(fmt.Sprintf("\033[07m\033[34mExplanation: %s\033[0m", h.Explanation))

	var auth []ssh.AuthMethod
	auth = append(auth, ssh.Password(password))

	sshConfig := &ssh.ClientConfig{
		User:            h.User,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", h.Host+":"+h.Port, sshConfig)
	if err != nil {
		log.Printf("%s error: %v\n", "dial", err)
		return
	}

	session, err := client.NewSession()
	if err != nil {
		log.Printf("%s error: %v\n", "new session", err)
		return
	}
	defer session.Close()

	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		fmt.Println(err)
	}
	defer terminal.Restore(fd, state)

	w, H, err := terminal.GetSize(fd)
	if err != nil {
		fmt.Println(err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	term := os.Getenv("TERM")
	err = session.RequestPty(term, H, w, modes)
	if err != nil {
		log.Printf("%s error: %v\n", "request pty", err)
		return
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	err = session.Shell()
	if err != nil {
		log.Printf("%s error: %v\n", "start shell", err)
		return
	}

	session.Wait()
}
