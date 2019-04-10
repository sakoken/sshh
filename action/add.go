package action

import (
	"github.com/chzyer/readline"
	"github.com/sakoken/sshh/connector"
	"github.com/sakoken/sshh/interactive"
)

func Add() error {
	l, _ := readline.NewEx(&readline.Config{
		Prompt:              "\033[36msshh-add»\033[0m ",
		InterruptPrompt:     "\n",
		EOFPrompt:           "exit",
		FuncFilterInputRune: interactive.FilterInput,
	})

	defer l.Close()

	host := &connector.Connector{}
	host.Host = interactive.Question(l, "HostName:", true, "")
	host.User = interactive.Question(l, "UserName:", false, "")
	host.Port = interactive.Question(l, "Port:", true, "22")
	if has, _ := connector.SshhData().Has(host); has {
		println("\033[31mAlready exists\033[00m")
		return nil
	}
	host.Password, _ = interactive.Password(l, "Password:", false)
	//host.Key = Question("SSHKey:", true, host.Key)
	host.Explanation = interactive.Question(l, "Explanation:", false, "")

	return connector.SshhData().Add(host).Save()
}
