package action

import (
	"github.com/sakoken/sshh/global"
)

func Delete(id int) error {
	var tempHosts []global.Host
	for k, v := range global.SshhData.Hosts {
		if k != id {
			tempHosts = append(tempHosts, v)
		}
	}
	global.SshhData.Hosts = tempHosts

	return global.SaveJson(global.SshhData)
}
