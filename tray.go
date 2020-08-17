package tray

/*
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
	"io/ioutil"
	"os"
	"runtime"
	"unsafe"

	"github.com/mattn/go-pointer"
)

var Seperator = &Menu{Text: "-"}

//export tray_menu_cb
func tray_menu_cb(cm *C.struct_tray_menu) {
	if cm == nil || cm.context == nil {
		return
	}
	m := pointer.Restore(cm.context).(*Menu)
	if m != nil && m.Callback != nil {
		m.Callback(m.tray, m)
	}
}

type alloc struct {
	strings []*C.char
	ptrs    []unsafe.Pointer
}

func (a *alloc) String(str string) *C.char {
	cs := C.CString(str)
	a.strings = append(a.strings, cs)
	return cs
}

func (a *alloc) Pointer(o interface{}) unsafe.Pointer {
	ptr := pointer.Save(o)
	a.ptrs = append(a.ptrs, ptr)
	return ptr
}

func (a *alloc) Free() {
	for _, cs := range a.strings {
		C.free(unsafe.Pointer(cs))
	}
	for _, ptr := range a.ptrs {
		pointer.Unref(ptr)
	}
}

type Menu struct {
	tray     *Tray
	Text     string
	Checked  bool
	Disabled bool
	Callback func(*Tray, *Menu)
	SubMenu  []*Menu
}

type Tray struct {
	Icon   []byte
	Menu   []*Menu
	allocs []*alloc
}

func (t *Tray) syncC() {
	C.tray_menu_reset()
	a := &alloc{
		strings: []*C.char{},
		ptrs:    []unsafe.Pointer{},
	}
	t.allocs = append(t.allocs, a)
	ct := C.tray_get()
	if runtime.GOOS != "linux" {
		ct.icon.data = (*C.char)(unsafe.Pointer(&t.Icon[0]))
		ct.icon.length = (C.int)(len(t.Icon))
	} else {
		// On linux instead we store the content of the icon the temp file,
		// which is created in the .Run method and cleaned up via .Quit.
		err := ioutil.WriteFile(C.GoString(ct.icon.data), t.Icon, 0644)
		if err != nil {
			panic(err)
		}
	}
	for i, m := range t.Menu {
		m.tray = t
		trayMenuSet(a, i, m)
	}
	offset := len(t.Menu) + 1
	for i, m := range t.Menu {
		if m.SubMenu == nil || len(m.SubMenu) == 0 {
			continue
		}
		cm := C.tray_menu_get(C.size_t(i))
		cm.submenu = C.tray_menu_get(C.size_t(offset))
		for _, s := range m.SubMenu {
			s.tray = t
			trayMenuSet(a, offset, s)
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
	if runtime.GOOS == "linux" {
		// On linux we'll have to store the icon in a tempfile, so let's
		// create the file. This file will be written to in .syncC and it will
		// be cleaned up by .Quit.
		ct := C.tray_get()
		tmp, err := ioutil.TempFile("", "icon")
		if err != nil {
			panic(err)
		}
		ct.icon.data = C.CString(tmp.Name())
	}
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
	if runtime.GOOS == "linux" {
		// Let's cleanup or temporary file.
		ct := C.tray_get()
		os.Remove(C.GoString(ct.icon.data))
	}
	C.tray_exit()
}

func Insert(a []*Menu, i int, m *Menu) []*Menu {
	return append(a[:i], append([]*Menu{m}, a[i:]...)...)
}

func Remove(a []*Menu, i int) []*Menu {
	return append(a[:i-1], a[i:]...)
}

func trayMenuSet(a *alloc, i int, m *Menu) {
	cm := C.tray_menu_get(C.size_t(i))
	cm.text = a.String(m.Text)
	cm.checked = boolToInt(m.Checked)
	cm.disabled = boolToInt(m.Disabled)
	cm.context = a.Pointer(m)
	C.tray_menu_set_cb(cm)
}

func boolToInt(b bool) C.int {
	if b {
		return 1
	}
	return 0
}
