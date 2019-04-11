package action

import (
	"github.com/chzyer/readline"
	"github.com/sakoken/sshh/connector"
	"github.com/sakoken/sshh/interactive"
)

func Add() error {
	i, _ := interactive.NewEx(&readline.Config{
		Prompt:              "\033[36msshh-add»\033[0m ",
		InterruptPrompt:     "\n",
		EOFPrompt:           "exit",
		FuncFilterInputRune: interactive.FilterInput,
	})

	defer i.Close()

	host := &connector.Connector{}
	host.Host = i.Question("HostName:", true, "")
	host.User = i.Question("UserName:", false, "")
	host.Port = i.Question("Port:", true, "22")
	if has, _ := connector.SshhData().Has(host); has {
		println("\033[31mAlready exists\033[00m")
		return nil
	}
	host.Password, _ = i.Password("Password:", false)
	//host.Key = Question("SSHKey:", true, host.Key)
	host.Explanation = i.Question("Explanation:", false, "")

	return connector.SshhData().Add(host).Save()
}
