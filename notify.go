// +build windows

package wslenv

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	HWND_BROADCAST          = windows.HWND(0xFFFF)
	WM_SETTINGCHANGE        = uint(0x001A)
	SMTO_NORMAL             = uint(0x0000)
	SMTO_ABORTIFHUNG        = uint(0x0002)
	SMTO_NOTIMEOUTIFNOTHUNG = uint(0x0008)
)

func Notify() error {
	user32 := windows.NewLazySystemDLL("user32")
	sendMessageTimeoutW := user32.NewProc("SendMessageTimeoutW")

	// https://docs.microsoft.com/en-us/windows/win32/winmsg/wm-settingchange
	environment, err := windows.UTF16PtrFromString(ENVIRONMENT)
	if err != nil {
		return err
	}

	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendmessagetimeoutw
	ret, _, err := sendMessageTimeoutW.Call(
		uintptr(HWND_BROADCAST),
		uintptr(WM_SETTINGCHANGE),
		uintptr(0),
		uintptr(unsafe.Pointer(environment)),
		uintptr(SMTO_NORMAL|SMTO_ABORTIFHUNG|SMTO_NOTIMEOUTIFNOTHUNG),
		uintptr(1000),
		uintptr(0))
	if uint32(ret) == 0 {
		return err
	}
	return nil
}
