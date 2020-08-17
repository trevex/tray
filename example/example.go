package main

import (
	"fmt"

	. "github.com/trevex/tray"
	"github.com/trevex/tray/example/icon"
)

func main() {
	t := Tray{
		Icon: icon.Data,
		Menu: []*TrayMenu{
			{
				Text: "Hello",
				Callback: func(m *TrayMenu) {
					fmt.Println("Hello, tray!")
				},
			},
			{
				Text: "-",
			},
			{
				Text:    "Checked",
				Checked: true,
				Callback: func(m *TrayMenu) {
					m.Checked = !m.Checked
					m.Tray.Update()
				},
			},
			{
				Text:     "Disabled",
				Disabled: true,
			},
			{
				Text: "Sub1",
				SubMenu: []*TrayMenu{
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
				SubMenu: []*TrayMenu{
					{
						Text: "Baz2",
					},
					{
						Text: "Bar2",
					},
				},
			},
			{
				Text: "-",
			},
			{
				Text: "Add",
				Callback: func(m *TrayMenu) {
					m.Tray.Menu = Insert(m.Tray.Menu, len(m.Tray.Menu)-2, &TrayMenu{Text: "Bizzbuzz"})
					m.Tray.Update()
				},
			},
			{
				Text: "Remove",
				Callback: func(m *TrayMenu) {
					if len(m.Tray.Menu) > 9 {
						m.Tray.Menu = Remove(m.Tray.Menu, len(m.Tray.Menu)-2)
						m.Tray.Update()
					}
				},
			},
			{
				Text: "-",
			},
			{
				Text: "Quit",
				Callback: func(m *TrayMenu) {
					m.Tray.Quit()
				},
			},
		},
	}
	t.Run()
}
