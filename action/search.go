package action

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/sakoken/sshh/global"
	"github.com/sakoken/sshh/model"
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
	showingHostsList  []*model.Host
	readLine          *readline.Instance
	positionList      []string
	selectedWithArrow int
}

func (s *Search) Do(query string) error {
	cfg := &readline.Config{
		Prompt:              "\033[36msshh»\033[0m ",
		InterruptPrompt:     "\n",
		EOFPrompt:           "exit",
		FuncFilterInputRune: s.filterInput,
		AutoComplete:        s.completer(),
	}
	s.readLine, _ = readline.NewEx(cfg)
	defer func(rl *readline.Instance) {
		if rl != nil {
			err := rl.Close()
			if err != nil {
				println(err.Error())
			}
		}
	}(s.readLine)
	s.showHostsTable(query)
	selectedNo, password := s.searchLoop(s.readLine)

	err := s.readLine.Close()
	if err != nil {
		println(err.Error())
	}

	if selectedNo >= 0 {
		host := s.showingHostsList[selectedNo]
		global.SshhData.SetTopPosition(host)
		global.SaveJson(global.SshhData)
		s.sshConnection(password, host)
	}

	return nil
}

func (s *Search) filterInput(r rune) (rune, bool) {
	switch r {
	case readline.CharPrev:
		s.selectWithArrowUp()
		return r, false
	case readline.CharNext:
		s.selectWithArrowDown()
		return r, false
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func (s *Search) selectWithArrowUp() {
	if s.selectedWithArrow <= 0 {
		return
	}
	s.selectedWithArrow--
	s.readLine.Operation.SetBuffer(s.positionList[s.selectedWithArrow])
}

func (s *Search) selectWithArrowDown() {
	if s.selectedWithArrow >= len(s.positionList)-1 {
		return
	}
	s.selectedWithArrow++
	s.readLine.Operation.SetBuffer(s.positionList[s.selectedWithArrow])
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
		case line == "exit":
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

			println("No password has been set for this host")
		default:
			s.showHostsTable(line)
		}
	}

	return
}

func (s *Search) showHostsTable(keyword string) {
	s.find(keyword)
	s.resetPositionList()
	global.PrintTable(s.showingHostsList)
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

func (s *Search) resetPositionList() {
	s.selectedWithArrow = -1
	s.positionList = []string{}
	for k := range s.showingHostsList {
		s.positionList = append(s.positionList, fmt.Sprintf("#%d", k))
	}
}

func (s *Search) find(keyword string) {
	var hosts []*model.Host
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

func (s *Search) sshConnection(password string, host *model.Host) {
	println(fmt.Sprintf("\033[07m\033[34m%s\033[0m", host.SshCommand()))
	println(fmt.Sprintf("\033[07m\033[34mExplanation: %s\033[0m", host.Explanation))

	var auth []ssh.AuthMethod
	auth = append(auth, ssh.Password(password))

	sshConfig := &ssh.ClientConfig{
		User:            host.User,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host.Host+":"+host.Port, sshConfig)
	if err != nil {
		log.Printf("%s error: %v\n", "dial", err)
		return
	}

	session, err := client.NewSession()
	if err != nil {
		log.Printf("%s error: %v\n", "new session", err)
		return
	}
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
	if err != nil {
		log.Printf("%s error: %v\n", "request pty", err)
		return
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	err = session.Shell()
	if err != nil {
		log.Printf("%s error: %v\n", "start shell", err)
		return
	}

	session.Wait()
}
