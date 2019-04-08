package global

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

type Sshh struct {
	Hosts []*Host `json:"hosts"`
}

func (s *Sshh) SetTopPosition(h *Host) {
	position := h.Position
	var tempHosts []*Host
	tempHosts = append(tempHosts, h)
	for k, v := range s.Hosts {
		if k != position {
			tempHosts = append(tempHosts, v)
		}
	}
	s.Hosts = tempHosts
	s.ResetPosition()
}

func (s *Sshh) ResetPosition() {
	for k, v := range s.Hosts {
		v.Position = k
	}
}

func (s *Sshh) Has(h *Host) (bool, *Host) {
	ch := h.Clone()
	ch.Password = nil
	ch.Explanation = ""
	hb, _ := json.Marshal(ch)
	for _, v := range s.Hosts {
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
	var tempHosts []*Host
	for k, v := range s.Hosts {
		if k != position {
			tempHosts = append(tempHosts, v)
		}
	}
	s.Hosts = tempHosts
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
	return ioutil.WriteFile(SshhJson(), buf.Bytes(), 0777)
}
