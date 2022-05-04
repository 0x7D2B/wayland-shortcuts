# wayland-shortcuts

wayland shortcut remapper for recovering mac users.

1. `Win`+(`Shift`)+`Tab`             => Next (Previous) Tab
2. `Alt`+`Left`/`Right`              => Next/Previous Word
3. `Ctrl`+`Left`/`Right`             => Beginning/End of Line
4. `Ctrl`+`Up`/`Down`                => Beginning/End of Page
5. `Ctrl`/`Alt`+`Backspace`/`Delete` => Delete Word/Line
6. `Ctrl`+`Alt`+`I`                  => Browser Dev Tools

in case everything goes to hell, press `LeftCtrl`+`F1`+`F12` to exit

## Build

`go get`
`go build -ldflags="-s -w"`
                                                                         
## Use

you can run this as root but that's very scary, do you trust people on the internet not to keylog you???

try [SETUP.md](SETUP.md) maybe it can install it as a systemd service for you (also send me a pr to do this properly thx)

----

a long time ago this was https://github.com/arnarg/waybind
