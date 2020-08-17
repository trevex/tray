package main

import (
	"fmt"

	"github.com/trevex/tray"
	"github.com/trevex/tray/example/icon"
)

func main() {
	t := tray.Tray{
		Icon: icon.Data,
		Menu: []*tray.Menu{
			{
				Text: "Hello",
				Callback: func(t *tray.Tray, m *tray.Menu) {
					fmt.Println("Hello, tray!")
				},
			},
			tray.Seperator,
			{
				Text:    "Checked",
				Checked: true,
				Callback: func(t *tray.Tray, m *tray.Menu) {
					m.Checked = !m.Checked
					t.Update()
				},
			},
			{
				Text:     "Disabled",
				Disabled: true,
			},
			{
				Text: "Sub1",
				SubMenu: []*tray.Menu{
					{
						Text: "Baz1",
					},
					{
						Text: "Bar1",
					},
				},
			},
			{
				Text: "Sub2",
				SubMenu: []*tray.Menu{
					{
						Text: "Baz2",
					},
					{
						Text: "Bar2",
					},
				},
			},
			tray.Seperator,
			{
				Text: "Add",
				Callback: func(t *tray.Tray, m *tray.Menu) {
					t.Menu = tray.Insert(t.Menu, len(t.Menu)-2, &tray.Menu{Text: "Bizzbuzz"})
					t.Update()
				},
			},
			{
				Text: "Remove",
				Callback: func(t *tray.Tray, m *tray.Menu) {
					if len(t.Menu) > 9 {
						t.Menu = tray.Remove(t.Menu, len(t.Menu)-2)
						t.Update()
					}
				},
			},
			tray.Seperator,
			{
				Text: "Quit",
				Callback: func(t *tray.Tray, m *tray.Menu) {
					t.Quit()
				},
			},
		},
	}
	t.Run()
}
