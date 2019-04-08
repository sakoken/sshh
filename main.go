package main

import (
	"github.com/sakoken/sshh/action"
	"github.com/sakoken/sshh/config"
	"github.com/sakoken/sshh/connector"
	"github.com/sakoken/sshh/json"
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
		search := action.NewSearch()
		return search.Do(c.String("query"))
	}
	app.Usage = `https://github.com/sakoken/sshh/blob/master/README.md
	 [After exec sshh]
	 sshh>> #[positionNo]    :Do ssh connection
	 sshh>> add              :Add a new host 
	 sshh>> mod [positionNo] :Modify a host 
	 sshh>> del [positionNo] :Delete a host`
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
	config.UserHome = home
	if err := os.MkdirAll(config.SshhHome(), 0777); err != nil {
		log.Fatal(err)
	}
	json.CreateJson(config.SshhJson())
	json.ReadJson(config.SshhJson(), connector.SshhData)
	connector.SshhData.ResetPosition()
}
