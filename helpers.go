package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"github.com/gonutz/w32/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/ini.v1"
)

// Gets the saved order of Characters from characters.ini
func (a *App) loadCharacterList(cfg *ini.File) ([]string, error) {
	section := cfg.Section("Characters")

	var characterNames []string
	for _, key := range section.Keys() {
		characterNames = append(characterNames, key.Name())
	}

	return characterNames, nil
}

// this deletes our section and re create it
// idk if its a good idea but it works :)
func (a *App) SaveCharacterList(dofusWindows []WindowInfo) error {
	iniFile, _, _ := loadINIFile(charactersFilePath)

	iniFile.DeleteSection("Characters")

	section := iniFile.Section("Characters")
	// runtime.LogPrintf(a.ctx, "Saving character list: %v\n", dofusWindows)

	for _, window := range dofusWindows {
		if !strings.Contains(window.Title, "Dofus") {
			for _, key := range keybindMap {
				if key.Action == window.CharacterName {
					section.Key(window.CharacterName).SetValue(key.Action)
				} else {
					section.Key(window.CharacterName).SetValue("")
				}
			}
		}
	}

	err := iniFile.SaveTo(charactersFilePath)
	if err != nil {
		runtime.LogPrintf(a.ctx, "saving INI file: %v", err)
	}

	// runtime.LogPrint(a.ctx, "Dofus windows order updated successfully!\n")
	a.DofusWindows = dofusWindows
	return nil
}

// Populate our config.ini Sections and add Default keybinds
func (a *App) CreateConfigSection(cfg *ini.File, exeDir string) {
	section, err := cfg.GetSection("KeyBindings")
	if err != nil {
		section = cfg.Section("KeyBindings")
	}
	section.Key("StopOrganizer").SetValue("115,F4")
	section.Key("PreviousChar").SetValue("113,F2")
	section.Key("NextChar").SetValue("114,F3")

	// Not used for now, might be for V2
	// windowSection, err := cfg.GetSection("Window")
	// if err != nil {
	// 	windowSection = cfg.Section("Window")
	// }
	// windowSection.Key("MonitorHandle")
	// windowSection.Key("FullPosition")
	// windowSection.Key("OverlayPosition")

	// if section.Key("StopOrganizer").String() == "" {
	// }
	// if section.Key("PreviousChar").String() == "" {
	// }
	// if section.Key("NextChar").String() == "" {
	// }

	err = cfg.SaveTo(filepath.Join(exeDir, "config.ini"))
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error saving config file: %v", err)
	}
}

// Load an ini file, if does not exists, we create and return it
func loadINIFile(filePath string) (*ini.File, error, bool) {
	configFileMutex.Lock()
	defer configFileMutex.Unlock()

	if _, err := os.Stat(filePath); err == nil {
		// File exists, load it
		cfg, err := ini.Load(filePath)
		if err != nil {
			return nil, err, false
		}
		return cfg, nil, true
	} else {
		// File doesn't exist, create a new one
		cfg := ini.Empty()
		return cfg, nil, false
	}
}

// gets the .exe dir and returns it as string
func getExecutableDir() {
	exePath, err := os.Executable()

	exeDirTemp := filepath.Dir(exePath)

	if err != nil {
		fmt.Printf("error while getting exe dir %v\n", err)
	} else {
		exeDir = exeDirTemp
		configFilePath = filepath.Join(exeDirTemp, "config.ini")
		charactersFilePath = filepath.Join(exeDirTemp, "characters.ini")
	}
}

// Random bullshit to get .exe name of a window by using it's HWND
func GetExecutableName(hwnd w32.HWND) (string, error) {
	_, pid := w32.GetWindowThreadProcessId(hwnd)

	handle, _, _ := syscall.NewLazyDLL("kernel32.dll").
		NewProc("OpenProcess").
		Call(PROCESS_QUERY_LIMITED_INFORMATION, 0, uintptr(pid))
	if handle == 0 {
		return "", fmt.Errorf("unable to open process for PID %d", pid)
	}
	defer syscall.CloseHandle(syscall.Handle(handle))

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

	return filepath.Base(syscall.UTF16ToString(buffer[:size])), nil
}

// Extracts char name and class from Dofus Window Title
func parseTitleComponents(title string) (string, string) {
	parts := strings.Split(title, " - ")
	if len(parts) < 2 {
		return "Unknown", "Unknown"
	}
	return parts[0], parts[1]
}

// Set main window to always be on top
func (a *App) SetAlwaysOnTop() {
	isAlwaysOnTop = !isAlwaysOnTop
	runtime.WindowSetAlwaysOnTop(a.ctx, isAlwaysOnTop)
}

// later used
// func getVKStringFromVKCode(key types.VKCode) string {
// 	return key.String()
// }
