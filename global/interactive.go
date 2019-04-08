package global

import (
	"fmt"
	"github.com/chzyer/readline"
	"io"
	"os"
	"strconv"
	"strings"
)

func Password(l *readline.Instance, q string, required bool) ([]byte, string) {
	fmt.Println(q)

	result := PasswordQuestion(l, "Enter password", required, 1000)
	if len(result) == 0 {
		return []byte{}, ""
	}
	secretKey := PasswordQuestion(l, "Enter secret key for encrypt", true, 16)

	pswd, err := Encrypt([]byte(result), secretKey)
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

func Question(l *readline.Instance, q string, required bool, def string) string {
	fmt.Println(q)
	result := def

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

func PrintTable(hosts []*Host) {
	hostLen, portLen, userLen := MaxLength(hosts, 4)

	hostStr := "Host"
	portStr := "Port"
	userStr := "User"

	hostStr = hostStr + strings.Repeat(" ", hostLen-len(hostStr))
	portStr = portStr + strings.Repeat(" ", portLen-len(portStr))
	userStr = userStr + strings.Repeat(" ", userLen-len(userStr))

	lenOfNo := len(strconv.Itoa(len(hosts)))
	noSpace := strings.Repeat(" ", lenOfNo-1)
	println(fmt.Sprintf("[#]"+noSpace+" %s  %s  %s  %s", hostStr, userStr, portStr, "Explanation"))
	noStr := "[%d]"
	for k, v := range hosts {
		noSpace = strings.Repeat(" ",lenOfNo-len(strconv.Itoa(k)))
		host := v.Host + strings.Repeat(" ", hostLen-len(v.Host))
		port := v.Port + strings.Repeat(" ", portLen-len(v.Port))
		user := v.User + strings.Repeat(" ", userLen-len(v.User))
		println(fmt.Sprintf(noStr + noSpace + " %s  %s  %s  %s", k, host, user, port, v.Explanation))
	}
}

func MaxLength(hosts []*Host, def int) (hostLen, portLen, userLen int) {
	hostLen = def
	portLen = def
	userLen = def
	for _, v := range hosts {
		if len(v.Host) > hostLen {
			hostLen = len(v.Host)
		}
		if len(v.Port) > portLen {
			portLen = len(v.Port)
		}
		if len(v.User) > userLen {
			userLen = len(v.User)
		}
	}
	return
}
