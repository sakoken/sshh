package config

import (
	"sync"

	"github.com/sakoken/sshh/connector"
)

var sshhData *Sshh
var once sync.Once

func SshhData() *Sshh {
	once.Do(func() {
		ReadJson(SshhJson(), &sshhData)
		sshhData.ResetPosition()
	})
	return sshhData
}

type Sshh struct {
	*connector.ConnectorCollection
}

func (s *Sshh) Save() error {
	return WriteJson(SshhJson(), s)
}
