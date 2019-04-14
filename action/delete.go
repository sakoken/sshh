package action

import (
	"github.com/sakoken/sshh/config"
	"github.com/sakoken/sshh/connector"
)

func Delete(position int) {
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
