package main

import (
	"github.com/gonutz/w32/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) KeybindUpdatedEvent() {
	runtime.LogPrintf(a.ctx, "Running KeybindUpdatedEvent %v", keybindMap)
	runtime.EventsEmit(a.ctx, "KeybindsUpdate")
}

func (a *App) CharSelectedEvent(activeChar w32.HWND) {
	runtime.LogPrintf(a.ctx, "CharSelectedEvent %v", activeChar)
	runtime.EventsEmit(a.ctx, "CharSelectedEvent", activeChar)
}

func (a *App) UpdateOrganizerRunning() {
	runtime.LogPrintf(a.ctx, "Running updateMainHookState %t", isOrganizerRunning)
	runtime.EventsEmit(a.ctx, "updateOrganizerRunningState", isOrganizerRunning)
}

func (a *App) UpdateDofusWindows() {
	runtime.EventsEmit(a.ctx, "updatedCharacterOrder", a.DofusWindows)
}
