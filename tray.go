package tray

/*
#cgo darwin CFLAGS: -DTRAY_APPKIT=1 -DOBJC_OLD_DISPATCH_PROTOTYPES=1 -x objective-c
#cgo darwin LDFLAGS: -framework Cocoa -framework AppKit
#cgo darwin CFLAGS: -DDARWIN -x objective-c -fobjc-arc
#cgo darwin LDFLAGS: -framework Cocoa
#cgo linux pkg-config: gtk+-3.0 appindicator3-0.1

#include <stdlib.h>
#include "tray.h"

static struct tray_menu menu[128];

static struct tray tray = {
	.menu = menu
};

static struct tray* tray_get() {
	return &tray;
}

static struct tray_menu* tray_menu_get(size_t i) {
	return &menu[i];
}

static void tray_menu_reset() {
	for (int i = 0; i < 64; ++i) {
		menu[i].text = NULL;
		menu[i].cb = NULL;
		menu[i].context = NULL;
	}
}

extern void tray_menu_cb(struct tray_menu *);

static void tray_menu_set_cb(struct tray_menu *s) {
	s->cb = &tray_menu_cb;
}

*/
import "C"
import (
	"fmt"
	"github.com/mattn/go-pointer"
	"runtime"
	"unsafe"
)

//export tray_menu_cb
func tray_menu_cb(cm *C.struct_tray_menu) {
	if cm == nil || cm.context == nil {
		return
	}
	m := pointer.Restore(cm.context).(*TrayMenu)
	if m != nil && m.Callback != nil {
		m.Callback(m)
	}
}

type alloc struct {
	strings []*C.char
}

func (a *alloc) String(str string) *C.char {
	cs := C.CString(str)
	a.strings = append(a.strings, cs)
	return cs
}

func (a *alloc) Free() {
	for _, cs := range a.strings {
		C.free(unsafe.Pointer(cs))
	}
}

type TrayMenu struct {
	Tray     *Tray
	Text     string
	Checked  bool
	Disabled bool
	Callback func(*TrayMenu)
	SubMenu  []*TrayMenu
}

type Tray struct {
	Icon   []byte
	Menu   []*TrayMenu
	allocs []*alloc
}

func (t *Tray) syncC() {
	C.tray_menu_reset()
	a := &alloc{
		strings: []*C.char{},
	}
	t.allocs = append(t.allocs, a)
	ct := C.tray_get()
	ct.icon.data = (*C.char)(unsafe.Pointer(&t.Icon[0]))
	ct.icon.length = (C.int)(len(t.Icon))
	for i, m := range t.Menu {
		m.Tray = t
		cm := C.tray_menu_get(C.size_t(i))
		cm.text = a.String(m.Text)
		cm.checked = boolToInt(m.Checked)
		cm.disabled = boolToInt(m.Disabled)
		cm.context = pointer.Save(m)
		C.tray_menu_set_cb(cm)
	}
	offset := len(t.Menu) + 1
	for i, m := range t.Menu {
		if m.SubMenu == nil || len(m.SubMenu) == 0 {
			continue
		}
		cm := C.tray_menu_get(C.size_t(i))
		cm.submenu = C.tray_menu_get(C.size_t(offset))
		for _, s := range m.SubMenu {
			cs := C.tray_menu_get(C.size_t(offset))
			cs.text = a.String(s.Text)
			cs.checked = boolToInt(s.Checked)
			cs.disabled = boolToInt(s.Disabled)
			cs.context = pointer.Save(s)
			C.tray_menu_set_cb(cs)
			offset += 1
		}
		offset += 1
	}
	if len(t.allocs) > 2 {
		a, t.allocs = t.allocs[0], t.allocs[1:]
		a.Free()
	}
}

func (t *Tray) Run() error {
	runtime.LockOSThread()
	t.syncC()
	if C.tray_init(C.tray_get()) < 0 {
		return fmt.Errorf("tray init failed")
	}
	for C.tray_loop(1) == 0 {
	}
	return nil
}

func (t *Tray) Update() {
	t.syncC()
	C.tray_update(C.tray_get())
}

func (t *Tray) Quit() {
	C.tray_exit()
}

func Insert(a []*TrayMenu, i int, m *TrayMenu) []*TrayMenu {
	return append(a[:i], append([]*TrayMenu{m}, a[i:]...)...)
}

func Remove(a []*TrayMenu, i int) []*TrayMenu {
	return append(a[:i-1], a[i:]...)
}

func boolToInt(b bool) C.int {
	if b {
		return 1
	}
	return 0
}
