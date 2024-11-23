package main

import (
	"context"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/go-vgo/robotgo"
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
	isHungAppWindow       = user32.NewProc("IsHungAppWindow")
	user32                = syscall.NewLazyDLL("user32.dll")
)

var ATTEMPT_SET_FORE bool

// NEED DEBOUNCING
func (a *App) ActivateNextChar() {
	isDofus, index := a.IsWindowDofus()
	if isDofus {
		runtime.LogPrintf(a.ctx, "activating next char")
		nextIndex := (index + 1) % len(a.DofusWindows)
		nextWindow := win.HWND(a.DofusWindows[nextIndex].Hwnd)
		a.WinActivate(w32.HWND(nextWindow))
	}
}

func (a *App) UpdateTemporaryDofusWindows(tempChars []WindowInfo) {
	if tempChars != nil || len(tempChars) != 0 {
		a.DofusWindows = tempChars
	} else {
		runtime.LogErrorf(a.ctx, "error while updating temporary chars to a.DofusWindows %v", tempChars)
	}
}

// NEED DEBOUNCING
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
		runtime.LogPrintf(a.ctx, "Target window is not valid. %v", targetWindow)
		return origForegroundWindow
	}

	return a.setForegroundWindowEx(targetWindow, origForegroundWindow)
}

// Iterate through all active Windows and populate a.DofusWindows with them
func EnumWindowsCallback(ctx context.Context, hwnd w32.HWND, a *App) bool {
	// Get the window title
	title := w32.GetWindowText(hwnd)
	processName, _ := w32.GetClassName(hwnd)
	exeName, _ := GetExecutableName(hwnd)

	// We check if exe is Dofus, this runs once, should not cause any issues
	if exeName == "Dofus.exe" && processName == "UnityWndClass" && !strings.Contains(title, "Dofus") {
		characterName, class := parseTitleComponents(title)

		// TODO :
		// need to re order it self before append
		// if characters.ini has any

		a.DofusWindows = append(a.DofusWindows, WindowInfo{
			Title:         title,
			Hwnd:          uint64(hwnd),
			CharacterName: characterName,
			Class:         class,
		})

		runtime.LogPrintf(ctx, "Processed window: %s ", title)
	}
	return true
}

// Updates the order of the list of Characters
func (a *App) UpdateDofusWindowsOrder(loggedInCharacters []WindowInfo) ([]WindowInfo, error) {
	// TODO : use this after fetching, maybe try to set exeDir as global var again

	if loggedInCharacters == nil || len(loggedInCharacters) == 0 {
		return a.DofusWindows, nil
	}

	// Load the INI file
	exeDir, _ = getExecutableDir()
	charactersIniFilePath := filepath.Join(exeDir, "characters.ini")
	// no error handling because i dont have time
	iniFile, _, _ := loadINIFile(charactersIniFilePath)

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

	for _, savedChar := range savedOrder {
		if _, exists := processed[savedChar]; exists {
			continue
		}

		if loggedChar, exists := loggedInMap[savedChar]; exists {
			newOrderKnown = append(newOrderKnown, loggedChar)
			processed[savedChar] = true
		} else {
			processed[savedChar] = true
		}

		processed[savedChar] = true
	}

	for _, loggedChar := range loggedInCharacters {
		if _, exists := processed[loggedChar.CharacterName]; !exists {
			newOrderUnknown = append(newOrderUnknown, loggedChar)
			processed[loggedChar.CharacterName] = true
		}
	}

	newOrderKnown = append(newOrderKnown, newOrderUnknown...)

	a.DofusWindows = newOrderKnown

	return newOrderKnown, nil
}

// Used by the frontend to fetch the array, I think it might be useless now?
func (a *App) GetDofusWindows() []WindowInfo {
	_, err, exists := loadINIFile(charactersIniFilePath)
	runtime.LogPrintf(a.ctx, "characters exists %t", exists)
	if err != nil {
		runtime.LogError(a.ctx, "Error with the ini file")
	}

	// check if inifile exists

	runtime.LogPrintf(a.ctx, "Calling refreshAndUpdateCharacterList with exists : %t", exists)
	a.refreshAndUpdateCharacterList(exists)

	if len(a.DofusWindows) > 0 {
		return a.DofusWindows
	}
	return nil
}

// Set main window to always be on top
func (a *App) SetAlwaysOnTop() {
	isAlwaysOnTop = !isAlwaysOnTop
	runtime.WindowSetAlwaysOnTop(a.ctx, isAlwaysOnTop)
}

// Check if user forground window is a Dofus window ->
// return true and its index in list
// || return false and 0 if its not Dofus
func (a *App) IsWindowDofus() (bool, int) {
	activeWindowTitle := robotgo.GetTitle()
	// TODO: use isWindowValid to check before
	windowTitleMap := make(map[string]int)
	for i, window := range a.DofusWindows {
		windowTitleMap[window.Title] = i
	}
	var currentIndex int
	currentIndex, found := windowTitleMap[activeWindowTitle]
	if !found {
		return false, 0
	}
	return true, currentIndex
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

	// Instead we do a simple check against currently active Foreground Window
	newForeWindow := w32.GetForegroundWindow()
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
func detachThreadInputs(isAttachedMyToFore bool, isAttachedForeToTarget bool, mainThreadID, foreThread, targetThread uint32) {
	if isAttachedMyToFore {
		_, _, _ = procAttachThreadInput.Call(
			uintptr(mainThreadID),
			uintptr(foreThread),
			uintptr(0),
		)
	}
	if isAttachedForeToTarget {
		_, _, _ = procAttachThreadInput.Call(
			uintptr(foreThread),
			uintptr(targetThread),
			uintptr(0),
		)
	}
}

// This is done inside attemptSetForeground for now
// Check if our targetWindow was correctly brought forward
// func (a *App) hasWindowActivated(targetWindow w32.HWND) bool {
// 	return a.getForegroundWindow() == targetWindow
// }

func (a *App) setForegroundWindowEx(targetWindow w32.HWND, origForegroundWindow w32.HWND) w32.HWND {
	runtime.LogPrint(a.ctx, "Inside setForegroundWindowEx")

	targetThread := a.getWindowThreadProcessId(targetWindow)

	// Check if our window is already foreground, return if it is
	if targetWindow == origForegroundWindow {
		return targetWindow
	}

	// Check if minimized, restore if it is
	minimized := a.isWindowMinimized(targetWindow)
	if minimized {
		w32.ShowWindow(targetWindow, 9) // 9 == SW_RESTORE
	}

	newForeWindow := a.getForegroundWindow()

	runtime.LogPrint(a.ctx, "First Activation")
	// First attempt at SetForeground
	ATTEMPT_SET_FORE = a.attemptSetForeground(targetWindow)
	if ATTEMPT_SET_FORE {
		return targetWindow
	}

	// We failed so next we attach our mainThread to the targetWindow before trying again
	isAttachedToMyFore := false
	isAttachedForeToTarget := false
	mainThreadID := windows.GetCurrentThreadId()
	foreThread := a.getWindowThreadProcessId(newForeWindow)

	// Check that our original window still exists
	if int32(origForegroundWindow) != 0 {
		if foreThread != 0 && int32(mainThreadID) != int32(foreThread) && !a.isWindowHung(uintptr(origForegroundWindow)) {
			ret, _, _ := procAttachThreadInput.Call(
				uintptr(mainThreadID),
				uintptr(foreThread),
				BoolToBOOL(true),
			)
			isAttachedToMyFore = ret != 0
		}

		if foreThread != 0 && targetThread != 0 && foreThread != targetThread {
			ret, _, _ := procAttachThreadInput.Call(
				uintptr(foreThread),
				uintptr(targetThread),
				BoolToBOOL(true),
			)
			isAttachedForeToTarget = ret != 0
		}
	}

	var TriedKeyUp bool

	// Try up to 5 times
	for i := 0; i < 5; i++ {
		if !TriedKeyUp {
			TriedKeyUp = true
			// Send alt up to allow SetForeground
			robotgo.KeyToggle("alt", "up")
		}
		runtime.LogPrint(a.ctx, "For loop activation")
		// We try to SetForeground again up to 5 times
		ATTEMPT_SET_FORE = a.attemptSetForeground(targetWindow)
		// If success we return
		if ATTEMPT_SET_FORE {
			return newForeWindow
		}
	}

	// If it did not succeed we send double alt and we try again
	if !ATTEMPT_SET_FORE {
		runtime.LogPrint(a.ctx, "Double tap activation")
		robotgo.KeyTap("alt")
		robotgo.KeyTap("alt")

		// Last try
		ATTEMPT_SET_FORE = a.attemptSetForeground(targetWindow)

	}

	// !!! IMPORTANT !!! ---- Detach threads
	detachThreadInputs(isAttachedToMyFore, isAttachedForeToTarget, mainThreadID, uint32(foreThread), uint32(targetThread))

	// If success bring to top
	// Should not be needed
	// if ATTEMPT_SET_FORE {
	// 	currentWindow := a.getForegroundWindow()
	// 	procBringWindowToTop.Call(uintptr(currentWindow))
	// }

	if !ATTEMPT_SET_FORE {
		runtime.LogPrintf(a.ctx, "Failed to activate a window")
	}

	return newForeWindow
}
