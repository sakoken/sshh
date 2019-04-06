package model

import "fmt"

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
