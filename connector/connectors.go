package connector

import (
	"fmt"
	"strconv"
	"strings"
)

type ConnectorCollection struct {
	Connectors []*Connector `json:"hosts"`
}

func (s *ConnectorCollection) Add(h *Connector) *ConnectorCollection {
	s.Connectors = append(s.Connectors, h)
	s.ResetPosition()
	return s
}

func (s *ConnectorCollection) ToTopPosition(h *Connector) *ConnectorCollection {
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
	return s
}

func (s *ConnectorCollection) ResetPosition() *ConnectorCollection {
	for k, v := range s.Connectors {
		v.Position = k
	}
	return s
}

func (s *ConnectorCollection) Index(position int) *Connector {
	return s.Connectors[position]
}

func (s *ConnectorCollection) Count() int {
	return len(s.Connectors)
}

func (s *ConnectorCollection) Has(h *Connector) (bool, *Connector) {
	for _, v := range s.Connectors {
		if h.Equals(v) {
			return true, v
		}
	}
	return false, nil
}

func (s *ConnectorCollection) PrintTable() {
	hostLen, portLen, userLen := s.maxLength(4)

	hostStr := "Host"
	portStr := "Port"
	userStr := "User"

	hostStr = hostStr + strings.Repeat(" ", hostLen-len(hostStr))
	portStr = portStr + strings.Repeat(" ", portLen-len(portStr))
	userStr = userStr + strings.Repeat(" ", userLen-len(userStr))

	lenOfNo := len(strconv.Itoa(s.Count()))
	noSpace := strings.Repeat(" ", lenOfNo-1)
	println(fmt.Sprintf("[#]"+noSpace+" %s  %s  %s  %s", hostStr, userStr, portStr, "Explanation"))
	noStr := "[%d]"
	for k, v := range s.Connectors {
		noSpace = strings.Repeat(" ", lenOfNo-len(strconv.Itoa(k)))
		host := v.Host + strings.Repeat(" ", hostLen-len(v.Host))
		port := v.Port + strings.Repeat(" ", portLen-len(v.Port))
		user := v.User + strings.Repeat(" ", userLen-len(v.User))
		println(fmt.Sprintf(noStr+noSpace+" %s  %s  %s  %s", k, host, user, port, v.Explanation))
	}
}

func (s *ConnectorCollection) maxLength(def int) (hostLen, portLen, userLen int) {
	hostLen = def
	portLen = def
	userLen = def
	for _, v := range s.Connectors {
		if len(v.Host) > hostLen {
			hostLen = len(v.Host)
		}
		if len(v.Port) > portLen {
			portLen = len(v.Port)
		}
		if len(v.User) > userLen {
			userLen = len(v.User)
		}
	}
	return
}
