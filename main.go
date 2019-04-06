package main

import (
	"github.com/sakoken/sshh/action"
	"github.com/sakoken/sshh/global"
	"gopkg.in/urfave/cli.v2"
	"log"
	"os"
)

func main() {
	app := cli.App{}
	app.Name = "sshh"
	app.Description = "sshh is a management application of hosts for ssh"
	app.Version = "0.0.4"
	app.Before = before()
	app.Action = func(c *cli.Context) error {
		if c.Args().First() != "" {
			ssh := action.NewSsh()
			return ssh.CreateAndConnection(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2))
		}
		search := action.NewSeach()
		return search.Do(c.String("query"))
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "query",
			Aliases: []string{"q"},
			Usage:   "Query of search",
		},
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

func before() cli.BeforeFunc {
	return func(context *cli.Context) error {
		initSshh()
		return nil
	}
}

func initSshh() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	global.UserHome = home
	if err := os.MkdirAll(global.SshhHome(), 0777); err != nil {
		log.Fatal(err)
	}
	global.CreateJson(global.SshhJson())
	global.ReadJson(global.SshhJson(), &global.SshhData.Hosts)
	global.SshhData.ResetPosition()
}
