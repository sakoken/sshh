package action

import (
	"github.com/chzyer/readline"
	"github.com/sakoken/sshh/connector"
	"github.com/sakoken/sshh/encrypt"
	"github.com/sakoken/sshh/interactive"
	"strings"
)

func NewSsh() *Ssh {
	return &Ssh{}
}

type Ssh struct {
}

func (s *Ssh) CreateAndConnection(requestHost string, pOption string, port string) error {
	host := &connector.Connector{}
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
		FuncFilterInputRune: interactive.FilterInput,
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
	if has, resistedHost := connector.SshhData.Has(host); has {
		key := interactive.PasswordQuestion(l, "Enter secret key", true, 16)
		l.Close()
		pw, err := encrypt.Decrypt(resistedHost.Password, key)
		if err != nil {
			return err
		}
		resistedHost.SshConnection(string(pw))
		return nil
	}

	key := ""
	host.Password, key = interactive.Password(l, "Password:", true)
	l.Close()

	connector.SshhData.Connectors = append(connector.SshhData.Connectors, host)
	connector.SshhData.ResetPosition()
	connector.SshhData.Save()

	pw, err := encrypt.Decrypt(host.Password, key)
	if err != nil {
		return err
	}

	host.SshConnection(string(pw))
	return nil
}
