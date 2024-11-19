package main

import (
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/gonutz/w32/v2"
	"github.com/lxn/win"
	hook "github.com/robotn/gohook"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// we need two hooks main and pause

// HOOK
// Need a button to start/stop
func (a *App) addMainHook() {
	chanHook := hook.Start()
	defer hook.End()

	var isKeyPressed bool

	for ev := range chanHook {

		if ev.Rawcode == uint16(stopOrganizerKeybind) {
			if ev.Kind == hook.KeyDown || ev.Kind == hook.KeyHold && !isKeyPressed {
				isKeyPressed = true
				isMainHookActive = !isMainHookActive
				runtime.LogPrint(a.ctx, "invert bool")
				runtime.LogPrintf(a.ctx, "isMainHookActive : %t", isMainHookActive)
				a.UpdateMainHookState()
			}

			if ev.Kind == hook.KeyUp {
				isKeyPressed = false
				runtime.LogPrint(a.ctx, "up")
				a.UpdateMainHookState()
			}
		}

		if !isMainHookActive {
			continue
		}

		if ev.Rawcode == uint16(nextCharKeybind) {
			if ev.Kind == hook.KeyHold {
				if time.Since(lastKeyHoldTime) > 300*time.Millisecond { // This is not very good, need a better implementation
					// Update the last processed time
					lastKeyHoldTime = time.Now()

					// logs the global event for debug
					// runtime.LogPrintf(a.ctx, "%v", ev)

					// activeWindow := robotgo.GetHWND()
					activeWindowTitle := robotgo.GetTitle()

					var currentIndex int
					var nextWindow win.HWND

					// Need to separarte this logic so we can work with our array
					windowTitleMap := make(map[string]int)
					for i, window := range a.DofusWindows {
						windowTitleMap[window.Title] = i
					}

					currentIndex, found := windowTitleMap[activeWindowTitle]

					if !found {
						runtime.LogPrintf(a.ctx, "Current window not found")
						return
					}

					nextIndex := (currentIndex + 1) % len(a.DofusWindows)
					nextWindow = win.HWND(a.DofusWindows[nextIndex].Hwnd)
					a.winActivate(w32.HWND(nextWindow))
					// for i, window := range a.DofusWindows {
					// 	if window.Title == activeWindowTitle {
					// 		runtime.LogPrintf(a.ctx, "current char : %s", a.DofusWindows[i].CharacterName)
					// 		currentIndex = i
					// 		break
					// 	}
					// }

					// Not using this because might trigger anti cheat ?
					// Leave it here cuz might be helpful one day
					// exeName, _ := GetExecutableName(w32.HWND(activeWindow))
				}
			}
		}
		if ev.Rawcode == uint16(previousCharKeybind) {
			if ev.Kind == hook.KeyHold {
				if time.Since(lastKeyHoldTime) > 300*time.Millisecond { // This is not very good, need a better implementation
					// Update the last processed time
					lastKeyHoldTime = time.Now()

					// logs the global event for debug
					// runtime.LogPrintf(a.ctx, "%v", ev)

					// activeWindow := robotgo.GetHWND()
					activeWindowTitle := robotgo.GetTitle()

					var currentIndex int
					var nextWindow win.HWND

					// Need to separate this logic so we can work with our array
					windowTitleMap := make(map[string]int)
					for i, window := range a.DofusWindows {
						windowTitleMap[window.Title] = i
					}

					currentIndex, found := windowTitleMap[activeWindowTitle]

					if !found {
						runtime.LogPrintf(a.ctx, "Current window not found")
						return
					}

					// Reverse the index logic, decrement and wrap around if less than 0
					nextIndex := (currentIndex - 1 + len(a.DofusWindows)) % len(a.DofusWindows)
					nextWindow = win.HWND(a.DofusWindows[nextIndex].Hwnd)
					a.winActivate(w32.HWND(nextWindow))
				}
			}
		}
	}
}

// func (a *App) addPauseHook() {
// 	pauseHook := hook.Start()
// 	defer hook.End()

// 	var isKeyPressed bool

// 	for ev := range pauseHook {
// 		if ev.Rawcode == uint16(stopOrganizerKeybind) {
// 			if ev.Kind == hook.KeyDown || ev.Kind == hook.KeyHold && !isKeyPressed {
// 				isKeyPressed = true
// 				runtime.LogPrint(a.ctx, "invert bool")
// 				isMainHookActive = !isMainHookActive
// 				runtime.LogPrintf(a.ctx, "isMainHookActive : %t", isMainHookActive)
// 			}

// 			if ev.Kind == hook.KeyUp {
// 				runtime.LogPrint(a.ctx, "up")
// 				isKeyPressed = false
// 			}
// 		}
// 	}
// }

// TogglePauseHook
// This is abandonned for now because it makes the app Freeze for no reason

// func (a *App) StartToggleHook(stopChan chan struct{}) {
// 	toggleHook := hook.Start()
// 	defer hook.End()

// 	go func() {
// 		for ev := range toggleHook {
// 			if ev.Kind == hook.KeyDown || ev.Kind == hook.KeyHold {
// 				if _, found := blacklist[ev.Rawcode]; found {
// 					continue
// 				} else {
// 					if keyName, found := Keycode[ev.Rawcode]; found {
// 						// Process the key event
// 						toggleListenerRawCode = ev.Rawcode
// 						a.UpdateToggleKeybind(keyName)
// 						runtime.LogPrintf(a.ctx, "Captured event: %v\n", ev)

// 						// After capturing, close the stop channel to signal stopping
// 						close(stopChan)
// 						return
// 					} else {
// 						runtime.LogPrintf(a.ctx, "Unknown event with Rawcode: %d\n", ev.Rawcode)
// 						a.UpdateToggleKeybind("Invalid.. Try again")
// 					}
// 				}
// 			} else if ev.Kind == hook.MouseDown || ev.Kind == hook.MouseHold {
// 				runtime.LogPrintf(a.ctx, "Captured event: %v\n", ev)
// 				if _, found := blacklist[ev.Button]; found {
// 					continue
// 				} else {
// 					if keyName, found := Keycode[ev.Button]; found {
// 						toggleListenerRawCode = ev.Rawcode
// 						a.UpdateToggleKeybind(keyName)
// 						runtime.LogPrintf(a.ctx, "Captured event: %v\n", ev)

// 						// After capturing, close the stop channel to signal stopping
// 						close(stopChan)
// 						return
// 					} else {
// 						fmt.Printf("Unknown event with Rawcode: %d\n", ev.Rawcode)
// 						a.UpdateToggleKeybind("Invalid.. Try again")
// 					}
// 				}
// 			}
// 		}
// 	}()

// func (a *App) AddPauseHook() {
// 	chanHook := hook.Start()
// 	defer hook.End()

// 	for ev := range chanHook {
// 		if ev.Kind == hook.KeyHold || ev.Kind == hook.KeyDown || ev.Kind == hook.KeyUp ||
// 			ev.Kind == hook.MouseDown || ev.Kind == hook.MouseUp {
// 			toggleListenerKeybind = string(ev.Keychar)
// 			toggleListenerRawCode = int(ev.Rawcode)
// 			break
// 		}
// 	}
// }

func (a *App) SaveStopOrgaKeyBind(keycode int32, keyname string) {
	// runtime.LogPrintf(a.ctx, "keycode to save : %d", keycode)
	// runtime.LogPrintf(a.ctx, "keyname to save : %s", keyname)
	configFile, _, _ = loadINIFile(configFilePath)
	section, err := configFile.GetSection("KeyBindings")
	if err != nil {
		runtime.LogPrintf(a.ctx, "Error getting section: %v", err)
		return
	}

	value := fmt.Sprintf("%d,%s", keycode, keyname)

	section.Key("StopOrganizer").SetValue(value)

	err = configFile.SaveTo(configFilePath)
	if err != nil {
		runtime.LogPrintf(a.ctx, "Error saving INI file: %v", err)
	} else {
		runtime.LogPrintf(a.ctx, "Key saved successfully")
	}
}

func (a *App) SavePreviousCharKeybind(keycode int32, keyname string) {
	configFile, _, _ = loadINIFile(configFilePath)
	section, err := configFile.GetSection("KeyBindings")
	if err != nil {
		runtime.LogPrintf(a.ctx, "Error getting section: %v", err)
		return
	}
	value := fmt.Sprintf("%d,%s", keycode, keyname)

	section.Key("PreviousChar").SetValue(value)

	err = configFile.SaveTo(configFilePath)
	if err != nil {
		runtime.LogPrintf(a.ctx, "Error saving INI file: %v", err)
	} else {
		runtime.LogPrintf(a.ctx, "Key saved successfully")
	}
}

func (a *App) SaveNextCharKeybind(keycode int32, keyname string) {
	configFile, _, _ = loadINIFile(configFilePath)
	section, err := configFile.GetSection("KeyBindings")
	if err != nil {
		runtime.LogPrintf(a.ctx, "Error getting section: %v", err)
		return
	}

	value := fmt.Sprintf("%d,%s", keycode, keyname)

	section.Key("NextChar").SetValue(value)

	err = configFile.SaveTo(configFilePath)
	if err != nil {
		runtime.LogPrintf(a.ctx, "Error saving INI file: %v", err)
	} else {
		runtime.LogPrintf(a.ctx, "Key saved successfully")
	}
}

func (a *App) GetAllKeyBindings() (map[string]struct {
	KeyCode int32
	KeyName string
}, error,
) {
	// Reload the config file
	configFile, _, _ := loadINIFile(configFilePath)

	runtime.LogPrint(a.ctx, "INSIDE")

	// Get the KeyBindings section
	section, err := configFile.GetSection("KeyBindings")
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error getting section 'KeyBindings': %v", err)
		return nil, err
	}

	// Function to parse the key value
	parseKey := func(keyName string) (int32, string, error) {
		keyValue := section.Key(keyName).String()
		if keyValue == "" {
			err := fmt.Errorf("'%s' key not found", keyName)
			runtime.LogErrorf(a.ctx, "Error: %v", err) // Log error here
			return 0, "", err
		}

		var keycode int32
		var keyname string
		_, err := fmt.Sscanf(keyValue, "%d,%s", &keycode, &keyname)
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error parsing key value '%s': %v", keyValue, err)
			return 0, "", err
		}

		return keycode, keyname, nil
	}

	// Read all the keys
	keys := map[string]struct {
		KeyCode int32
		KeyName string
	}{}

	// Get StopOrganizer key
	stopCode, stopName, err := parseKey("StopOrganizer")
	if err != nil {
		return nil, err
	}
	keys["StopOrganizer"] = struct {
		KeyCode int32
		KeyName string
	}{stopCode, stopName}

	// Get PreviousChar key
	prevCode, prevName, err := parseKey("PreviousChar")
	if err != nil {
		return nil, err
	}
	keys["PreviousChar"] = struct {
		KeyCode int32
		KeyName string
	}{prevCode, prevName}

	// Get NextChar key
	nextCode, nextName, err := parseKey("NextChar")
	if err != nil {
		return nil, err
	}
	keys["NextChar"] = struct {
		KeyCode int32
		KeyName string
	}{nextCode, nextName}

	stopOrganizerKeybind = keys["StopOrganizer"].KeyCode
	previousCharKeybind = keys["PreviousChar"].KeyCode
	nextCharKeybind = keys["NextChar"].KeyCode

	return keys, nil
}

func (a *App) PauseHook() {
	runtime.LogPrint(a.ctx, "pausing hook")
	isMainHookActive = false
}

func (a *App) ResumeHook() {
	runtime.LogPrint(a.ctx, "resuming hook")
	isMainHookActive = true
}

// Simulate Alt, down+up, used to make dumbass microsoft windows let us use it's SetForeground api
func SimulateAltPress() {
	robotgo.KeyTap("alt")
	time.Sleep(50 * time.Millisecond)
}
