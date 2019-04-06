package global

import "github.com/sakoken/sshh/model"

const (
	SshhHomeName = ".sshh"
	SshhJsonName = "sshh.json"
)

var (
	UserHome = ""
	SshhHome = func() string { return UserHome + "/" + SshhHomeName }
	SshhJson = func() string { return SshhHome() + "/" + SshhJsonName }
	SshhData model.Sshh
)
