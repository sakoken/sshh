package action

import (
	"fmt"
	"github.com/sakoken/sshh/global"
)

func List() error {
	for k, v := range global.SshhData.Hosts {
		fmt.Printf("[%d] %s  explain:%s  user:%s  port:%s \n", k, v.Host, v.Explanation, v.User, v.Port)
	}
	return nil
}
