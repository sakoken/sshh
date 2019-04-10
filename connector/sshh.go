package connector

import (
	"bytes"
	"encoding/json"
	"github.com/sakoken/sshh/config"
	"io/ioutil"
	"sync"
)

var sshhData *Sshh
var once sync.Once

func SshhData() *Sshh {
	once.Do(func() {
		config.ReadJson(config.SshhJson(), &sshhData)
		sshhData.resetPosition()
	})

	return sshhData
}

type Sshh struct {
	Connectors []*Connector `json:"hosts"`
}

func (s *Sshh) Add(h *Connector) *Sshh {
	s.Connectors = append(s.Connectors, h)
	s.resetPosition()
	return s
}

func (s *Sshh) SetTopPosition(h *Connector) *Sshh {
	position := h.Position
	var tempHosts []*Connector
	tempHosts = append(tempHosts, h)
	for k, v := range s.Connectors {
		if k != position {
			tempHosts = append(tempHosts, v)
		}
	}
	s.Connectors = tempHosts
	s.resetPosition()
	return s
}

func (s *Sshh) resetPosition() *Sshh {
	for k, v := range s.Connectors {
		v.Position = k
	}
	return s
}

func (s *Sshh) Has(h *Connector) (bool, *Connector) {
	for _, v := range s.Connectors {
		if h.Equals(v) {
			return true, v
		}
	}
	return false, nil
}

func (s *Sshh) Delete(position int) {
	var tempHosts []*Connector
	for k, v := range s.Connectors {
		if k != position {
			tempHosts = append(tempHosts, v)
		}
	}
	s.Connectors = tempHosts
	s.resetPosition()
	s.Save()
}

func (s *Sshh) Save() error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err = json.Indent(&buf, []byte(b), "", "  "); err != nil {
		return err
	}
	return ioutil.WriteFile(config.SshhJson(), buf.Bytes(), 0777)
}
