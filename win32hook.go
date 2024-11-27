// CREDIT Jimmy Obonyo Abor
// https://stackoverflow.com/questions/41559026/go-lang-set-windows-events-hooks

package main

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"

	"github.com/gonutz/w32/v2"
	"golang.org/x/sys/windows"
)

var (
	modkernel32          = windows.NewLazyDLL("kernel32.dll")
	procSetWinEventHook  = user32.NewProc("SetWinEventHook")
	procUnhookWinEvent   = user32.NewProc("UnhookWinEvent")
	procGetMessage       = user32.NewProc("GetMessageW")
	procTranslateMessage = user32.NewProc("TranslateMessage")
	procDispatchMessage  = user32.NewProc("DispatchMessageW")
	procGetModuleHandle  = modkernel32.NewProc("GetModuleHandleW")
	winEvHook            HWINEVENTHOOK
)

type WINEVENTPROC func(hWinEventHook HWINEVENTHOOK, event uint32, hwnd HWND, idObject int32, idChild int32, idEventThread uint32, dwmsEventTime uint32) uintptr

type (
	HANDLE        uintptr
	HINSTANCE     HANDLE
	HHOOK         HANDLE
	HMODULE       HANDLE
	HWINEVENTHOOK HANDLE
	DWORD         uint32
	INT           int
	WPARAM        uintptr
	LPARAM        uintptr
	LRESULT       uintptr
	HWND          HANDLE
	UINT          uint32
	BOOL          int32
	ULONG_PTR     uintptr
	LONG          int32
	LPWSTR        *WCHAR
	WCHAR         uint16
)

type POINT struct {
	X, Y int32
}

type MSG struct {
	Hwnd    HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

const (
	EVENT_SYSTEM_FOREGROUND = 3
	WINEVENT_OUTOFCONTEXT   = 0
	WINEVENT_INCONTEXT      = 4
	WINEVENT_SKIPOWNPROCESS = 2
	WINEVENT_SKIPOWNTHREAD  = 1
)

// Starts a windows hook for foreground window changed event
// Emits Event to CharSelectedEvent
func (a *App) foregroundWindowsHook() {
	// log.Println("starting")
	// might block something?
	// hinst := GetModuleHandle("")
	// fmt.Println(hinst)

	eventChannel := make(chan w32.HWND)

	go func() {
		for hwnd := range eventChannel {
			a.CharSelectedEvent(hwnd)
		}
	}()

	ActiveWinEventHook := func(hWinEventHook HWINEVENTHOOK, event uint32, hwnd HWND, idObject int32, idChild int32, idEventThread uint32, dwmsEventTime uint32) uintptr {
		w32HWND := w32.HWND(hwnd)
		processName, _ := w32.GetClassName(w32HWND)
		if processName == "UnityWndClass" {
			// Send the hwnd to the channel
			eventChannel <- w32HWND
			// log.Printf("windowActivated : %s\n", w32.GetWindowText(w32HWND))
		}
		return 0
	}

	winEvHook = SetWinEventHook(EVENT_SYSTEM_FOREGROUND, EVENT_SYSTEM_FOREGROUND, 0, ActiveWinEventHook, 0, 0, WINEVENT_OUTOFCONTEXT|WINEVENT_SKIPOWNPROCESS)
	// log.Println("Windows Event Hook: ")
	log.Println("Windows Event Hook: ", winEvHook)

	for {

		var msg MSG
		if m := GetMessage(&msg, 0, 0, 0); m != 0 {
			TranslateMessage(&msg)
			DispatchMessage(&msg)
		}
	}
}

func SetWinEventHook(eventMin DWORD, eventMax DWORD, hmodWinEventProc HMODULE, pfnWinEventProc WINEVENTPROC, idProcess DWORD, idThread DWORD, dwFlags DWORD) HWINEVENTHOOK {
	// log.Println("procSetWinEventHook S")
	pfnWinEventProcCallback := syscall.NewCallback(pfnWinEventProc)
	ret, _, err := procSetWinEventHook.Call(
		uintptr(eventMin),
		uintptr(eventMax),
		uintptr(hmodWinEventProc),
		pfnWinEventProcCallback,
		uintptr(idProcess),
		uintptr(idThread),
		uintptr(dwFlags),
	)
	if err != nil {
		log.Printf("%#v", err)
	}
	// log.Printf("%#v", ret2)
	// log.Println("procSetWinEventHook E")
	return HWINEVENTHOOK(ret)
}

func UnhookWinEvent(hWinEventHook HWINEVENTHOOK) bool {
	// fmt.Println("UnhookWinEvent called")
	ret, _, err := procUnhookWinEvent.Call(
		uintptr(hWinEventHook),
	)
	if err != nil {
		fmt.Printf("Error in UnhookWinEvent: %v\n", err)
	}

	return ret != 0
}

func GetModuleHandle(modulename string) HINSTANCE {
	var mn uintptr
	if modulename == "" {
		mn = 0
	} else {
		mn = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(modulename)))
	}
	ret, _, _ := procGetModuleHandle.Call(mn)
	return HINSTANCE(ret)
}

func GetMessage(msg *MSG, hwnd HWND, msgFilterMin UINT, msgFilterMax UINT) int {
	ret, _, _ := procGetMessage.Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax))

	return int(ret)
}

func TranslateMessage(msg *MSG) bool {
	ret, _, _ := procTranslateMessage.Call(
		uintptr(unsafe.Pointer(msg)))
	return ret != 0
}

func DispatchMessage(msg *MSG) uintptr {
	ret, _, _ := procDispatchMessage.Call(
		uintptr(unsafe.Pointer(msg)))
	return ret
}
