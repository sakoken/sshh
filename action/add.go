package action

import (
	"github.com/chzyer/readline"
	"github.com/sakoken/sshh/global"
)

func Add() error {
	l, _ := readline.NewEx(&readline.Config{
		Prompt:              "\033[36msshh-add»\033[0m ",
		InterruptPrompt:     "\n",
		EOFPrompt:           "exit",
		FuncFilterInputRune: global.FilterInput,
	})

	defer l.Close()

	host := &global.Host{}
	host.Host = global.Question(l, "HostName:", true, "")
	host.User = global.Question(l, "UserName:", false, "")
	host.Password, _ = global.Password(l, "Password:", false)
	host.Port = global.Question(l, "Port:", true, "22")
	//host.Key = Question("SSHKey:", true, host.Key)
	host.Explanation = global.Question(l, "Explanation:", false, "")
	if has, _ := global.SshhData.Has(host); has {
		println("\033[31mAlready exists\033[00m")
		return nil
	}
	global.SshhData.Hosts = append(global.SshhData.Hosts, host)
	global.SshhData.ResetPosition()
	return global.SshhData.Save()
}
