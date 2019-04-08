package config

const (
	SshhHomeName = ".sshh"
	SshhJsonName = "sshh.json"
)

var (
	UserHome = ""
	SshhHome = func() string { return UserHome + "/" + SshhHomeName }
	SshhJson = func() string { return SshhHome() + "/" + SshhJsonName }
)
