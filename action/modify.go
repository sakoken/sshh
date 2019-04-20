package action

import (
	"github.com/chzyer/readline"
	"github.com/sakoken/sshh/config"
	"github.com/sakoken/sshh/interactive"
)

func Modify(position int) error {
	i, _ := interactive.NewEx(&readline.Config{
		InterruptPrompt:     "\n",
		EOFPrompt:           "exit",
		FuncFilterInputRune: interactive.FilterInput,
	})

	defer i.Close()

	var host = config.SshhData().Connectors[position]
	host = host.Clone()
	host.Host = i.Question("HostName", true, host.Host)
	host.User = i.Question("UserName", false, host.User)
	host.Port = i.Question("PortNumber", true, host.Port)

	if has, hasHost := config.SshhData().Has(host); has && host.Position != hasHost.Position {
		println("\033[31mAlready exists\033[00m")
		return nil
	}

	pswd, _ := i.Password("Password", false)
	if len(pswd) > 0 {
		host.Password = pswd
	}
	//host.Key = Question("SSHKey:", true, host.Key)
	host.Explanation = i.Question("Explanation", false, host.Explanation)
	config.SshhData().Connectors[position] = host
	return config.SshhData().Save()
}
