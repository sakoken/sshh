package action

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/sakoken/sshh/config"
	"github.com/sakoken/sshh/connector"
	"github.com/sakoken/sshh/interactive"
)

func Delete(position int) {
	i, _ := interactive.NewEx(&readline.Config{
		InterruptPrompt:     "\n",
		EOFPrompt:           "exit",
		FuncFilterInputRune: interactive.FilterInput,
	})
	defer i.Close()

	con := config.SshhData().Connectors[position]
	if !i.YesNo(fmt.Sprintf("Really you want delete [%s] y/n", con.CommandBase())) {
		return
	}

	var tempHosts []*connector.Connector
	for k, v := range config.SshhData().Connectors {
		if k != position {
			tempHosts = append(tempHosts, v)
		}
	}
	config.SshhData().Connectors = tempHosts
	config.SshhData().ResetPosition()
	config.SshhData().Save()
}
