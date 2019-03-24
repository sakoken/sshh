module github.com/sakoken/sshh

go 1.12

require (
	github.com/atotto/clipboard v0.1.2
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e
	github.com/marcusolsson/tui-go v0.4.0
	golang.org/x/crypto v0.0.0-20190320223903-b7391e95e576 // indirect
	gopkg.in/urfave/cli.v2 v2.0.0-20180128182452-d3ae77c26ac8
)

replace github.com/marcusolsson/tui-go v0.4.0 => github.com/sakoken/tui-go v0.4.1
