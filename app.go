package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/gonutz/w32/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/ini.v1"
)

// Errors are not handled properly at ANY point in this APP
// If something goes wrong, well good luck!

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
	Title         string `json:"title"`
	Hwnd          uint64 `json:"hwnd"`
	CharacterName string
	Class         string
	Order         int
}

// Create the App
func NewApp() *App {
	return &App{}
}

var (
	charactersIniFilePath         string
	configFilePath                string
	lastKeyHoldTime               time.Time
	modKernel32                   = syscall.NewLazyDLL("kernel32.dll")
	procQueryFullProcessImageName = modKernel32.NewProc("QueryFullProcessImageNameW")
	isMainHookActive              bool
	toggleListenerKeybind         string
	configFile                    *ini.File
	exeDir                        string
	stopOrganizerKeybind          int32
	previousCharKeybind           int32
	nextCharKeybind               int32
)

const (
	PROCESS_QUERY_LIMITED_INFORMATION = 0x1000
)

func (a *App) UpdateMainHookState() {
	runtime.EventsEmit(a.ctx, "updateMainHookState", isMainHookActive)
}

func (a *App) GetToggleListenerKeybind() string {
	if len(toggleListenerKeybind) > 0 {
		return toggleListenerKeybind
	}
	return "Invalid Keybind"
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// // Send alt once allows the APP to use SetForegroundWindow with no flashing
	// // Credit Lexikos from :
	// // https://github.com/AutoHotkey/AutoHotkey/blob/581114c1c7bb3890ff61cf5f6e1f1201cd8c8b78/source/window.cpp#L89
	SimulateAltPress()

	// get the current dir
	exeDir, err := getExecutableDir()
	if err != nil {
		runtime.LogPrintf(a.ctx, "Error retrieving executable directory: %v\n", err)
		return
	}
	runtime.LogPrintf(a.ctx, "exeDir %s", exeDir)

	configFilePath = filepath.Join(exeDir, "config.ini")
	// check if config  ini file exists
	configFile, err, exists := loadINIFile(configFilePath)
	if !exists {
		a.CreateConfigSection(configFile, exeDir)
	}
	runtime.LogPrintf(a.ctx, "config exists %t", exists)
	if err != nil {
		runtime.LogError(a.ctx, "Error with the ini file")
	}

	charactersIniFilePath = filepath.Join(exeDir, "characters.ini")
	// check if characters ini file exists
	// iniFile, err, exists := loadINIFile(charactersIniFilePath)
	// runtime.LogPrintf(a.ctx, "characters exists %t", exists)
	// if err != nil {
	// 	runtime.LogError(a.ctx, "Error with the ini file")
	// }

	// Initialize our array
	a.refreshAndUpdateCharacterList(exists)

	// Start main hook for input listener
	a.addMainHook()
}

func (a *App) refreshAndUpdateCharacterList(exists bool) {
	a.DofusWindows = []WindowInfo{}

	// Loop through windows and populate our array
	w32.EnumWindows(func(hwnd w32.HWND) bool {
		return EnumWindowsCallback(a.ctx, hwnd, a)
	})

	runtime.LogPrintf(a.ctx, "Looped through Windows and inside refreshAndUpdateCharacterList with exists : %t", exists)
	// This stinks
	if !exists {
		a.SaveCharacterList(a.DofusWindows)
	}

	runtime.LogPrintf(a.ctx, "end of refresh updating Dofus windows")
}

func getExecutableDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exePath), nil
}

func EnumWindowsCallback(ctx context.Context, hwnd w32.HWND, a *App) bool {
	// Get the window title
	title := w32.GetWindowText(hwnd)
	processName, _ := w32.GetClassName(hwnd)
	exeName, _ := GetExecutableName(hwnd)

	// We check if exe is Dofus, this runs once, should not cause any issues
	if exeName == "Dofus.exe" && processName == "UnityWndClass" && !strings.Contains(title, "Dofus") {
		characterName, class := parseTitleComponents(title)

		// this will need to change to check if user
		// changed order of list in frontend

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

func (a *App) loadCharacterList(cfg *ini.File) ([]string, error) {
	section := cfg.Section("Characters")

	var characterNames []string
	for _, key := range section.Keys() {
		characterNames = append(characterNames, key.Name())
	}

	return characterNames, nil
}

// logic to try and keep order as much as possible even if
// some characters are not currently logged in
func (a *App) UpdateOrder(loggedInNames []string, savedOrder []string) []string {
	runtime.LogPrintf(a.ctx, "INSIDE THIS UPDATE ORDER FUNCTION")
	runtime.LogPrintf(a.ctx, "LOGGED IN NAMES: %v", loggedInNames)
	runtime.LogPrintf(a.ctx, "SAVED ORDER: %v", savedOrder)

	// Step 1: Create a map of savedOrder for fast lookup
	savedMap := make(map[string]bool)
	for _, char := range savedOrder {
		savedMap[char] = true
	}

	// Step 2: Create a new list to hold the updated order
	var updatedOrder []string
	var newChars []string // This will hold new characters

	// Step 3: Add new characters first (those not in savedOrder)
	for _, loggedChar := range loggedInNames {
		if !savedMap[loggedChar] { // If character wasn't in savedOrder
			newChars = append(newChars, loggedChar)
		}
	}

	// Step 4: Merge new characters at a specific position (e.g., after the first X saved characters)
	// Example: Insert the new characters after the first 2 characters from savedOrder
	insertPosition := 2                                                 // This can be customized as per requirement
	updatedOrder = append(updatedOrder, savedOrder[:insertPosition]...) // Add first part of savedOrder

	// Insert new characters after the first X positions
	updatedOrder = append(updatedOrder, newChars...)

	// Add the remaining saved characters
	updatedOrder = append(updatedOrder, savedOrder[insertPosition:]...)

	runtime.LogPrintf(a.ctx, "RETURNING UPDATED ORDER: %v", updatedOrder)
	return updatedOrder
}

func (a *App) UpdateDofusWindowsOrder(loggedInCharacters []WindowInfo) ([]WindowInfo, error) {
	runtime.LogPrintf(a.ctx, "Start of Update Order")
	// Load the INI file
	exeDir, _ = getExecutableDir()
	charactersIniFilePath := filepath.Join(exeDir, "characters.ini")
	// runtime.LogPrintf(a.ctx, "charactersIniFilePath %s", charactersIniFilePath)
	// no error handling because i dont have time
	iniFile, _, _ := loadINIFile(charactersIniFilePath)

	// Load saved character names from the INI file
	savedOrder, err := a.loadCharacterList(iniFile)
	if err != nil {
		runtime.LogError(a.ctx, "Error loading character list")
		return nil, err
	}

	for _, charName := range savedOrder {
		runtime.LogPrintf(a.ctx, "saved char %s", charName)
	}

	// map of currrent char names
	loggedInNames := make([]string, len(loggedInCharacters))
	for i, window := range loggedInCharacters {
		loggedInNames[i] = window.CharacterName
		// runtime.LogPrintf(a.ctx, "char name : %s", window.CharacterName)
	}

	updatedOrder := a.UpdateOrder(loggedInNames, savedOrder)

	runtime.LogPrintf(a.ctx, "Saved order: %v", savedOrder)
	runtime.LogPrintf(a.ctx, "Updated order: %v", updatedOrder)

	// create new array of windows base on new order
	var reorderedWindows []WindowInfo
	for _, charName := range updatedOrder {
		for _, window := range loggedInCharacters {
			if window.CharacterName == charName {
				reorderedWindows = append(reorderedWindows, window)
			}
		}
	}

	// save that new order to the file
	a.SaveCharacterList(reorderedWindows)

	// for _, charName := range reorderedWindows {
	// 	runtime.LogPrintf(a.ctx, "charName.CharacterName %s", charName.CharacterName)
	// }

	// return the array of windows to the frontend
	return reorderedWindows, nil
}

// this deletes our section and re create it
// idk if its a good idea but it works :)
func (a *App) SaveCharacterList(dofusWindows []WindowInfo) error {
	charactersIniFilePath := filepath.Join(exeDir, "characters.ini")
	iniFile, _, _ := loadINIFile(charactersIniFilePath)

	iniFile.DeleteSection("Characters")

	section := iniFile.Section("Characters")
	runtime.LogPrintf(a.ctx, "Saving character list: %v", dofusWindows)
	for _, window := range dofusWindows {
		section.Key(window.CharacterName).SetValue("")
	}

	err := iniFile.SaveTo(charactersIniFilePath)
	if err != nil {
		runtime.LogPrintf(a.ctx, "saving INI file: %v", err)
	}

	runtime.LogPrint(a.ctx, "Dofus windows order updated successfully!\n")

	return nil
}

// used to populate our config.ini Sections
func (a *App) CreateConfigSection(cfg *ini.File, exeDir string) {
	section, err := cfg.GetSection("KeyBindings")
	if err != nil {
		section = cfg.Section("KeyBindings")
	}

	if section.Key("StopOrganizer").String() == "" {
		section.Key("StopOrganizer").SetValue("115,F4")
	}
	if section.Key("PreviousChar").String() == "" {
		section.Key("PreviousChar").SetValue("113,F2")
	}
	if section.Key("NextChar").String() == "" {
		section.Key("NextChar").SetValue("114,F3")
	}

	err = cfg.SaveTo(filepath.Join(exeDir, "config.ini"))
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error saving config file: %v", err)
	} else {
		runtime.LogPrintf(a.ctx, "Config file created/updated successfully")
	}
}

// Load if not exists, CREATE ini
func loadINIFile(filePath string) (*ini.File, error, bool) {
	if _, err := os.Stat(filePath); err == nil {
		// File exists, load it
		cfg, err := ini.Load(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to load INI file: %w", err), false
		}
		return cfg, nil, true
	} else {
		// File doesn't exist, create a new one
		cfg := ini.Empty()
		return cfg, nil, false
	}
}

func parseTitleComponents(title string) (string, string) {
	parts := strings.Split(title, " - ")
	if len(parts) < 2 {
		return "Unknown", "Unknown"
	}
	return parts[0], parts[1]
}

// Random bullshit to get .exe
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

	return filepath.Base(syscall.UTF16ToString(buffer[:size])), nil
}

// Used by the frontend to fetch the array
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

// Used by the frontend to update the array
// func (a *App) UpdateDofusWindows() []WindowInfo {
// 	a.DofusWindows = []WindowInfo{}
// 	w32.EnumWindows(func(hwnd w32.HWND) bool {
// 		return EnumWindowsCallback(a.ctx, hwnd, a)
// 	})
// 	if len(a.DofusWindows) > 0 {
// 		return a.DofusWindows
// 	}
// 	return nil
// }

func (a *App) UpdateDofusWindows() {
	runtime.EventsEmit(a.ctx, "updatedCharacterOrder", a.DofusWindows)
}

// Should not be needed, In theory at leastt
// Used to grant ourselves the right to SetForeground
// Very very very unreliable anyway
// func AllowSetForegroundWindow(pid uint32) error {
// 	user32 := windows.NewLazySystemDLL("user32.dll")
// 	proc := user32.NewProc("AllowSetForegroundWindow")

// 	// Call the function with the process ID
// 	r1, _, err := proc.Call(uintptr(pid))
// 	if r1 == 0 {
// 		return fmt.Errorf("AllowSetForegroundWindow failed: %v", err)
// 	}
// 	return nil
// }
