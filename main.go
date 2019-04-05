package main

import (
	"github.com/sakoken/sshh/action"
	"github.com/sakoken/sshh/global"
	"gopkg.in/urfave/cli.v2"
	"log"
	"os"
	"strconv"
)

func main() {
	app := cli.App{}
	app.Name = "sshh"
	app.Description = "sshh is a management application of hosts for ssh"
	app.Version = "0.0.2"
	app.Before = before()
	app.Action = func(c *cli.Context) error {
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
	app.Commands = commands()

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
	global.ReadJson(global.SshhJson(), &global.SshhData)
	global.SshhData.ResetPosition()
}

func commands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "list",
			Usage: "show hosts",
			Action: func(_ *cli.Context) error {
				return action.List()
			},
		},
		{
			Name:  "add",
			Usage: "add hosts",
			Action: func(_ *cli.Context) error {
				return action.Add()
			},
		},
		{
			Name:  "mod",
			Usage: "modify hosts",
			Action: func(c *cli.Context) error {
				position, err := strconv.Atoi(c.Args().First())
				if err != nil {
					return err
				}
				return action.Modify(position)
			},
		},
		{
			Name:  "del",
			Usage: "delete hosts",
			Action: func(c *cli.Context) error {
				position, err := strconv.Atoi(c.Args().First())
				if err != nil {
					return err
				}
				return action.Delete(position)
			},
		},
	}
}
