package main

import "github.com/wailsapp/wails/v2/pkg/runtime"

func (a *App) KeybindUpdatedEvent() {
	runtime.LogPrintf(a.ctx, "Running KeybindUpdatedEvent %v", keybindMap)
	runtime.EventsEmit(a.ctx, "KeybindsUpdate")
}

func (a *App) UpdateOrganizerRunning() {
	runtime.LogPrintf(a.ctx, "Running updateMainHookState %t", isOrganizerRunning)
	runtime.EventsEmit(a.ctx, "updateOrganizerRunningState", isOrganizerRunning)
}

func (a *App) UpdateDofusWindows() {
	runtime.EventsEmit(a.ctx, "updatedCharacterOrder", a.DofusWindows)
}
