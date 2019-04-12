package action

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/sakoken/sshh/connector"
	"github.com/sakoken/sshh/encrypt"
	"github.com/sakoken/sshh/interactive"
)

func NewSearch() *Search {
	return &Search{}
}

type Search struct {
	showingHostsList  []*connector.Connector
	readLine          *interactive.Interactive
	positionList      []string
	selectedWithArrow int
	lastSearchKeyWord string
}

func (s *Search) Do(query string) error {
	var err error
	s.readLine, err = interactive.NewEx(&readline.Config{
		Prompt:              "\033[36msshh»\033[0m ",
		InterruptPrompt:     "\n",
		EOFPrompt:           "exit",
		FuncFilterInputRune: s.filterInput,
	})
	if err != nil {
		return err
	}
	defer func(rl *interactive.Interactive) {
		if rl != nil {
			err := rl.Close()
			if err != nil {
				println(err.Error())
			}
		}
	}(s.readLine)
	s.search(query)

	selectedNo, password := s.searchLoop()

	err = s.readLine.Close()
	if err != nil {
		return err
	}
	s.readLine = nil

	if selectedNo >= 0 {
		host := s.showingHostsList[selectedNo]
		connector.SshhData().SetTopPosition(host).Save()
		host.SshConnection(password)
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

func (s *Search) searchLoop() (selectedNo int, password string) {
	for {
		selectedNo = -1
		password = ""

		line, err := s.readLine.Readline()
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
		if line == "add" {
			Add()
			s.search(s.lastSearchKeyWord)
			continue
		} else if line == "exit" {
			os.Exit(0)
		} else if ok, con := s.checkCommand(line, "mod "); ok {
			Modify(con.Position)
			s.search(s.lastSearchKeyWord)
			continue
		} else if ok, con := s.checkCommand(line, "del "); ok {
			connector.SshhData().Delete(con.Position)
			s.search(s.lastSearchKeyWord)
			continue
		} else if ok, con := s.checkCommand(line, "#"); ok {
			key := ""
			if len(con.Password) <= 0 {
				println("No password has been set for this host")
				con.Password, key = s.readLine.Password("Password:", false)
				//host.Key = Question("SSHKey:", true, host.Key)
				connector.SshhData().Save()
			}

			if key == "" {
				key = s.readLine.PasswordQuestion("Enter secret key", true, 16)
			}
			pw, err := encrypt.Decrypt(con.Password, key)
			if err != nil {
				println(err.Error())
				continue
			}
			password = string(pw)
			return
		} else if !strings.HasPrefix(line, "mod ") && !strings.HasPrefix(line, "del ") && !strings.HasPrefix(line, "#") {
			s.lastSearchKeyWord = line
			s.search(line)
			continue
		}
	}

	return
}

func (s *Search) checkCommand(line string, command string) (bool, *connector.Connector) {
	if !strings.HasPrefix(line, command) {
		return false, nil
	}

	line = strings.TrimSpace(line[len(command):])
	if !(len(line) >= 1 && regexp.MustCompile("[0-9]").Match([]byte(line))) {
		println("don't fond the number: " + line)
		return false, nil
	}
	sn, _ := strconv.Atoi(line)
	if len(s.showingHostsList)-1 < sn {
		return false, nil
	}

	return true, s.showingHostsList[sn]
}

func (s *Search) search(keyword string) {
	s.find(keyword)
	s.resetPositionList()
	s.readLine.PrintTable(s.showingHostsList)
}

func (s *Search) resetPositionList() {
	s.selectedWithArrow = -1
	s.positionList = []string{}
	for k := range s.showingHostsList {
		s.positionList = append(s.positionList, fmt.Sprintf("#%d", k))
	}
}

func (s *Search) find(keyword string) {
	var hosts []*connector.Connector
	for _, v := range connector.SshhData().Connectors {
		if strings.Index(v.Host, keyword) >= 0 ||
			strings.Index(v.User, keyword) >= 0 ||
			strings.Index(v.Port, keyword) >= 0 ||
			strings.Index(v.Explanation, keyword) >= 0 {
			hosts = append(hosts, v)
		}
	}
	s.showingHostsList = hosts
}
