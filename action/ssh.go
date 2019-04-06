package action

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/sakoken/sshh/global"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"strings"
)

func NewSsh() *Ssh {
	return &Ssh{}
}

type Ssh struct {
}

func (s *Ssh) CreateAndConnection(requestHost string, pOption string, port string) error {
	host := &global.Host{}
	host.Host = requestHost
	host.Port = "22"
	if index := strings.Index(requestHost, "@"); index > 0 {
		host.Host = requestHost[index+1:]
		host.User = requestHost[:index]
	}
	if pOption == "-p" && port != "" {
		host.Port = port
	}

	l, _ := readline.NewEx(&readline.Config{
		Prompt:              "\033[36msshh»\033[0m ",
		InterruptPrompt:     "\n",
		EOFPrompt:           "exit",
		FuncFilterInputRune: global.FilterInput,
	})

	defer func(rl *readline.Instance) {
		if rl != nil {
			err := rl.Close()
			if err != nil {
				println(err.Error())
			}
		}
	}(l)

	//すでに登録済みの場合はすぐにssh
	if has, resistedHost := global.SshhData.Has(host); has {
		key := global.PasswordQuestion(l, "Enter secret key", true, 16)
		l.Close()
		pw, err := global.Decrypt(resistedHost.Password, key)
		if err != nil {
			return err
		}
		SshConnection(string(pw), resistedHost)
		return nil
	}

	key := ""
	host.Password, key = global.Password(l, "Password:", true)
	l.Close()

	global.SshhData.Hosts = append(global.SshhData.Hosts, host)
	global.SshhData.ResetPosition()
	global.SshhData.Save()

	pw, err := global.Decrypt(host.Password, key)
	if err != nil {
		return err
	}

	SshConnection(string(pw), host)
	return nil
}

func SshConnection(password string, host *global.Host) {
	println(fmt.Sprintf("\033[07m\033[34m%s\033[0m", host.SshCommand()))
	println(fmt.Sprintf("\033[07m\033[34mExplanation: %s\033[0m", host.Explanation))

	var auth []ssh.AuthMethod
	auth = append(auth, ssh.Password(password))

	sshConfig := &ssh.ClientConfig{
		User:            host.User,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host.Host+":"+host.Port, sshConfig)
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

	w, h, err := terminal.GetSize(fd)
	if err != nil {
		fmt.Println(err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	term := os.Getenv("TERM")
	err = session.RequestPty(term, h, w, modes)
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
