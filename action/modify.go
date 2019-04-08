package action

import (
	"github.com/chzyer/readline"
	"github.com/sakoken/sshh/connector"
	"github.com/sakoken/sshh/interactive"
)

func Modify(position int) error {
	l, _ := readline.NewEx(&readline.Config{
		Prompt:              "\033[36msshh-mod»\033[0m ",
		InterruptPrompt:     "\n",
		EOFPrompt:           "exit",
		FuncFilterInputRune: interactive.FilterInput,
	})

	defer l.Close()

	var host = connector.SshhData.Connectors[position]
	host = host.Clone()
	host.Host = interactive.Question(l, "HostName:", true, host.Host)
	host.User = interactive.Question(l, "UserName:", false, host.User)
	host.Port = interactive.Question(l, "PortNumber:", true, host.Port)
	pswd, _ := interactive.Password(l, "Password:", false)
	if len(pswd) > 0 {
		host.Password = pswd
	}

	//host.Key = Question("SSHKey:", true, host.Key)
	host.Explanation = interactive.Question(l, "Explanation:", false, host.Explanation)
	if has, hasHost := connector.SshhData.Has(host); has && host.Position != hasHost.Position {
		println("\033[31mAlready exists\033[00m")
		return nil
	}
	connector.SshhData.Connectors[position] = host
	return connector.SshhData.Save()
}
