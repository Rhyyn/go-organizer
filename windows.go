package main

import (
	"context"
	"strings"
	"syscall"
	"time"

	"github.com/gonutz/w32/v2"
	"github.com/lxn/win"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/sys/windows"
)

// This is basically a port of winActivate from AHK :
// - https://github.com/AutoHotkey/AutoHotkey/blob/581114c1c7bb3890ff61cf5f6e1f1201cd8c8b78/source/window.cpp
// Slimmer version because we don't need as many checks
// All credit goes to their contributors

var (
	procAttachThreadInput = user32.NewProc("AttachThreadInput")
	procBringWindowToTop  = user32.NewProc("BringWindowToTop")
	isHungAppWindow       = user32.NewProc("IsHungAppWindow")
	user32                = syscall.NewLazyDLL("user32.dll")
	// enumDisplayMonitors   = user32.NewProc("EnumDisplayMonitors")
	// monitorFromWindow     = user32.NewProc("MonitorFromWindow")
	ATTEMPT_SET_FORE bool
	// basePosition             WindowPosition
	// MONITOR_DEFAULTTOPRIMARY DWORD
	// monitor                  HMONITOR
)

// type WindowPosition struct {
// 	X int
// 	Y int
// }
// type RECT struct {
// 	Left, Top, Right, Bottom int32
// }
// type (
// 	HMONITOR HANDLE
// 	HDC      HANDLE
// )

// func (a *App) monitorEnumProc(hMonitor HMONITOR, hdcMonitor HDC, lprcMonitor *RECT, dwData LPARAM) uintptr {
// 	fmt.Printf("Monitor Handle: %v\n", hMonitor)
// 	fmt.Printf("Monitor Rect: Left=%d, Top=%d, Right=%d, Bottom=%d\n",
// 		lprcMonitor.Left, lprcMonitor.Top, lprcMonitor.Right, lprcMonitor.Bottom)

// 	return 1
// }

// func EnumDisplayMonitors(hdc HDC, clip *RECT, lpfnEnum, data uintptr) error {
// 	ret, _, _ := enumDisplayMonitors.Call(
// 		uintptr(hdc),
// 		uintptr(unsafe.Pointer(clip)),
// 		lpfnEnum,
// 		data,
// 	)
// 	if ret == 0 {
// 		return fmt.Errorf("EnumDisplayMonitors returned FALSE")
// 	}
// 	return nil
// }

// func (a *App) getAllMonitors() {
// 	enumProc := syscall.NewCallback(a.monitorEnumProc)
// 	err := EnumDisplayMonitors(0, nil, enumProc, 0)
// 	if err != nil {
// 		log.Fatalf("Failed to enumerate monitors: %v", err)
// 	}
// }

// func MonitorFromWindow(hwnd HWND, dwFlags uint32) HMONITOR {
// 	ret, _, _ := monitorFromWindow.Call(monitorFromWindow.Addr(), 2,
// 		uintptr(hwnd),
// 		uintptr(dwFlags),
// 		0)

// 	return HMONITOR(ret)
// }

// // This needs to trigger at startup
// // Need to check DPI
// func (a *App) handleWindowPosition() {
// 	a.getAllMonitors()

// 	// We save the base position on startup
// 	basePosition.X, basePosition.Y = runtime.WindowGetPosition(a.ctx)
// 	// save current monitor that window is on, if no monitor, sets to primary
// 	hwnd := HWND(w32.GetCurrentProcess()) // could do native, but cba
// 	monitor = MonitorFromWindow(hwnd, 3)

// 	// func to compare to config saved
// 	// if != enumMonitors
// 	// if monitor handle not present in monitors
// 	// sets to primary and save
// 	// MonitorHandle = x
// 	// FullWindowPosition = x, y
// 	// OverlayWindowPosition = x, y
// 	// else check for saved pos
// 	// if pos not of bound
// 	// sets pos check bounds
// 	// if out set to center or bound?
// }

func (a *App) UpdateTemporaryDofusWindows(tempChars []WindowInfo) {
	if tempChars != nil || len(tempChars) != 0 {
		a.DofusWindows = tempChars
	} else {
		runtime.LogErrorf(a.ctx, "error while updating temporary chars to a.DofusWindows %v", tempChars)
	}
}

func (a *App) ActivateCharacter(characterName string) {
	isDofus, _ := a.IsWindowDofus()
	if isDofus {
		for _, char := range a.DofusWindows {
			if char.CharacterName == characterName {
				a.WinActivate(w32.HWND(char.Hwnd))
			}
		}
	}
}

func (a *App) ActivateNextChar() {
	isDofus, index := a.IsWindowDofus()
	if isDofus {
		nextIndex := (index + 1) % len(a.DofusWindows)
		nextWindow := win.HWND(a.DofusWindows[nextIndex].Hwnd)
		a.WinActivate(w32.HWND(nextWindow))
	}
}

func (a *App) ActivatePreviousChar() {
	isDofus, index := a.IsWindowDofus()
	if isDofus {
		nextIndex := (index - 1 + len(a.DofusWindows)) % len(a.DofusWindows)
		nextWindow := win.HWND(a.DofusWindows[nextIndex].Hwnd)
		a.WinActivate(w32.HWND(nextWindow))
	}
}

func (a *App) WinActivate(targetWindow w32.HWND) w32.HWND {
	origForegroundWindow := a.getForegroundWindow()
	// Check if our window is valid, returns original if not
	if !a.isWindowValid(targetWindow) {
		// runtime.LogPrintf(a.ctx, "Target window is not valid. %v", targetWindow)
		return origForegroundWindow
	}
	// fmt.Printf("targetWindow : %v\n origForegroundWindow: %v\n", targetWindow, origForegroundWindow)

	return a.setForegroundWindowEx(targetWindow, origForegroundWindow)
}

// Fetch Windows to see if any new Dofus windows appeared
func (a *App) refreshAndUpdateCharacterList(exists bool) {
	a.DofusWindows = []WindowInfo{}

	// Loop through windows and populate our array
	w32.EnumWindows(func(hwnd w32.HWND) bool {
		return EnumWindowsCallback(a.ctx, hwnd, a)
	})

	// runtime.LogPrintf(a.ctx, "Looped through Windows and inside refreshAndUpdateCharacterList with exists : %t", exists)

	// This stinks
	if !exists {
		a.SaveCharacterList(a.DofusWindows)
	}

	// Should order before updating front
	a.UpdateDofusWindowsOrder(a.DofusWindows)

	// runtime.LogPrintf(a.ctx, "end of refresh updating Dofus windows")
}

// Iterate through all active Windows and populate a.DofusWindows with them
func EnumWindowsCallback(ctx context.Context, hwnd w32.HWND, a *App) bool {
	// Get the window title
	title := w32.GetWindowText(hwnd)
	processName, _ := w32.GetClassName(hwnd)
	exeName, _ := GetExecutableName(hwnd)

	if exeName == "Dofus.exe" && processName == "UnityWndClass" {
		characterName, class := parseTitleComponents(title)

		keybind := ""
		for _, key := range keybindMap {
			if key.Action == characterName {
				keybind = key.KeyName
			}
		}

		a.DofusWindows = append(a.DofusWindows, WindowInfo{
			Title:         title,
			Hwnd:          uint64(hwnd),
			CharacterName: characterName,
			Class:         class,
			Keybind:       keybind,
		})

		// runtime.LogPrintf(ctx, "Processed window: %s ", title)
	}
	return true
}

// Updates the order of the list of Characters
func (a *App) UpdateDofusWindowsOrder(loggedInCharacters []WindowInfo) ([]WindowInfo, error) {
	if len(loggedInCharacters) == 0 {
		return a.DofusWindows, nil
	}

	for i := range loggedInCharacters {
		char := &loggedInCharacters[i] // reference

		for _, key := range keybindMap {
			if key.Action == char.CharacterName {
				char.Keybind = key.KeyName
			}
		}
	}

	// no error handling because i dont have time
	iniFile, _, _ := loadINIFile(charactersFilePath)

	// Load saved character names from the INI file
	savedOrder, err := a.loadCharacterList(iniFile)
	if err != nil {
		runtime.LogError(a.ctx, "Error loading character list")
		return nil, err
	}

	// array of known char from our saved order
	var newOrderKnown []WindowInfo
	// array of unknown char from our saved order
	var newOrderUnknown []WindowInfo

	loggedInMap := make(map[string]WindowInfo)

	for _, char := range loggedInCharacters {
		loggedInMap[char.CharacterName] = char
	}

	processed := make(map[string]bool)
	processedHWND := make(map[uint64]bool)

	for _, loggedChar := range loggedInCharacters {
		if !strings.Contains(loggedChar.CharacterName, "Dofus") {
			for _, savedChar := range savedOrder {
				if _, exists := processed[savedChar]; exists {
					continue
				}

				if loggedChar, exists := loggedInMap[savedChar]; exists {
					newOrderKnown = append(newOrderKnown, loggedChar)
					processed[savedChar] = true
					processedHWND[loggedChar.Hwnd] = true
				} else {
					processedHWND[loggedChar.Hwnd] = true
					processed[savedChar] = true
				}
				processedHWND[loggedChar.Hwnd] = true
				processed[savedChar] = true
			}

			for _, loggedChar := range loggedInCharacters {
				if _, exists := processed[loggedChar.CharacterName]; !exists {
					newOrderUnknown = append(newOrderUnknown, loggedChar)
					processed[loggedChar.CharacterName] = true
					processedHWND[loggedChar.Hwnd] = true
				}
			}
		} else {
			if _, exists := processedHWND[loggedChar.Hwnd]; !exists {
				processed[loggedChar.CharacterName] = true
				processedHWND[loggedChar.Hwnd] = true
				newOrderUnknown = append(newOrderUnknown, loggedChar)
			}
		}
	}

	// fmt.Printf("newOrderUnknown : %v\n", newOrderUnknown)
	// fmt.Printf("newOrderKnown : %v\n", newOrderKnown)
	newOrderKnown = append(newOrderKnown, newOrderUnknown...)
	// fmt.Printf("newOrderKnown after combining : %v\n", newOrderKnown)

	a.DofusWindows = newOrderKnown

	// if we set [0] but user activate manually other char it fucks everything
	// a.WinActivate(w32.HWND(a.DofusWindows[0].Hwnd))

	return newOrderKnown, nil
}

// Used by the frontend to fetch the array, I think it might be useless now?
func (a *App) GetDofusWindows() []WindowInfo {
	_, err, exists := loadINIFile(charactersFilePath)
	if err != nil {
		runtime.LogError(a.ctx, "Error with the ini file")
	}

	a.refreshAndUpdateCharacterList(exists)

	if len(a.DofusWindows) > 0 {
		return a.DofusWindows
	}
	return nil
}

// Check if user forground window is a Dofus window ->
// return true and its index in list
// || return false and 0 if its not Dofus
func (a *App) IsWindowDofus() (bool, int) {
	activeWindowHandle := w32.GetForegroundWindow()
	if activeWindowHandle == 0 {
		return false, 0
	}

	for i, window := range a.DofusWindows {
		if window.Hwnd == uint64(activeWindowHandle) {
			return true, i
		}
	}

	return false, 0
}

// checks if our window can be activated
func (a *App) isWindowValid(targetWindow w32.HWND) bool {
	isValid := w32.IsWindow(targetWindow)
	return isValid
}

// gets the current foreground window
func (a *App) getForegroundWindow() w32.HWND {
	return w32.GetForegroundWindow()
}

// gets the tread id of a window
func (a *App) getWindowThreadProcessId(window w32.HWND) w32.HANDLE {
	foreThread, _ := w32.GetWindowThreadProcessId(window)
	return foreThread
}

// Checks if the window is unresponsive
func (a *App) isWindowHung(hwnd uintptr) bool {
	ret, _, _ := isHungAppWindow.Call(hwnd)
	return ret != 0
}

// Checks if window is minimized
const WS_MINIMIZE = 0x20000000

func (a *App) isWindowMinimized(hwnd w32.HWND) bool {
	style := w32.GetWindowLong(hwnd, w32.GWL_STYLE)
	return style&WS_MINIMIZE != 0
}

// Attempt to set the targetWindow to the foreground
func (a *App) attemptSetForeground(targetWindow w32.HWND) bool {
	// We do not use the returning bool because from AHK contributors it is unreliable
	_ = w32.SetForegroundWindow(targetWindow)

	time.Sleep(30 * time.Millisecond)
	// Instead we do a simple check against currently active Foreground Window
	newForeWindow := w32.GetForegroundWindow()
	// fmt.Printf("active Window : %v with HWND : %v\n", w32.GetWindowText(newForeWindow), newForeWindow)
	// fmt.Printf("targetWindow : %v with HWND : %v\n", w32.GetWindowText(targetWindow), targetWindow)
	return newForeWindow == targetWindow
}

// Window bool to Go
func BoolToBOOL(value bool) uintptr {
	if value {
		return 1
	}
	return 0
}

// !! IMPORTANT, detaching threads we added earlier
func detachThreadInputs(currentThread, activeThread, windowThread w32.HANDLE) {
	if currentThread != 0 && currentThread != activeThread {
		procAttachThreadInput.Call(
			uintptr(currentThread),
			uintptr(activeThread),
			BoolToBOOL(false),
		)
	}

	if windowThread != 0 && currentThread != 0 && windowThread != currentThread {
		procAttachThreadInput.Call(
			uintptr(windowThread),
			uintptr(currentThread),
			BoolToBOOL(false),
		)
	}
}

// This is done inside attemptSetForeground for now
// Check if our targetWindow was correctly brought forward
// func (a *App) hasWindowActivated(targetWindow w32.HWND) bool {
// 	return a.getForegroundWindow() == targetWindow
// }

func (a *App) setForegroundWindowEx(targetWindow w32.HWND, origForegroundWindow w32.HWND) w32.HWND {
	// runtime.LogPrint(a.ctx, "Inside setForegroundWindowEx")

	// Check if our window is already foreground, return if it is
	if targetWindow == origForegroundWindow {
		return targetWindow
	}

	// Check if minimized, restore if it is
	minimized := a.isWindowMinimized(targetWindow)
	if minimized {
		// runtime.LogPrint(a.ctx, "Window was minimized.. restoring..")
		w32.ShowWindow(targetWindow, 9) // 9 == SW_RESTORE
	}

	newForeWindow := a.getForegroundWindow()

	// runtime.LogPrint(a.ctx, "First Activation")
	// First attempt at SetForeground
	ATTEMPT_SET_FORE = false
	ATTEMPT_SET_FORE = a.attemptSetForeground(targetWindow)
	if ATTEMPT_SET_FORE {
		// runtime.LogPrint(a.ctx, "First Activation sucess..")
		// currentWindow := a.getForegroundWindow()
		// procBringWindowToTop.Call(uintptr(currentWindow))
		return targetWindow
	}

	// runtime.LogPrint(a.ctx, "First Activation failed, attaching threads...")
	// We failed so next we attach our mainThread to the targetWindow before trying again
	currentThread := w32.HANDLE(windows.GetCurrentThreadId())
	activeThread := a.getWindowThreadProcessId(newForeWindow)
	windowThread := a.getWindowThreadProcessId(targetWindow)

	// Check that our original window still exists
	if int32(origForegroundWindow) != 0 {
		if currentThread != 0 && currentThread != activeThread && !a.isWindowHung(uintptr(origForegroundWindow)) {
			procAttachThreadInput.Call(
				uintptr(currentThread),
				uintptr(activeThread),
				BoolToBOOL(true),
			)
		}

		if windowThread != 0 && currentThread != 0 && windowThread != currentThread {
			procAttachThreadInput.Call(
				uintptr(windowThread),
				uintptr(currentThread),
				BoolToBOOL(true),
			)
		}
	}

	// runtime.LogPrint(a.ctx, "Activation with threads..\n")
	// robotgo.KeyTap("alt")
	SimulateAltPress()
	ATTEMPT_SET_FORE = a.attemptSetForeground(targetWindow)
	// If success we return
	if ATTEMPT_SET_FORE {
		// runtime.LogPrint(a.ctx, "Activation with threads success.. bringing window to top")
		// !!! IMPORTANT !!! ---- Detach threads
		detachThreadInputs(currentThread, activeThread, windowThread)
		return newForeWindow
	} else {
		// If it did not succeed we send double alt and we try again
		// runtime.LogPrint(a.ctx, "Activation with threads failed.. sending double alt..")
		// runtime.LogPrint(a.ctx, "Double tap activation")
		SimulateAltPress()
		SimulateAltPress()

		// Last try
		ATTEMPT_SET_FORE = a.attemptSetForeground(targetWindow)
		// !!! IMPORTANT !!! ---- Detach threads
		detachThreadInputs(currentThread, activeThread, windowThread)
	}

	// If success bring to top
	// Should not be needed
	if ATTEMPT_SET_FORE {
		// runtime.LogPrint(a.ctx, "Activation after double tap alt with threads success.. bringing window to top")
		// runtime.LogPrintf(a.ctx, "procBringWindowToTop")
		currentWindow := a.getForegroundWindow()
		procBringWindowToTop.Call(uintptr(currentWindow))
	} else {
		runtime.LogPrintf(a.ctx, "Failed to activate a window")
	}

	return newForeWindow
}

// var TriedKeyUp bool

// Try up to 5 times
// for i := 0; i < 5; i++ {
// 	if !TriedKeyUp {
// 		TriedKeyUp = true
// 		// Send alt up to allow SetForeground
// 		robotgo.KeyToggle("alt", "up")
// 	}
// 	runtime.LogPrint(a.ctx, "For loop activation")
// 	// We try to SetForeground again up to 5 times
// 	ATTEMPT_SET_FORE = a.attemptSetForeground(targetWindow)
// 	// If success we return
// 	if ATTEMPT_SET_FORE {
// 		return newForeWindow
// 	}
// }
