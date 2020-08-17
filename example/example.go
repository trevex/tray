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
					a := m.Tray.Menu
					i := len(a) - 2
					m.Tray.Menu = append(a[:i], append([]*TrayMenu{{Text: fmt.Sprintf("Foo%d", i)}}, a[i:]...)...)
					m.Tray.Update()
				},
			},
			{
				Text: "Remove",
				Callback: func(m *TrayMenu) {
					a := m.Tray.Menu
					if len(a) > 9 {
						i := len(a) - 2
						m.Tray.Menu = append(a[:i-1], a[i:]...)
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
