package model

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
