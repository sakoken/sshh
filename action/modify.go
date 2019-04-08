package action

import (
	"github.com/chzyer/readline"
	"github.com/sakoken/sshh/global"
)

func Modify(position int) error {
	l, _ := readline.NewEx(&readline.Config{
		Prompt:              "\033[36msshh-mod»\033[0m ",
		InterruptPrompt:     "\n",
		EOFPrompt:           "exit",
		FuncFilterInputRune: global.FilterInput,
	})

	defer l.Close()

	var host = global.SshhData.Hosts[position]
	host = host.Clone()
	host.Host = global.Question(l, "HostName:", true, host.Host)
	host.User = global.Question(l, "UserName:", false, host.User)
	host.Port = global.Question(l, "PortNumber:", true, host.Port)
	pswd, _ := global.Password(l, "Password:", false)
	if len(pswd) > 0 {
		host.Password = pswd
	}

	//host.Key = Question("SSHKey:", true, host.Key)
	host.Explanation = global.Question(l, "Explanation:", false, host.Explanation)
	if has, hasHost := global.SshhData.Has(host); has && host.Position != hasHost.Position {
		println("\033[31mAlready exists\033[00m")
		return nil
	}
	global.SshhData.Hosts[position] = host
	return global.SshhData.Save()
}
