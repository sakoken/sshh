package action

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/sakoken/sshh/global"
	"gopkg.in/urfave/cli.v2"
	"io"
	"os"
	"strings"
)

func Add(_ *cli.Context) error {
	var host global.Host
	host.Host = Question("HostName:", true, "")
	host.User = Question("UserName:", false, "")
	host.Port = Question("HostName:", true, "22")
	host.Explain = Question("Explain:", false, "")

	global.SshhData.Hosts = append(global.SshhData.Hosts, host)

	return SaveJson(global.SshhData)
}

func Question(q string, required bool, def string) string {
	fmt.Println(q)
	result := def

	l, _ := readline.NewEx(&readline.Config{
		Prompt:              "\033[31m»\033[0m ",
		InterruptPrompt:     "\n",
		EOFPrompt:           "exit",
		FuncFilterInputRune: filterInput,
	})

	defer l.Close()

	setPasswordCfg := l.GenPasswordConfig()
	setPasswordCfg.SetListener(func(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
		l.SetPrompt(fmt.Sprintf("Enter password(%v): ", len(line)))
		l.Refresh()
		return nil, 0, false
	})

	var err error
	for {
		result, err = l.ReadlineWithDefault(def)
		result = strings.TrimSpace(result)

		if err == readline.ErrInterrupt {
			if len(result) != 0 {
				break
			} else {
				os.Exit(0)
			}
		} else if err == io.EOF {
			break
		}

		if len(result) == 0 && required {
			fmt.Println("Please enter correct answer.")
			continue
		}

		break
	}
	return result
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}
