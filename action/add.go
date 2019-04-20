package action

import (
	"github.com/chzyer/readline"
	"github.com/sakoken/sshh/config"
	"github.com/sakoken/sshh/connector"
	"github.com/sakoken/sshh/interactive"
)

func Add() error {
	i, _ := interactive.NewEx(&readline.Config{
		InterruptPrompt:     "\n",
		EOFPrompt:           "exit",
		FuncFilterInputRune: interactive.FilterInput,
	})

	defer i.Close()

	host := &connector.Connector{}
	host.Host = i.Question("HostName", true, "")
	host.User = i.Question("UserName", false, "")
	host.Port = i.Question("Port", true, "22")
	if has, _ := config.SshhData().Has(host); has {
		println("\033[31mAlready exists\033[00m")
		return nil
	}
	host.Password, _ = i.ServerPassword(false)
	//host.Key = Question("SSHKey:", true, host.Key)
	host.Explanation = i.Question("Explanation", false, "")
	config.SshhData().Add(host)
	return config.SshhData().Save()
}
