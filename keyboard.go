package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// gohook fails at registering mouse down, its always mouse up
// gohook does not have keychar for mouse events
// gohook does not have KeyPressed event
// gohook KeyDown does not work with keys that do not output text (f1, page up..)
// --
// gohook is based on an open source C library which has all these things but
// somehow did not make the port from C to GO :)

// mouse button not work because of formatting when saving I think ||
// need heavy debounce
func (a *App) mainHook() {
	runtime.LogPrintf(a.ctx, "isMainHookActive %t", isMainHookActive)
	if isMainHookActive {
		for action, keybind := range keybindsList {
			if keybind.KeyCode == 3 || keybind.KeyCode == 4 || keybind.KeyCode == 5 {
				if action == "StopOrganizer" {
					hook.Register(hook.MouseDown, []string{""}, func(e hook.Event) {
						if e.Button == uint16(keybindsList[action].KeyCode) {
							isOrganizerRunning = !isOrganizerRunning
							a.UpdateOrganizerRunning()
							runtime.LogPrintf(a.ctx, "%s mouse : pressed %v\n ------------\n", action, e)
						}
					})
				} else {
					hook.Register(hook.MouseDown, []string{""}, func(e hook.Event) {
						if e.Button == uint16(keybindsList[action].KeyCode) && isOrganizerRunning {
							runtime.LogPrintf(a.ctx, "%s mouse : pressed %v\n ------------\n ", action, e)
						}
					})
				}
			} else {
				if action == "StopOrganizer" {
					hook.Register(hook.KeyHold, []string{keybindsList[action].KeyName}, func(e hook.Event) {
						runtime.LogPrintf(a.ctx, "%s key : pressed %v\n ------------\n", action, e)
						isOrganizerRunning = !isOrganizerRunning
						a.UpdateOrganizerRunning()
					})
				} else {
					hook.Register(hook.KeyHold, []string{keybindsList[action].KeyName}, func(e hook.Event) {
						if isOrganizerRunning {
							a.ActivateAction(action)
							runtime.LogPrintf(a.ctx, "%s key : pressed %v\n ------------\n", action, e)
						}
					})
				}
			}
		}

		s := hook.Start()
		<-hook.Process(s)
	} else {
		hook.End()
		runtime.LogPrint(a.ctx, "hook ended..")
	}
}

func (a *App) ActivateAction(action string) {
	switch action {
	case "NextChar":
		a.ActivateNextChar()
	case "PreviousChar":
		a.ActivatePreviousChar()
	}
}

// Check if user forground window is a Dofus window ->
// return true and its index in list
// || return false and 0 if its not Dofus
func (a *App) IsWindowDofus() (bool, int) {
	activeWindowTitle := robotgo.GetTitle()
	// Need to separarte this logic so we can work with our array
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

func (a *App) GetIndexOfCharacter() {
}

// Start Main Hook
func (a *App) StartHook() {
	runtime.LogPrint(a.ctx, "Starting hook..")
	isMainHookActive = true
	go a.mainHook()
}

// Stop Main Hook
func (a *App) StopHook() {
	runtime.LogPrint(a.ctx, "Stopping hook..")
	isMainHookActive = false
	a.mainHook()
}

// Used to Pause Organizer
func (a *App) PauseHook() {
	runtime.LogPrint(a.ctx, "pausing hook")
	isMainHookActive = false
}

// Used to Resume Organizer
func (a *App) ResumeHook() {
	runtime.LogPrint(a.ctx, "resuming hook")
	isMainHookActive = true
}

// Generic SaveKeybind
func (a *App) SaveKeybind(keycode int32, keyname string, keybindName string) (string, error) {
	configFile, _, _ = loadINIFile(configFilePath)
	section, _ := configFile.GetSection("KeyBindings")

	value := fmt.Sprintf("%d,%s", keycode, keyname)

	section.Key(keybindName).SetValue(value)

	err := configFile.SaveTo(configFilePath)
	if err != nil {
		runtime.LogPrintf(a.ctx, "Error saving INI file: %v", err)
		return "", err
	}

	keybindsList[keybindName] = Keybinds{
		KeyCode: keycode,
		KeyName: strings.ToLower(keyname),
	}

	runtime.LogPrintf(a.ctx, "Updated Keybinds to : %v", keybindsList)
	a.KeybindUpdatedEvent()

	return "sucess", nil
}

// no error handling
// need rework because space in names like mouse 4 does not properly gets parsed
func (a *App) GetAllKeyBindings() map[string]Keybinds {
	// Reload the config file
	configFile, _, _ := loadINIFile(configFilePath)

	// Get the KeyBindings section
	section, err := configFile.GetSection("KeyBindings")
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error getting section 'KeyBindings': %v", err)
		return nil
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
	// keys := map[string]struct {
	// 	KeyCode int32
	// 	KeyName string
	// }{}

	// Get StopOrganizer key
	stopCode, stopName, err := parseKey("StopOrganizer")
	if err != nil {
		return nil
	}
	keybindsList["StopOrganizer"] = Keybinds{
		KeyCode: stopCode,
		KeyName: strings.ToLower(stopName),
	}

	// Get PreviousChar key
	prevCode, prevName, err := parseKey("PreviousChar")
	if err != nil {
		return nil
	}
	keybindsList["PreviousChar"] = Keybinds{
		KeyCode: prevCode,
		KeyName: strings.ToLower(prevName),
	}

	// Get NextChar key
	nextCode, nextName, err := parseKey("NextChar")
	if err != nil {
		return nil
	}
	keybindsList["NextChar"] = Keybinds{
		KeyCode: nextCode,
		KeyName: strings.ToLower(nextName),
	}

	runtime.LogPrintf(a.ctx, "Update keybinds to %v", keybindsList)
	return keybindsList
}

// Simulate Alt, down+up, used to make dumbass microsoft windows let us use it's SetForeground api
func SimulateAltPress() {
	robotgo.KeyTap("alt")
	time.Sleep(50 * time.Millisecond)
}

// easy hook, works only once
// adde := hook.AddEvent(keybindsList["StopOrganizer"].KeyName)
// 		runtime.LogPrintf(a.ctx, "adde : %v", adde)
// 		if adde {
// 			runtime.LogPrintf(a.ctx, "adde triggered")
// 		}
