# `tray`

`tray` is a minimalistic go-library based on [zserge/tray](https://github.com/zserge/tray)
providing an updatable system tray on all major platforms.

## Getting started

On Mac and Windows no further dependencies are required. Although make sure
`gcc` is available on Windows to be able to use [`cgo`](https://golang.org/cmd/cgo/),
e.g. install [tdm-gcc](https://jmeubank.github.io/tdm-gcc/).
On Linux gtk3 and libappindicator are required with their headers.
Make sure to install them beforehand.

See [example](./example/example.go) for usage instructions.

## Known issues

* No sub menu icons are supported
* Will use a temporary file for the tray icon on Linux and Windows
* To avoid memory leaks and issues with the several system libraries a pre-allocated array is used for the menu entries. This introduces a limit of up to 128 menu entries
* Both the use of static allocation and the requirement to run on the main thread limit each executable to a single system tray menu
* To also provide a GUI consider forking a process, because GUI libraries will require to run in the main thread as well, e.g. [andlabs/ui](github.com/andlabs/ui)

## Acknowledgements

This library would not be possible without [zserge/tray](https://github.com/zserge/tray).

If you are looking for a more mature system tray, [getlantern/systray](github.com/getlantern/systray)
might be a better choice. However it can not properly add/remove menu entries,
which is why this library was created.


