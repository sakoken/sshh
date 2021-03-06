package interactive

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/sakoken/sshh/encrypt"
)

type Interactive struct {
	readline.Instance
}

type defaultPainter struct{}

func (p *defaultPainter) Paint(line []rune, _ int) []rune {
	return line
}

func NewEx(cfg *readline.Config) (*Interactive, error) {
	t, err := readline.NewTerminal(cfg)
	if err != nil {
		return nil, err
	}
	rl := t.Readline()
	if cfg.Painter == nil {
		cfg.Painter = &defaultPainter{}
	}
	i := &Interactive{}
	i.Config = cfg
	i.Terminal = t
	i.Operation = rl
	return i, nil
}

func (i *Interactive) PreparePrompt(q string) {
	i.SetPrompt("\033[36m" + q + "»\033[0m ")
}

func (i *Interactive) ServerPassword(required bool) ([]byte, string) {
	result := i.PasswordQuestion("Enter server password", required)
	if len(result) == 0 {
		return []byte{}, ""
	}
	secretKey := i.PasswordQuestion("Enter secret phrase for encrypt", true)

	pswd, err := encrypt.Encrypt([]byte(result), secretKey)
	if err != nil {
		println(err.Error())
		return []byte{}, ""
	}

	return pswd, secretKey
}

func FilterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func (i *Interactive) PasswordQuestion(msg string, required bool) (result string) {
	cfg := i.GenPasswordConfig()
	cfg.SetListener(func(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
		i.PreparePrompt(fmt.Sprintf(msg+"(%v)", len(line)))
		i.Refresh()
		return nil, 0, false
	})

	pw, err := i.ReadPasswordWithConfig(cfg)
	if err == readline.ErrInterrupt {
		if len(pw) != 0 {
			return
		} else {
			os.Exit(0)
		}
	} else if err == io.EOF {
		return
	}
	result = strings.TrimSpace(string(pw))
	if result == "" && required {
		result = i.PasswordQuestion(msg, required)
	}
	return
}

func (i *Interactive) Confirm(q string) bool {
	res := i.Question(q, true, "")
	return "yes" == res || "Yes" == res || "y" == res || "Y" == res || "YES" == res
}

func (i *Interactive) Question(q string, required bool, def string) string {
	result := def
	i.PreparePrompt(q)
	var err error
	for {
		result, err = i.ReadlineWithDefault(def)
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
