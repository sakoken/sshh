package action

import (
	"github.com/sakoken/sshh/global"
)

func List() error {
	global.PrintTable(global.SshhData.Hosts)
	return nil
}
