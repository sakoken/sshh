package action

import (
	"github.com/sakoken/sshh/global"
)

func Modify(id int) error {
	var host = &global.SshhData.Hosts[id]
	host.Host = Question("HostName:", true, host.Host)
	host.User = Question("UserName:", false, host.User)
	host.Port = Question("HostName:", true, host.Port)
	host.Explain = Question("Explain:", false, host.Explain)
	return SaveJson(global.SshhData)
}
