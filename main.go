package main

import (
	"encoding/json"
	"fmt"
	"github.com/sakoken/sshh/action"
	"github.com/sakoken/sshh/global"
	"gopkg.in/urfave/cli.v2"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func main() {
	app := cli.App{}
	app.Name = "sshh"
	app.Description = "sshh is a management application of hosts for ssh"
	app.Version = "0.0.1"
	app.Before = before()
	app.Action = action.Search
	app.Commands = commands()

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

func before() cli.BeforeFunc {
	return func(context *cli.Context) error {
		initSshh()
		readSshhFile()
		return nil
	}
}

func initSshh() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	global.UserHome = home

	if _, err := os.Stat(global.SshhJson()); !os.IsNotExist(err) {
		return
	}

	if err := os.MkdirAll(global.SshhHome(), 0777); err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(global.SshhJson(), os.O_WRONLY|os.O_CREATE, 0700)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fmt.Fprintln(file, "{}")
}

func readSshhFile() {
	dat, err := ioutil.ReadFile(global.SshhJson())
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(dat, &global.SshhData)
}

func commands() []*cli.Command {
	return []*cli.Command{
		{
			Name:   "list",
			Usage:  "show hosts",
			Action: action.List,
		},
		{
			Name:   "add",
			Usage:  "add hosts",
			Action: action.Add,
		},
		{
			Name:  "mod",
			Usage: "modify hosts",
			Action: func(c *cli.Context) error {
				id, err := strconv.Atoi(c.Args().First())
				if err != nil {
					return err
				}
				return action.Modify(id)
			},
		},
		{
			Name:  "del",
			Usage: "delete hosts",
			Action: func(c *cli.Context) error {
				id, err := strconv.Atoi(c.Args().First())
				if err != nil {
					return err
				}
				return action.Delete(id)
			},
		},
	}
}
