package controller

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/sakoken/sshh/action"
	"github.com/sakoken/sshh/config"
	"github.com/sakoken/sshh/connector"
	"github.com/sakoken/sshh/encrypt"
	"github.com/sakoken/sshh/interactive"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func NewController() *Controller {
	return &Controller{}
}

type Controller struct {
	showingCollection connector.ConnectorCollection
	readLine          *interactive.Interactive
	positionList      []string
	selectedWithArrow int
	lastSearchKeyWord string
}

func (s *Controller) Do(query string) error {
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

	con, password := s.loop()

	err = s.readLine.Close()
	if err != nil {
		return err
	}
	s.readLine = nil

	if con != nil && password != "" {
		config.SshhData().ToTopPosition(con)
		config.SshhData().Save()
		con.SshConnection(password)
	}

	return nil
}

func (s *Controller) loop() (con *connector.Connector, password string) {
	for {
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

		ok := false
		line = strings.TrimSpace(line)
		if line == "add" {
			action.Add()
			s.search(s.lastSearchKeyWord)
			continue
		} else if line == "exit" {
			os.Exit(0)
		} else if ok, con = s.checkCommand(line, "mod "); ok {
			action.Modify(con.Position)
			s.search(s.lastSearchKeyWord)
			continue
		} else if ok, con = s.checkCommand(line, "del "); ok {
			action.Delete(con.Position)
			s.search(s.lastSearchKeyWord)
			continue
		} else if ok, con = s.checkCommand(line, "#"); ok {
			key := ""
			if len(con.Password) <= 0 {
				println("No password has been set for this host")
				con.Password, key = s.readLine.ServerPassword(false)
				//host.Key = Question("SSHKey:", true, host.Key)
				config.SshhData().Save()
			}

			if key == "" {
				key = s.readLine.PasswordQuestion("Enter secret phrase", true)
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
func (s *Controller) checkCommand(line string, command string) (bool, *connector.Connector) {
	if !strings.HasPrefix(line, command) {
		return false, nil
	}

	line = strings.TrimSpace(line[len(command):])
	if !(len(line) >= 1 && regexp.MustCompile("[0-9]").Match([]byte(line))) {
		println("don't fond the number: " + line)
		return false, nil
	}
	sn, _ := strconv.Atoi(line)
	if s.showingCollection.Count()-1 < sn {
		return false, nil
	}

	return true, s.showingCollection.Index(sn)
}

func (s *Controller) filterInput(r rune) (rune, bool) {
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

func (s *Controller) selectWithArrowUp() {
	if s.selectedWithArrow <= 0 {
		return
	}
	s.selectedWithArrow--
	s.readLine.Operation.SetBuffer(s.positionList[s.selectedWithArrow])
}

func (s *Controller) selectWithArrowDown() {
	if s.selectedWithArrow >= len(s.positionList)-1 {
		return
	}
	s.selectedWithArrow++
	s.readLine.Operation.SetBuffer(s.positionList[s.selectedWithArrow])
}

func (s *Controller) search(keyword string) {
	s.showingCollection.Connectors = action.Search(keyword)
	s.resetPositionList()
	s.showingCollection.PrintTable()
}

func (s *Controller) resetPositionList() {
	s.selectedWithArrow = -1
	s.positionList = []string{}
	for k := range s.showingCollection.Connectors {
		s.positionList = append(s.positionList, fmt.Sprintf("#%d", k))
	}
}
