package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/go-vgo/robotgo"
	"github.com/gonutz/w32/v2"
	"github.com/lxn/win"
	hook "github.com/robotn/gohook"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/ini.v1"
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
	iniFilePath                   string
	lastKeyHoldTime               time.Time
	modKernel32                   = syscall.NewLazyDLL("kernel32.dll")
	procQueryFullProcessImageName = modKernel32.NewProc("QueryFullProcessImageNameW")
	isMainHookActive              bool
)

const (
	PROCESS_QUERY_LIMITED_INFORMATION = 0x1000
)

// HOOK
// Need a button to start/stop
func (a *App) addMainHook() {
	chanHook := hook.Start()
	defer hook.End()

	for ev := range chanHook {
		if !isMainHookActive {
			continue
		}
		if ev.Rawcode == 114 {
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

					// for i, window := range a.DofusWindows {
					// 	if window.Title == activeWindowTitle {
					// 		runtime.LogPrintf(a.ctx, "current char : %s", a.DofusWindows[i].CharacterName)
					// 		currentIndex = i
					// 		break
					// 	}
					// }

					nextIndex := (currentIndex + 1) % len(a.DofusWindows)
					nextWindow = win.HWND(a.DofusWindows[nextIndex].Hwnd)
					a.winActivate(w32.HWND(nextWindow))

					// Not using this because might trigger anti cheat ?
					// Leave it here cuz might be helpful one day
					// exeName, _ := GetExecutableName(w32.HWND(activeWindow))
				}
			}
		}
	}
}

func (a *App) PauseHook() {
	runtime.LogPrint(a.ctx, "pausing hook")
	isMainHookActive = false
}

func (a *App) ResumeHook() {
	runtime.LogPrint(a.ctx, "resuming hook")
	isMainHookActive = true
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
		fmt.Printf("Error retrieving executable directory: %v\n", err)
		return
	}

	// ini file
	iniFilePath = filepath.Join(exeDir, "characters.ini")

	// check if ini file exists
	iniFile, err, exists := loadINIFile()
	runtime.LogPrintf(a.ctx, "exists %t", exists)
	if err != nil {
		runtime.LogError(a.ctx, "Error with the ini file")
	}

	// Initialize our array
	a.DofusWindows = []WindowInfo{}

	// Loop through windows and populate our array
	w32.EnumWindows(func(hwnd w32.HWND) bool {
		return EnumWindowsCallback(ctx, hwnd, a)
	})

	// This stinks
	if exists {
		a.checkAndAddNewCharacters(iniFile)
	} else {
		a.saveCharacterList(iniFile, a.DofusWindows)
	}

	// Start main hook for keyboard listener
	a.addMainHook()
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
	if exeName == "Dofus.exe" && processName == "UnityWndClass" {
		characterName, class := parseTitleComponents(title)

		// this will need to change to check if user
		// changed order of list in frontend
		order := len(a.DofusWindows) + 1

		a.DofusWindows = append(a.DofusWindows, WindowInfo{
			Title:         title,
			Hwnd:          uint64(hwnd),
			CharacterName: characterName,
			Class:         class,
			Order:         order,
		})

		runtime.LogPrintf(ctx, "Processed window: %s (Order: %d)", title, order)
	}
	return true
}

func (a *App) UpdateDofusWindowsOrder(newOrder []WindowInfo) error {
	// Create a map to quickly look up the index of each character name
	characterNameToIndex := make(map[string]int)
	for i, window := range a.DofusWindows {
		characterNameToIndex[window.CharacterName] = i
	}

	var reorderedWindows []WindowInfo

	// Iterate through the new order and reorder the DofusWindows
	for _, window := range newOrder {
		// find the corresponding window by character name
		index, exists := characterNameToIndex[window.CharacterName]
		if !exists {
			runtime.LogPrint(a.ctx, "here")
			runtime.LogPrintf(a.ctx, "Character name %s not found in the current DofusWindows list\n", window.CharacterName)
			continue
		}

		reorderedWindows = append(reorderedWindows, a.DofusWindows[index])
	}

	// Update
	a.DofusWindows = reorderedWindows

	// Update the Order field for each window based on its new position
	for i := range a.DofusWindows {
		a.DofusWindows[i].Order = i + 1
	}

	// DEBUG
	// for _, window := range a.DofusWindows {
	// 	runtime.LogDebugf(a.ctx, "windowName :%s", window.CharacterName)
	// }

	iniFile, err, _ := loadINIFile()
	if err != nil {
		runtime.LogError(a.ctx, "Erreur dans le loading de l'INI ligne 273")
	}
	a.saveCharacterList(iniFile, a.DofusWindows)
	return nil
}

// Only used once during inital startup or If user save his current order
// This is destructive, might not be the best idea
func (a *App) saveCharacterList(cfg *ini.File, dofusWindows []WindowInfo) error {
	cfg.DeleteSection("Characters")

	section := cfg.Section("Characters")

	for _, window := range dofusWindows {
		section.Key(window.CharacterName).SetValue(fmt.Sprintf("%d", window.Order))
	}

	err := cfg.SaveTo(iniFilePath)
	if err != nil {
		runtime.LogPrintf(a.ctx, "saving INI file: %v", err)
	}

	runtime.LogPrint(a.ctx, "Dofus windows order updated successfully!\n")

	return nil
}

// Check if need to add new chars
// Unknown if it works for now
func (a *App) checkAndAddNewCharacters(cfg *ini.File) error {
	section := cfg.Section("Characters")

	// Create a set of currently active characters
	for _, window := range a.DofusWindows {
		runtime.LogPrintf(a.ctx, "name in order: %s", window.CharacterName)
	}

	activeNamesSet := make(map[string]bool)
	for _, window := range a.DofusWindows {
		activeNamesSet[window.CharacterName] = true
	}

	runtime.LogPrintf(a.ctx, "activeNamesSet : %v", activeNamesSet)

	// Load the current list of characters already in the INI file
	existingNames, err := a.loadCharacterList(cfg)
	if err != nil {
		return fmt.Errorf("error loading character list: %v", err)
	}

	// Create a set of existing character names from the INI file for fast lookup
	existingNamesSet := make(map[string]bool)
	for _, name := range existingNames {
		existingNamesSet[name] = true
	}

	// Loop through the currently logged-in windows and check for new characters
	for _, window := range a.DofusWindows {
		// If the character is not already in the INI file, add it
		if !existingNamesSet[window.CharacterName] && activeNamesSet[window.CharacterName] {
			// Find the correct position based on the character's order
			order := window.Order

			// Iterate over the section to check if the order already exists
			for _, key := range section.Keys() { // Keys() returns a slice of `ini.Key`
				keyName := key.Name()

				// If a character with the same order is found, shift its order
				if existingOrder, err := strconv.Atoi(key.String()); err == nil && existingOrder >= order {
					// Increment the order of the existing character to avoid conflict
					newOrder := existingOrder + 1
					section.Key(keyName).SetValue(fmt.Sprintf("%d", newOrder))
				}
			}

			// Now add the new character with its order
			section.Key(window.CharacterName).SetValue(fmt.Sprintf("%d", order))
		}
	}

	// Save the INI file after adding any new characters
	err = cfg.SaveTo(iniFilePath)
	if err != nil {
		return fmt.Errorf("error saving INI file: %v", err)
	}

	return nil
}

// load characters from ini
func (a *App) loadCharacterList(cfg *ini.File) ([]string, error) {
	section := cfg.Section("Characters")

	var characterNames []string
	for _, key := range section.Keys() {
		characterNames = append(characterNames, key.Name())
	}

	return characterNames, nil
}

// Load if not exists, CREATE ini
func loadINIFile() (*ini.File, error, bool) {
	if _, err := os.Stat(iniFilePath); err == nil {
		// File exists, load it
		cfg, err := ini.Load(iniFilePath)
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

// Simulate Alt, down+up, used to make dumbass microsoft windows let us use it's SetForeground api
func SimulateAltPress() {
	robotgo.KeyTap("alt")
	time.Sleep(50 * time.Millisecond)
}

// Used by the frontend to fetch the array
func (a *App) GetDofusWindows() []WindowInfo {
	if len(a.DofusWindows) > 0 {
		return a.DofusWindows
	}
	return nil
}

// Used by the frontend to update the array
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
