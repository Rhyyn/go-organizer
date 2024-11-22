package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Generic SaveKeybind
func (a *App) SaveKeybind(keycode int32, keyname string, keybindName string) (string, error) {
	configFile, _, _ = loadINIFile(configFilePath)
	section, _ := configFile.GetSection("KeyBindings")

	value := fmt.Sprintf("%d,%s", keycode, keyname)

	runtime.LogPrintf(a.ctx, "value saved : %v", value)

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

		parts := strings.SplitN(keyValue, ",", 2)
		if len(parts) != 2 {
			err := fmt.Errorf("invalid key value format for '%s': '%s'", keyName, keyValue)
			runtime.LogErrorf(a.ctx, "Error: %v", err)
			return 0, "", err
		}

		keycode, err := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 32)
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error parsing keycode '%s': %v", parts[0], err)
			return 0, "", err
		}

		keyname := strings.TrimSpace(parts[1])

		return int32(keycode), keyname, nil
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
	runtime.LogPrintf(a.ctx, "prevCode  : %d, prevName : %s", prevCode, prevName)
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
