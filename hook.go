package main

import (
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

// slight delay after startup, cause unkown ||
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

// easy hook, works only once
// adde := hook.AddEvent(keybindsList["StopOrganizer"].KeyName)
// 		runtime.LogPrintf(a.ctx, "adde : %v", adde)
// 		if adde {
// 			runtime.LogPrintf(a.ctx, "adde triggered")
// 		}
