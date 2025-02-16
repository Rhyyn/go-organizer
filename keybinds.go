package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Generic SaveKeybind
func (a *App) SaveKeybind(keycode int32, keyname string, keybindName string) (string, error) {
	// Open Config
	configFile, _, _ = loadINIFile(configFilePath)
	section, _ := configFile.GetSection("KeyBindings")

	if _, exists := keybindMap[keycode]; exists {
		return "failed", nil
	}

	// Create a string combination of the two
	value := fmt.Sprintf("%d,%s", keycode, strings.ToUpper(keyname))

	// runtime.LogPrintf(a.ctx, "value saved : %v", value)
	// delete existing
	for existingKeycode, keybind := range keybindMap {
		if keybind.Action == keybindName {
			delete(keybindMap, existingKeycode)
			break
		}
	}

	section.Key(keybindName).SetValue(value)

	err := configFile.SaveTo(configFilePath)
	if err != nil {
		runtime.LogPrintf(a.ctx, "Error saving INI file: %v", err)
		return "", err
	}

	keybindMap[keycode] = Keybinds{
		Action:  keybindName,
		KeyName: strings.ToUpper(keyname),
	}

	runtime.LogPrintf(a.ctx, "Updated Saved Keybinds to : %v", keybindMap)
	a.KeybindUpdatedEvent()

	return "sucess", nil
}

// no error handling
func (a *App) GetAllKeyBindings() map[int32]Keybinds {
	configFile, err, _ := loadINIFile(configFilePath)
	if err != nil && configFile != nil {
		fmt.Printf("Error loading config file: %v\n", err)
		return nil
	}

	// Get the KeyBindings section
	section, err := configFile.GetSection("KeyBindings")
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error getting section 'KeyBindings': %v", err)
		return nil
	}

	for _, key := range section.Keys() {
		fmt.Println(key)
		keyValue := key.Value()

		// 114,F
		parts := strings.SplitN(keyValue, ",", 2)
		// if len(parts) != 2 {
		// 	err := fmt.Errorf("invalid key value format for '%s': '%s'", key, keyValue)
		// 	runtime.LogErrorf(a.ctx, "Error: %v", err)
		// 	return 0, "", err
		// }

		// keycode = 114
		keycode, _ := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 32)
		// if err != nil {
		// 	runtime.LogErrorf(a.ctx, "Error parsing keycode '%s': %v", parts[0], err)
		// 	return 0, "", err
		// }
		// keyname = F
		keyname := strings.TrimSpace(parts[1])

		keybindMap[int32(keycode)] = Keybinds{
			Action:  key.Name(),
			KeyName: keyname,
		}
	}

	fmt.Printf("Update keybinds to %v\n", keybindMap)
	return keybindMap
}

func (a *App) FetchKeybindsFromBack() map[int32]Keybinds {
	return keybindMap
}

var procKeybdEvent = user32.NewProc("keybd_event")

const (
	VK_MENU         = 0x12   // Virtual Key Code for Alt
	KEYEVENTF_KEYUP = 0x0002 // Key release flag
)

// Simulate Alt, down+up, used to make dumbass microsoft windows let us use it's SetForeground api
func SimulateAltPress() {
	procKeybdEvent.Call(
		uintptr(VK_MENU),
		0,
		uintptr(KEYEVENTF_KEYUP),
		0,
	)
	time.Sleep(50 * time.Millisecond)
}

// Simulate Alt, down+up, used to make dumbass microsoft let us use it's SetForeground api,
// We're not using native because for some reason it makes it worse?
// func SimulateAltPress() {
// 	robotgo.KeyTap("alt")
// 	time.Sleep(50 * time.Millisecond)
// }
