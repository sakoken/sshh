package action

import (
	"github.com/sakoken/sshh/global"
)

func Add() error {
	var host global.Host
	host.Host = global.Question("HostName:", true, "")
	host.User = global.Question("UserName:", false, "")
	host.Password = global.Password("Password:", false)
	host.Port = global.Question("Port:", true, "22")
	//host.Key = Question("SSHKey:", true, host.Key)
	host.Explanation = global.Question("Explanation:", false, "")
	global.SshhData.Hosts = append(global.SshhData.Hosts, &host)
	global.SshhData.ResetPosition()

	return global.SaveJson(global.SshhData)
}
