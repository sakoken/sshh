package action

import (
	"github.com/chzyer/readline"
	"github.com/sakoken/sshh/connector"
	"github.com/sakoken/sshh/interactive"
)

func Modify(position int) error {
	i, _ := interactive.NewEx(&readline.Config{
		Prompt:              "\033[36msshh-mod»\033[0m ",
		InterruptPrompt:     "\n",
		EOFPrompt:           "exit",
		FuncFilterInputRune: interactive.FilterInput,
	})

	defer i.Close()

	var host = connector.SshhData().Connectors[position]
	host = host.Clone()
	host.Host = i.Question("HostName:", true, host.Host)
	host.User = i.Question("UserName:", false, host.User)
	host.Port = i.Question("PortNumber:", true, host.Port)

	if has, hasHost := connector.SshhData().Has(host); has && host.Position != hasHost.Position {
		println("\033[31mAlready exists\033[00m")
		return nil
	}

	pswd, _ := i.Password("Password:", false)
	if len(pswd) > 0 {
		host.Password = pswd
	}
	//host.Key = Question("SSHKey:", true, host.Key)
	host.Explanation = i.Question("Explanation:", false, host.Explanation)
	connector.SshhData().Connectors[position] = host
	return connector.SshhData().Save()
}
