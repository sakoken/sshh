package action

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/sakoken/sshh/global"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func NewSeach() *Search {
	return &Search{}
}

type Search struct {
	showingHostsList []*global.Host
}

func (s *Search) Do(query string) error {
	cfg := &readline.Config{
		Prompt:              "\033[31msshh»\033[0m ",
		InterruptPrompt:     "\n",
		EOFPrompt:           "exit",
		FuncFilterInputRune: global.FilterInput,
		AutoComplete:        s.completer(),
	}
	l, _ := readline.NewEx(cfg)
	defer func(l *readline.Instance) {
		if l != nil {
			err := l.Close()
			if err != nil {
				println(err.Error())
			}
		}
	}(l)
	s.showHostsTable(query)
	selectedNo, password := s.searchLoop(l)

	err := l.Close()
	if err != nil {
		println(err.Error())
	}

	if selectedNo >= 0 {
		host := s.showingHostsList[selectedNo]
		global.SshhData.SetTopPosition(host)
		global.SaveJson(global.SshhData)
		s.sshConnection(password, host.Host, host.Port, host.User)
	}

	return nil
}

func (s *Search) searchLoop(l *readline.Instance) (selectedNo int, password string) {
	for {
		selectedNo = -1
		password = ""

		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "exit"):
			os.Exit(0)
		case strings.HasPrefix(line, "#") && len(line) >= 2 && regexp.MustCompile("[0-9]").Match([]byte(line[1:])):
			selectedNo, _ = strconv.Atoi(line[1:])
			if len(s.showingHostsList)-1 < selectedNo {
				continue
			}
			host := s.showingHostsList[selectedNo]
			if len(host.Password) > 0 {
				key := global.PasswordQuestion(l, "Enter secret key", true, 16)
				pw, err := global.Decrypt(host.Password, key)
				if err != nil {
					println(err.Error())
					continue
				}
				password = string(pw)
				return
			}
		default:
			s.showHostsTable(line)
		}
	}

	return
}

func (s *Search) showHostsTable(keyword string) {
	s.find(keyword)
	s.printTable()
}

func (s *Search) printTable() {
	for k, v := range s.showingHostsList {
		println(fmt.Sprintf("[%d] %s %s %s %s", k, v.Host, v.Port, v.User, v.Explanation))
	}
}

func (s *Search) completer() *readline.PrefixCompleter {
	var child []readline.PrefixCompleterInterface
	prefix := readline.NewPrefixCompleter()
	for _, v := range global.SshhData.Hosts {
		child = append(child, readline.PcItem(v.Host))
		for _, v := range strings.Split(v.Explanation, " ") {
			child = append(child, readline.PcItem(v))
		}
	}
	prefix.SetChildren(child)
	return prefix
}

func (s *Search) find(keyword string) {
	var hosts []*global.Host
	for _, v := range global.SshhData.Hosts {
		if strings.Index(v.Host, keyword) >= 0 ||
			strings.Index(v.User, keyword) >= 0 ||
			strings.Index(v.Port, keyword) >= 0 ||
			strings.Index(v.Explanation, keyword) >= 0 {
			hosts = append(hosts, v)
		}
	}
	s.showingHostsList = hosts
}

func (s *Search) sshConnection(password string, host string, port string, user string) {
	ce := func(err error, msg string) {
		if err != nil {
			log.Printf("%s error: %v\n", msg, err)
		}
	}

	var auth []ssh.AuthMethod
	auth = append(auth, ssh.Password(password))

	sshConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host+":"+port, sshConfig)
	ce(err, "dial")

	session, err := client.NewSession()
	ce(err, "new session")
	defer session.Close()

	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		fmt.Println(err)
	}
	defer terminal.Restore(fd, state)

	w, h, err := terminal.GetSize(fd)
	if err != nil {
		fmt.Println(err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	term := os.Getenv("TERM")
	err = session.RequestPty(term, h, w, modes)
	ce(err, "request pty")

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	err = session.Shell()
	ce(err, "start shell")

	err = session.Wait()
	ce(err, "return")
}
