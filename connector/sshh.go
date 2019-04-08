package connector

import (
	"bytes"
	"encoding/json"
	"github.com/sakoken/sshh/config"
	"io/ioutil"
)

var SshhData Sshh

type Sshh struct {
	Connectors []*Connector `json:"hosts"`
}

func (s *Sshh) SetTopPosition(h *Connector) {
	position := h.Position
	var tempHosts []*Connector
	tempHosts = append(tempHosts, h)
	for k, v := range s.Connectors {
		if k != position {
			tempHosts = append(tempHosts, v)
		}
	}
	s.Connectors = tempHosts
	s.ResetPosition()
}

func (s *Sshh) ResetPosition() {
	for k, v := range s.Connectors {
		v.Position = k
	}
}

func (s *Sshh) Has(h *Connector) (bool, *Connector) {
	ch := h.Clone()
	ch.Password = nil
	ch.Explanation = ""
	hb, _ := json.Marshal(ch)
	for _, v := range s.Connectors {
		v2 := v.Clone()
		v2.Password = nil
		v2.Explanation = ""
		b, _ := json.Marshal(v2)
		if string(hb) == string(b) {
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
	s.ResetPosition()
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
