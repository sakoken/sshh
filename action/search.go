package action

import (
	"github.com/sakoken/sshh/config"
	"github.com/sakoken/sshh/connector"
	"strings"
)

func Search(keyword string) (hosts []*connector.Connector) {
	for _, v := range config.SshhData().Connectors {
		if strings.Index(v.Host, keyword) >= 0 ||
			strings.Index(v.User, keyword) >= 0 ||
			strings.Index(v.Port, keyword) >= 0 ||
			strings.Index(v.Explanation, keyword) >= 0 {
			hosts = append(hosts, v)
		}
	}
	return
}
