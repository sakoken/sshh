package global

import (
	"fmt"
)

const (
	SshhHomeName = ".sshh"
	SshhJsonName = "sshh.json"
)

var (
	UserHome = ""
	SshhHome = func() string { return UserHome + "/" + SshhHomeName }
	SshhJson = func() string { return SshhHome() + "/" + SshhJsonName }
	SshhData Sshh
)

type Sshh struct {
	Hosts []*Host `json:"hosts"`
}

type Host struct {
	Host        string `json:"host"`
	User        string `json:"user"`
	Port        string `json:"port"`
	Password    []byte `json:"password"`
	Key         string `json:"ssh_key"`
	Explanation string `json:"explanation"`
	Position    int    `json:"-"`
}

func (h Host) SshCommand() string {
	user := ""
	port := ""

	if h.User != "" {
		user = h.User + "@"
	}

	if h.Port != "" {
		port = "-p " + h.Port
	}
	return fmt.Sprintf("ssh %s%s %s", user, h.Host, port)
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
