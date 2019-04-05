package action

import (
	"github.com/sakoken/sshh/global"
)

func Modify(position int) error {
	var host = global.SshhData.Hosts[position]
	host.Host = global.Question("HostName:", true, host.Host)
	host.User = global.Question("UserName:", false, host.User)
	host.Port = global.Question("PortNumber:", true, host.Port)
	//host.Key = Question("SSHKey:", true, host.Key)
	host.Explanation = global.Question("Explanation:", false, host.Explanation)
	return global.SaveJson(global.SshhData)
}
