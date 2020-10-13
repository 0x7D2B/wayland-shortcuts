package main

import (
	"fmt"
	"os"

	ev "github.com/gvalkov/golang-evdev"
	"gopkg.in/bendahl/uinput.v1"
)

//
// shortcut remapper for recovering mac users. assuming you want ctrl as cmd
//
// Win+(Shift)+Tab           => Next (Previous) Tab
// Alt+Left/Right           => next/previous word
// Ctrl+Left/Right           => beginning/end of line
// Ctrl+Up/Down              => beginning/end of page
// Ctrl/Alt+Backspace/Delete => Delete word/line
// Ctrl+Alt+I                => browser dev tools
// 
// `go get github.com/gvalkov/golang-evdev gopkg.in/bendahl/uinput.v1`
// `go build -ldflags="-s -w"`
//
// run with sudo or better yet set up a dedicated user
// in case everything goes to hell, press LeftCtrl+F1+RightCtrl+F12 to exit
//
// a long time ago this was https://github.com/arnarg/waybind
//

// might need to change this depending on machine
const device_path = "/dev/input/event0"




const KEY_MAX = ev.KEY_MICMUTE

// shift
func LShift(key uint16) bool { return key == ev.KEY_LEFTSHIFT }
func RShift(key uint16) bool { return key == ev.KEY_RIGHTSHIFT }
func Shift(key uint16) bool { return LShift(key) || RShift(key) }
// ctrl
func LCtrl(key uint16) bool { return key == ev.KEY_LEFTCTRL }
func RCtrl(key uint16) bool { return key == ev.KEY_RIGHTCTRL }
func Ctrl(key uint16) bool { return LCtrl(key) || RCtrl(key) }
// super
func LSuper(key uint16) bool { return key == ev.KEY_LEFTMETA }
func RSuper(key uint16) bool { return key == ev.KEY_RIGHTMETA }
func Super(key uint16) bool { return LSuper(key) || RSuper(key) }
// alt
func LAlt(key uint16) bool { return key == ev.KEY_LEFTALT }
func RAlt(key uint16) bool { return key == ev.KEY_RIGHTALT }
func Alt(key uint16) bool { return LAlt(key) || RAlt(key) }

func ShiftSuper(key1, key2 uint16) bool { 
    return (Shift(key1) && Super(key2)) || (Super(key1) && Shift(key2))
}
func ShiftAlt(key1, key2 uint16) bool { 
    return (Shift(key1) && Alt(key2)) || (Alt(key1) && Shift(key2))
}
func CtrlShift(key1, key2 uint16) bool { 
    return (Shift(key1) && Ctrl(key2)) || (Ctrl(key1) && Shift(key2))
}
func CtrlAlt(key1, key2 uint16) bool { 
    return (Ctrl(key1) && Alt(key2)) || (Alt(key1) && Ctrl(key2))
}

func Tab(key uint16) bool { return key == ev.KEY_TAB }

// text keys
func Left(key uint16) bool { return key == ev.KEY_LEFT }
func Right(key uint16) bool { return key == ev.KEY_RIGHT }
func Up(key uint16) bool { return key == ev.KEY_UP }
func Down(key uint16) bool { return key == ev.KEY_DOWN }
func Direction(key uint16) bool { 
    return Left(key) || Right(key) || Up(key) || Down(key) 
}
func Delete(key uint16) bool { return key == ev.KEY_DELETE }
func Backspace(key uint16) bool { return key == ev.KEY_BACKSPACE }
func TextKeys(key uint16) bool { 
    return Direction(key) || Delete(key) || Backspace(key)
}

func main() {
	// Create virtual keyboard
	keyboard, err := uinput.CreateKeyboard("/dev/uinput", []byte("Virtual Keyboard Shortcuts"))
	if err != nil {
		fmt.Printf("Could not create virtual keyboard: %s\n", err)
		os.Exit(1)
	}
	defer keyboard.Close()
    var kbDown = keyboard.KeyDown
    var kbUp = keyboard.KeyUp

	// Open real keyboard
	device, err := ev.Open(device_path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get exclusive access to keyboard
	device.Grab()
	defer device.Release()

    // keep track of keys pressed
	var kPressed [KEY_MAX + 1]bool
    // keep track of keys in the order they're pressed
    var k = make([]uint16, 0, KEY_MAX + 1)

    // shortcuts
    var sSuperTab = false
    var sAltTextKeys = false
    var sCtrlLeft = false
    var sCtrlRight = false
    var sCtrlUp = false
    var sCtrlDown = false
    var sCtrlAltI = false

	for {
		// Read keyboard events
		events, err := device.Read()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

        for _, event := range events {
            // skip event 2
            if event.Type == ev.EV_KEY && event.Value != 2 {
                if event.Value == 1 { // key down
                    k = append(k, event.Code)
                    kPressed[event.Code] = true
                } else if event.Value == 0 { // key up
                    for i, key := range k {
                        if key == event.Code {
                            k = append(k[:i], k[i+1:]...)
                            break
                        }
                    }
                    kPressed[event.Code] = false
                }


                // Killswitch
                if (kPressed[ev.KEY_LEFTCTRL] && kPressed[ev.KEY_RIGHTCTRL] &&
                    kPressed[ev.KEY_F1]       && kPressed[ev.KEY_F12]) {
                    fmt.Println("Exit sequence pressed")
                    os.Exit(0)
                }

                
                // Super+Tab -> Ctrl+Tab
                if (len(k) == 2 &&      Super(k[0])       && Tab(k[1])) || 
                   (len(k) == 3 && ShiftSuper(k[0], k[1]) && Tab(k[2])) {
                    if !sSuperTab {
                        sSuperTab = true
                        // clear super key
                        kbUp(ev.KEY_LEFTMETA)
                        kbUp(ev.KEY_RIGHTMETA)
                        // replace it with ctrl
                        kbDown(ev.KEY_LEFTCTRL)
                    }
                } else if sSuperTab && !Shift(event.Code) { // allow Super+Tab+Shift 
                    sSuperTab = false
                    kbUp(ev.KEY_LEFTCTRL)
                    if kPressed[ev.KEY_LEFTMETA] { kbDown(ev.KEY_LEFTMETA) }
                    if kPressed[ev.KEY_RIGHTMETA] { kbDown(ev.KEY_RIGHTMETA) }


                // Alt+Left/Right/Up/Down/Backspace/Delete -> Ctrl+Left/Right/Up/Down/Backspace/Delete
                } else if (len(k) == 2 &&      Alt(k[0])       && TextKeys(k[1])) || 
                          (len(k) == 3 && ShiftAlt(k[0], k[1]) && TextKeys(k[2])) {
                    if !sAltTextKeys {
                        sAltTextKeys = true
                        // clear alt key
                        kbUp(ev.KEY_LEFTALT)
                        kbUp(ev.KEY_RIGHTALT)
                        // replace it with ctrl
                        kbDown(ev.KEY_LEFTCTRL)
                    }
                } else if sAltTextKeys && !Shift(event.Code) {
                    sAltTextKeys = false
                    kbUp(ev.KEY_LEFTCTRL)
                    if kPressed[ev.KEY_LEFTALT] { kbDown(ev.KEY_LEFTALT) }
                    if kPressed[ev.KEY_RIGHTALT] { kbDown(ev.KEY_RIGHTALT) }


                // Ctrl+Left -> Home
                } else if (len(k) == 2 &&      Ctrl(k[0])       && Left(k[1])) || 
                          (len(k) == 3 && CtrlShift(k[0], k[1]) && Left(k[2])) {
                    if !sCtrlLeft {
                        sCtrlLeft = true
                        // clear all keys
                        kbUp(ev.KEY_LEFTCTRL)
                        kbUp(ev.KEY_RIGHTCTRL)
                        kbUp(ev.KEY_LEFT)
                        // send home
                        kbDown(ev.KEY_HOME)
                        continue // don't pass text keys through
                    }
                } else if sCtrlLeft && !Shift(event.Code) {
                    sCtrlLeft = false
                    kbUp(ev.KEY_HOME)
                    if kPressed[ev.KEY_LEFTCTRL] { kbDown(ev.KEY_LEFTCTRL) }
                    if kPressed[ev.KEY_RIGHTCTRL] { kbDown(ev.KEY_RIGHTCTRL) }
                    if kPressed[ev.KEY_LEFT] { kbDown(ev.KEY_LEFT) }


                // Ctrl+Right -> End
                } else if (len(k) == 2 &&      Ctrl(k[0])       && Right(k[1])) || 
                          (len(k) == 3 && CtrlShift(k[0], k[1]) && Right(k[2])) {
                    if !sCtrlRight {
                        sCtrlRight = true
                        // clear all keys
                        kbUp(ev.KEY_LEFTCTRL)
                        kbUp(ev.KEY_RIGHTCTRL)
                        kbUp(ev.KEY_RIGHT)
                        // send home
                        kbDown(ev.KEY_END)
                        continue // don't pass text keys
                    }
                } else if sCtrlRight && !Shift(event.Code) {
                    sCtrlRight = false
                    kbUp(ev.KEY_END)
                    if kPressed[ev.KEY_LEFTCTRL] { kbDown(ev.KEY_LEFTCTRL) }
                    if kPressed[ev.KEY_RIGHTCTRL] { kbDown(ev.KEY_RIGHTCTRL) }
                    if kPressed[ev.KEY_RIGHT] { kbDown(ev.KEY_RIGHT) }


                // Ctrl+Up -> PgUp
                } else if (len(k) == 2 &&      Ctrl(k[0])       && Up(k[1])) || 
                          (len(k) == 3 && CtrlShift(k[0], k[1]) && Up(k[2])) {
                    if !sCtrlUp {
                        sCtrlUp = true
                        // clear all keys
                        kbUp(ev.KEY_LEFTCTRL)
                        kbUp(ev.KEY_RIGHTCTRL)
                        kbUp(ev.KEY_UP)
                        // send home
                        kbDown(ev.KEY_PAGEUP)
                        continue // don't pass text keys
                    }
                } else if sCtrlUp && !Shift(event.Code) {
                    sCtrlUp = false
                    kbUp(ev.KEY_PAGEUP)
                    if kPressed[ev.KEY_LEFTCTRL] { kbDown(ev.KEY_LEFTCTRL) }
                    if kPressed[ev.KEY_RIGHTCTRL] { kbDown(ev.KEY_RIGHTCTRL) }
                    if kPressed[ev.KEY_UP] { kbDown(ev.KEY_UP) }


                // Ctrl+Down -> PgDown
                } else if (len(k) == 2 &&      Ctrl(k[0])       && Down(k[1])) || 
                          (len(k) == 3 && CtrlShift(k[0], k[1]) && Down(k[2])) {
                    if !sCtrlDown {
                        sCtrlDown = true
                        // clear all keys
                        kbUp(ev.KEY_LEFTCTRL)
                        kbUp(ev.KEY_RIGHTCTRL)
                        kbUp(ev.KEY_DOWN)
                        // send home
                        kbDown(ev.KEY_PAGEDOWN)
                        continue // don't pass text keys
                    }
                } else if sCtrlDown && !Shift(event.Code) {
                    sCtrlDown = false
                    kbUp(ev.KEY_PAGEDOWN)
                    if kPressed[ev.KEY_LEFTCTRL] { kbDown(ev.KEY_LEFTCTRL) }
                    if kPressed[ev.KEY_RIGHTCTRL] { kbDown(ev.KEY_RIGHTCTRL) }
                    if kPressed[ev.KEY_DOWN] { kbDown(ev.KEY_DOWN) }


                // Ctrl+Alt+I -> Ctrl+Shift+I
                } else if len(k) == 3 && CtrlAlt(k[0], k[1]) && k[2] == ev.KEY_I {
                    if !sCtrlAltI {
                        sCtrlAltI = true
                        // clear alt
                        kbUp(ev.KEY_LEFTALT)
                        kbUp(ev.KEY_RIGHTALT)
                        // send shift
                        kbDown(ev.KEY_LEFTSHIFT)
                    }
                } else if sCtrlAltI {
                    sCtrlAltI = false
                    kbUp(ev.KEY_LEFTSHIFT)
                    if kPressed[ev.KEY_LEFTALT] { kbDown(ev.KEY_LEFTALT) }
                    if kPressed[ev.KEY_RIGHTALT] { kbDown(ev.KEY_RIGHTALT) }
                }


                // press the actual key
                if event.Value == 1 {
                    kbDown(int(event.Code))
                } else if event.Value == 0 {
                    kbUp(int(event.Code))
                }
            }
        }
	}
}
