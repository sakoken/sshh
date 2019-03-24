package action

import (
	"encoding/json"
	"github.com/sakoken/sshh/global"
	"io/ioutil"
)

func Delete(id int) error {
	var tempHosts []global.Host
	for k, v := range global.SshhData.Hosts {
		if k != id {
			tempHosts = append(tempHosts, v)
		}
	}
	global.SshhData.Hosts = tempHosts

	return SaveJson(global.SshhData)
}

func SaveJson(data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(global.SshhJson(), b, 0777)
}
