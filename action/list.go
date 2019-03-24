package action

import (
	"fmt"
	"github.com/sakoken/sshh/global"
	"gopkg.in/urfave/cli.v2"
)

func List(_ *cli.Context) error {
	for k, v := range global.SshhData.Hosts {
		fmt.Printf("[%d] %s  explain:%s  user:%s  port:%s \n", k, v.Host, v.Explain, v.User, v.Port)
	}
	return nil
}
