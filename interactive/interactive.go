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

func (i *Interactive) Password(q string, required bool) ([]byte, string) {
	i.PreparePrompt(q)
	result := i.PasswordQuestion("Enter password", required, 1000)
	if len(result) == 0 {
		return []byte{}, ""
	}
	secretKey := i.PasswordQuestion("Enter secret key for encrypt", true, 16)

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

func (i *Interactive) PasswordQuestion(msg string, required bool, maxLength int) (result string) {
	cfg := i.GenPasswordConfig()
	cfg.SetListener(func(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
		i.PreparePrompt(fmt.Sprintf(msg+"(%v) ", len(line)))
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
	if len(result) >= maxLength && required {
		if len(result) >= maxLength {
			fmt.Printf("Please enter less then %d\n", maxLength)
		}
		result = i.PasswordQuestion(msg, required, maxLength)
	}
	return
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
