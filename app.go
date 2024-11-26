package main

import (
	"context"
	"log"
	"sync"
	"syscall"
	"time"

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
	charactersFilePath            string
	configFilePath                string
	modKernel32                   = syscall.NewLazyDLL("kernel32.dll")
	procQueryFullProcessImageName = modKernel32.NewProc("QueryFullProcessImageNameW")
	configFile                    *ini.File
	exeDir                        string
	isAlwaysOnTop                 bool
	keybindMap                    map[int32]Keybinds
	isOrganizerRunning            bool
	isKeyPressed                  map[int32]bool
	isMousePressed                map[int32]bool
	characterFileMutex            sync.Mutex
	configFileMutex               sync.Mutex
	keyboardChan                  chan types.KeyboardEvent
	mouseChan                     chan types.MouseEvent
	WM_XBUTTONDOWN                types.Message = 0x020B // 523 -> XButton Down
	WM_XBUTTONUP                  types.Message = 0x020C // 524 -> XButton UP
	lastInputTime                 time.Time
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

	// Gets exe dir of program and store as var
	getExecutableDir()

	runtime.LogPrintf(a.ctx, "configFilePath : %s\n", configFilePath)
	runtime.LogPrintf(a.ctx, "charactersFilePath : %s\n", charactersFilePath)

	// check if config  ini file exists
	configFile, err, exists := loadINIFile(configFilePath)
	if !exists {
		a.CreateConfigSection(configFile, exeDir)
	}
	if err != nil {
		runtime.LogError(a.ctx, "Error with the ini file")
	}

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

	// Send alt once allows the APP to use SetForegroundWindow with no flashing
	// Credit Lexikos from :
	// https://github.com/AutoHotkey/AutoHotkey/blob/581114c1c7bb3890ff61cf5f6e1f1201cd8c8b78/source/window.cpp#L89
	SimulateAltPress()

	a.handleWindowPosition()

	// Start hooks
	go a.foregroundWindowsHook()

	if err := a.InstallHook(); err != nil {
		log.Fatal(err)
	}

	if err := a.GoHook(); err != nil {
		log.Fatal(err)
	}
}
