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

// HOOK
func (a *App) addMainHook() {
	chanHook := hook.Start()
	defer hook.End()

	var isKeyPressed bool
	var lastToggleTime time.Time

	for ev := range chanHook {

		if ev.Rawcode == uint16(stopOrganizerKeybind) {
			now := time.Now()

			if ev.Kind == hook.KeyDown || ev.Kind == hook.KeyHold && !isKeyPressed {
				if !isKeyPressed && now.Sub(lastToggleTime) > 200*time.Millisecond {
					lastToggleTime = now

					isKeyPressed = true
					isMainHookActive = !isMainHookActive
					runtime.LogPrint(a.ctx, "invert bool")
					runtime.LogPrintf(a.ctx, "isMainHookActive : %t", isMainHookActive)
					a.UpdateMainHookState()
				}
			}

			if ev.Kind == hook.KeyUp {
				if isKeyPressed {
					runtime.LogPrint(a.ctx, "up")
					isKeyPressed = false
					a.UpdateMainHookState()
				}
			}
		}

		if !isMainHookActive {
			continue
		}

		if ev.Kind == hook.KeyDown || ev.Kind == hook.KeyHold || ev.Kind == hook.MouseDown {
			activeWindowTitle := robotgo.GetTitle()
			// Need to separarte this logic so we can work with our array
			windowTitleMap := make(map[string]int)
			for i, window := range a.DofusWindows {
				windowTitleMap[window.Title] = i
			}
			var currentIndex int
			currentIndex, found := windowTitleMap[activeWindowTitle]
			if !found {
				// runtime.LogPrintf(a.ctx, "Current window not dofus")
				continue
			}

			if ev.Rawcode == uint16(nextCharKeybind) {
				if ev.Kind == hook.KeyHold {
					if time.Since(lastKeyHoldTime) > 300*time.Millisecond { // This is not very good, need a better implementation
						// Update the last processed time
						lastKeyHoldTime = time.Now()

						// logs the global event for debug
						// runtime.LogPrintf(a.ctx, "%v", ev)
						var nextWindow win.HWND
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
}

// TODO : Extract logic from those 3

func (a *App) SaveStopOrgaKeyBind(keycode int32, keyname string) (string, error) {
	configFile, _, _ = loadINIFile(configFilePath)
	section, _ := configFile.GetSection("KeyBindings")

	value := fmt.Sprintf("%d,%s", keycode, keyname)

	section.Key("StopOrganizer").SetValue(value)

	err := configFile.SaveTo(configFilePath)
	if err != nil {
		runtime.LogPrintf(a.ctx, "Error saving INI file: %v", err)
		return "", err
	}
	return "sucess", nil
}

func (a *App) SavePreviousCharKeybind(keycode int32, keyname string) (string, error) {
	configFile, _, _ = loadINIFile(configFilePath)
	// ingore err for now cuz I CBA
	section, _ := configFile.GetSection("KeyBindings")
	value := fmt.Sprintf("%d,%s", keycode, keyname)

	section.Key("PreviousChar").SetValue(value)

	err := configFile.SaveTo(configFilePath)
	if err != nil {
		runtime.LogPrintf(a.ctx, "Error saving INI file: %v", err)
		return "", err
	}
	return "success", nil
}

func (a *App) SaveNextCharKeybind(keycode int32, keyname string) (string, error) {
	configFile, _, _ = loadINIFile(configFilePath)
	section, _ := configFile.GetSection("KeyBindings")

	value := fmt.Sprintf("%d,%s", keycode, keyname)

	section.Key("NextChar").SetValue(value)

	err := configFile.SaveTo(configFilePath)
	if err != nil {
		runtime.LogPrintf(a.ctx, "Error saving INI file: %v", err)
		return "", err
	}
	return "sucess", nil
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
