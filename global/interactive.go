package global

import (
	"fmt"
	"github.com/chzyer/readline"
	"io"
	"os"
	"strings"
)

func Password(q string, required bool) []byte {
	fmt.Println(q)

	l, _ := readline.NewEx(&readline.Config{
		Prompt:              "\033[31m»\033[0m ",
		InterruptPrompt:     "\n",
		EOFPrompt:           "exit",
		FuncFilterInputRune: FilterInput,
	})

	defer l.Close()

	result := PasswordQuestion(l, "Enter password", required, 0)
	if len(result) == 0 {
		return []byte{}
	}
	secretKey := PasswordQuestion(l, "Enter secret key for encrypt", true, 16)

	pswd, err := Encrypt([]byte(result), secretKey)
	if err != nil {
		println(err.Error())
	}

	return pswd
}

func FilterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func PasswordQuestion(l *readline.Instance, msg string, required bool, maxLength int) (result string) {
	cfg := l.GenPasswordConfig()
	cfg.SetListener(func(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
		l.SetPrompt(fmt.Sprintf(msg+"(%v): ", len(line)))
		l.Refresh()
		return nil, 0, false
	})

	pw, err := l.ReadPasswordWithConfig(cfg)
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
		result = PasswordQuestion(l, msg, required, maxLength)
	}
	return
}

func Question(q string, required bool, def string) string {
	fmt.Println(q)
	result := def

	l, _ := readline.NewEx(&readline.Config{
		Prompt:              "\033[31m»\033[0m ",
		InterruptPrompt:     "\n",
		EOFPrompt:           "exit",
		FuncFilterInputRune: FilterInput,
	})

	defer l.Close()

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
