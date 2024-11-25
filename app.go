package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/gonutz/w32/v2"
	"github.com/lxn/win"
	"github.com/moutend/go-hook/pkg/types"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/ini.v1"
)

// TODO :Errors are not handled properly at ANY point in this APP
// If something goes wrong, well good luck!

// TODO: Separate struct / var file?

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

type Keybinds struct {
	Action  string
	KeyName string
}

var (
	charactersIniFilePath         string
	configFilePath                string
	modKernel32                   = syscall.NewLazyDLL("kernel32.dll")
	procQueryFullProcessImageName = modKernel32.NewProc("QueryFullProcessImageNameW")
	procGetWindowTextW            = user32.NewProc("GetWindowTextW")
	configFile                    *ini.File
	exeDir                        string
	isAlwaysOnTop                 bool
	keybindMap                    map[int32]Keybinds
	isOrganizerRunning            bool
	isKeyPressed                  map[int32]bool
	isMousePressed                map[int32]bool
	mapMutex                      sync.Mutex
	keyboardChan                  chan types.KeyboardEvent
	mouseChan                     chan types.MouseEvent
	WM_XBUTTONDOWN                types.Message = 0x020B // 523 -> XButton Down
	WM_XBUTTONUP                  types.Message = 0x020C // 524 -> XButton UP
	lastInputTime                 time.Time
	msg                           win.MSG
)

const (
	PROCESS_QUERY_LIMITED_INFORMATION = 0x1000
)

// Create the App
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// get the current dir
	exeDir, err := getExecutableDir()
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error retrieving executable directory: %v\n", err)
		return
	}

	configFilePath = filepath.Join(exeDir, "config.ini")
	// check if config  ini file exists
	configFile, err, exists := loadINIFile(configFilePath)
	if !exists {
		a.CreateConfigSection(configFile, exeDir)
	}
	if err != nil {
		runtime.LogError(a.ctx, "Error with the ini file")
	}

	charactersIniFilePath = filepath.Join(exeDir, "characters.ini")

	// Initialize our array
	a.refreshAndUpdateCharacterList(exists)

	// Initialize map of keybinds and load our saved keybinds
	keybindMap = make(map[int32]Keybinds)
	isKeyPressed = make(map[int32]bool)
	isMousePressed = make(map[int32]bool)
	lastInputTime = time.Now()
	a.GetAllKeyBindings()

	// Start of our Observers
	runtime.EventsOn(a.ctx, "KeybindsUpdate", func(optionalData ...interface{}) {
		runtime.LogPrint(a.ctx, "Keybinds Updated...  Restarting Hooks with updated keybinds...")

		// Stop / Start hook to update their keybinds
		a.handleInterrupt()

		err := a.InstallHook()
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error installing Hook.. %v", err)
		}
	})

	// // Send alt once allows the APP to use SetForegroundWindow with no flashing
	// // Credit Lexikos from :
	// // https://github.com/AutoHotkey/AutoHotkey/blob/581114c1c7bb3890ff61cf5f6e1f1201cd8c8b78/source/window.cpp#L89
	SimulateAltPress()

	go a.testHook()
	// Start main hook for input listener
	if err := a.InstallHook(); err != nil {
		log.Fatal(err)
	}

	if err := a.GoHook(); err != nil {
		log.Fatal(err)
	}
}

// Fetch Windows to see if any new Dofus windows appeared
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

	// Should order before updating front
	a.UpdateDofusWindowsOrder(a.DofusWindows)

	runtime.LogPrintf(a.ctx, "end of refresh updating Dofus windows")
}

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
	exeDir, _ := getExecutableDir()

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
	a.DofusWindows = dofusWindows
	return nil
}

// Populate our config.ini Sections and add Default keybinds
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

// Load an ini file, if does not exists, we create and return it
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
