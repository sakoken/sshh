package action

import (
	"github.com/atotto/clipboard"
	"github.com/marcusolsson/tui-go"
	"github.com/sakoken/sshh/global"

	"gopkg.in/urfave/cli.v2"
	"strings"
)

var hosts *tui.Table
var selectIndex int
var showingHostsList []global.Host

func Search(c *cli.Context) error {
	arg := c.Args().First()

	ui, err := tui.New(layout(arg))
	if err != nil {
		return err
	}
	ui.SetKeybinding("Esc", func() { ui.Quit() })
	ui.SetKeybinding("Ctrl+c", func() { ui.Quit() })

	hosts.OnItemActivated(func(table *tui.Table) {
		// 未選択中は何もしない
		if table.Selected() < 0 {
			return
		}
		h := showingHostsList[table.Selected()]

		clipboard.WriteAll(h.SshCommand())
		ui.Quit()
	})

	//初回検索
	find(arg)

	return ui.Run()
}

func layout(searchKeyword string) tui.Widget {

	hosts = tui.NewTable(0, 0)
	hosts.UseRuneEvent = false
	hosts.SetColumnStretch(0, 2)
	hosts.SetColumnStretch(1, 5)
	hosts.SetColumnStretch(2, 1)
	hosts.SetColumnStretch(3, 1)
	hosts.SetSizePolicy(tui.Maximum, tui.Maximum)
	hosts.SetFocused(true)
	hosts.OnSelectionChanged(func(table *tui.Table) {
		// 現在位置を把握するために取得
		selectIndex = table.Selected()
	})

	//初期位置
	selectIndex = -1
	hosts.SetSelected(selectIndex)

	for _, s := range global.SshhData.Hosts {
		hosts.AppendRow(
			tui.NewLabel(s.Host),
			tui.NewLabel(s.Explain),
			tui.NewLabel(s.User),
			tui.NewLabel(s.Port),
		)
	}

	searchWidget := tui.NewEntry()
	searchWidget.SetText(searchKeyword)
	searchWidget.SetFocused(true)
	searchWidget.SetSizePolicy(tui.Expanding, tui.Maximum)
	searchWidget.OnSubmit(func(entry *tui.Entry) {
		find(entry.Text())
	})
	searchWidget.OnChanged(func(entry *tui.Entry) {
		// 何かしら入力されたら選択を初期化
		selectIndex = -1
		hosts.SetSelected(selectIndex)
	})

	searchLabel := tui.NewLabel("QUERY>>")
	searchLabel.SetSizePolicy(tui.Minimum, tui.Maximum)
	searchBox := tui.NewHBox(searchLabel, searchWidget)

	mainFrame := tui.NewVBox(searchBox, hosts, tui.NewSpacer())
	mainFrame.SetSizePolicy(tui.Maximum, tui.Maximum)

	return mainFrame
}

func find(keyword string) {
	//何か選択中の場合は検索を行わない
	if selectIndex >= 0 {
		return
	}

	showingHostsList = []global.Host{}
	for _, v := range global.SshhData.Hosts {
		if strings.Index(v.Host, keyword) >= 0 ||
			strings.Index(v.User, keyword) >= 0 ||
			strings.Index(v.Port, keyword) >= 0 ||
			strings.Index(v.Explain, keyword) >= 0 {
			showingHostsList = append(showingHostsList, v)
		}
	}

	hosts.RemoveRows()
	for _, s := range showingHostsList {
		hosts.AppendRow(
			tui.NewLabel(s.Host),
			tui.NewLabel(s.Explain),
			tui.NewLabel(s.User),
			tui.NewLabel(s.Port),
		)
	}
}
