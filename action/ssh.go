package action

import (
	"github.com/chzyer/readline"
	"github.com/sakoken/sshh/connector"
	"github.com/sakoken/sshh/encrypt"
	"github.com/sakoken/sshh/interactive"
	"strings"
)

func CreateAndConnection(requestHost string, pOption string, port string) error {
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

	i, _ := interactive.NewEx(&readline.Config{
		Prompt:              "\033[36msshh»\033[0m ",
		InterruptPrompt:     "\n",
		EOFPrompt:           "exit",
		FuncFilterInputRune: interactive.FilterInput,
	})

	defer func(rl *interactive.Interactive) {
		if rl != nil {
			err := rl.Close()
			if err != nil {
				println(err.Error())
			}
		}
	}(i)

	//すでに登録済みの場合はすぐにssh
	if has, resistedHost := connector.SshhData().Has(host); has {
		key := i.PasswordQuestion("Enter secret key", true, 16)
		i.Close()
		pw, err := encrypt.Decrypt(resistedHost.Password, key)
		if err != nil {
			return err
		}
		connector.SshhData().SetTopPosition(resistedHost).Save()
		resistedHost.SshConnection(string(pw))
		return nil
	}

	key := ""
	host.Password, key = i.Password("Password:", true)
	i.Close()

	connector.SshhData().Add(host).Save()

	pw, err := encrypt.Decrypt(host.Password, key)
	if err != nil {
		return err
	}
	connector.SshhData().SetTopPosition(host).Save()
	host.SshConnection(string(pw))
	return nil
}
