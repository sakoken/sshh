package action

import (
	"github.com/sakoken/sshh/global"
	"github.com/sakoken/sshh/model"
)

func Delete(position int) error {
	var tempHosts []*model.Host
	for k, v := range global.SshhData.Hosts {
		if k != position {
			tempHosts = append(tempHosts, v)
		}
	}
	global.SshhData.Hosts = tempHosts
	global.SshhData.ResetPosition()

	return global.SaveJson(global.SshhData)
}
