package main

import (
	"context"
	"fmt"
	"path/filepath"
	"syscall"
	"unsafe"

	"github.com/go-vgo/robotgo"
	"github.com/gonutz/w32/v2"
	hook "github.com/robotn/gohook"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// function to iterate saved chars and find their window hwnd

// Fetch Windows and store Name - Class - Handle
// Move button for re order
// Save to json
// Ask for name
// name :
// - 1 : name: ddd, class: ddd, order: 1

// Button to activate / desactivate
// Button check if window active exe is dofus and processName is unity
// Listener
// Save on exit

type Account struct {
	Name  string `json:"name"`
	Class string `json:"class"`
	Order int    `json:"order"`
}

// App struct
type App struct {
	ctx          context.Context
	DofusWindows []WindowInfo
}

type WindowInfo struct {
	Title string `json:"title"`
	Hwnd  uint64 `json:"hwnd"` // Use uint64 instead of uintptr
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) Add() {
	chanHook := hook.Start()
	defer hook.End()

	for ev := range chanHook {
		if ev.Rawcode == 114 {
			// need to store hwnd when launched or refresh
			// load hwnd from list
			// use char name to know if window is dofus instead of exe
			activeRobotWindow := robotgo.GetHWND()
			exeName, _ := GetExecutableName(w32.HWND(activeRobotWindow))
			if exeName == "Dofus.exe" {
				a.SetWindowForeground(a.DofusWindows[0].Hwnd)
				// 	}
			}
		}
	}
	// hook.AddEvent(hook.KeyDown, []string{"61"}, func(e hook.Event) {
	// hookTest := hook.AddEvent("f3")
	// if hookTest {
	// 	runtime.LogPrint(a.ctx, "f3 pressed")
	// }

	// hook.Register(hook.KeyDown, []string{"f3"}, func(e hook.Event) {
	// 	runtime.LogPrintf(a.ctx, "keychar %v", e.Keychar)
	// })

	// {
	// 	activeRobotWindow := robotgo.GetHWND()
	// 	runtime.LogPrintf(a.ctx, "Key pressed: %v\n", e)
	// 	// runtime.LogPrintf(a.ctx, "Active ROBOT window hwnd: %d", int32(activeRobotWindow))
	// 	runtime.LogPrint(a.ctx, "test")
	// 	exeName, _ := GetExecutableName(w32.HWND(activeRobotWindow))
	// 	runtime.LogPrintf(a.ctx, "Active window exe: %s", exeName)
	// 	if exeName == "Dofus.exe" {
	// 		a.SetWindowForeground(a.DofusWindows[0].Hwnd)
	// 	}
	// })

	// hook: 2024-11-16 20:30:21.1541953 +0100 CET m=+6.623024101 - Event: {Kind: MouseHold, Button: 1, X: -490, Y: 901, Clicks: 1}
	// hook: 2024-11-16 20:30:23.4710445 +0100 CET m=+8.939873301 - Event: {Kind: MouseHold, Button: 4, X: -490, Y: 901, Clicks: 1}
	// hook: 2024-11-16 20:30:26.7923629 +0100 CET m=+12.261191701 - Event: {Kind: KeyHold, Rawcode: 114, Keychar: 65535}
	// s := hook.Start()
	// <-hook.Process(s)
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods

var lastCharActive string

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.DofusWindows = []WindowInfo{}
	w32.EnumWindows(func(hwnd w32.HWND) bool {
		return EnumWindowsCallback(ctx, hwnd, a)
	})
	a.Add()
	// runtime.LogPrintf(ctx, "Dofus pid: %d", pid)
}

func EnumWindowsCallback(ctx context.Context, hwnd w32.HWND, a *App) bool {
	// Get the window title
	title := w32.GetWindowText(hwnd)
	processName, _ := w32.GetClassName(hwnd)
	exeName, _ := GetExecutableName(hwnd)

	// a.DofusWindows = []WindowInfo{}

	if exeName == "Dofus.exe" && processName == "UnityWndClass" {
		runtime.LogPrintf(ctx, "title is %s, processName is: %s, exeName is : %s\n", title, processName, exeName)
		a.DofusWindows = append(a.DofusWindows, WindowInfo{Title: title, Hwnd: uint64(hwnd)})
		if len(lastCharActive) == 0 {
			lastCharActive = a.DofusWindows[0].Title
		}
	}
	return true // Continue enumeration
}

var (
	modKernel32                   = syscall.NewLazyDLL("kernel32.dll")
	procQueryFullProcessImageName = modKernel32.NewProc("QueryFullProcessImageNameW")
)

const (
	PROCESS_QUERY_LIMITED_INFORMATION = 0x1000
)

func GetExecutableName(hwnd w32.HWND) (string, error) {
	// Get the process ID
	_, pid := w32.GetWindowThreadProcessId(hwnd)

	// Open the process
	handle, _, _ := syscall.NewLazyDLL("kernel32.dll").
		NewProc("OpenProcess").
		Call(PROCESS_QUERY_LIMITED_INFORMATION, 0, uintptr(pid))
	if handle == 0 {
		return "", fmt.Errorf("unable to open process for PID %d", pid)
	}
	defer syscall.CloseHandle(syscall.Handle(handle))

	// Query the executable name
	buffer := make([]uint16, syscall.MAX_PATH)
	size := uint32(len(buffer))
	ret, _, err := procQueryFullProcessImageName.Call(
		handle,
		0,
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&size)),
	)
	if ret == 0 {
		return "", fmt.Errorf("error retrieving executable name: %v", err)
	}

	// Convert UTF-16 buffer to a Go string
	return filepath.Base(syscall.UTF16ToString(buffer[:size])), nil
}

// // Greet returns a greeting for the given name
// func (a *App) Greet(p Person) string {
// 	return fmt.Sprintf("Hello %s (Age: %d), It's show time!", p.Name, p.Age)
// }

// func (a *App) TestReturn() string {
// 	if len(a.DofusWindows) > 0 {
// 		runtime.LogPrintf(a.ctx, "Found %d windows", len(a.DofusWindows))
// 		for _, window := range a.DofusWindows {
// 			runtime.LogPrintf(a.ctx, "Window Title: %s, Hwnd: %v", window.Title, window.Hwnd)
// 		}
// 	}
// 	return "test"
// }

func (a *App) SetWindowForeground(hwnd uint64) {
	hwndW32 := w32.HWND(hwnd)
	w32.SetForegroundWindow(hwndW32)
}

func (a *App) GetDofusWindows() []WindowInfo {
	if len(a.DofusWindows) > 0 {
		return a.DofusWindows
	}
	return nil
}

func (a *App) UpdateDofusWindows() []WindowInfo {
	a.DofusWindows = []WindowInfo{}
	w32.EnumWindows(func(hwnd w32.HWND) bool {
		return EnumWindowsCallback(a.ctx, hwnd, a)
	})
	if len(a.DofusWindows) > 0 {
		return a.DofusWindows
	}
	return nil
}
