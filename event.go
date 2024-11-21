package main

import "github.com/wailsapp/wails/v2/pkg/runtime"

func (a *App) KeybindUpdatedEvent() {
	runtime.LogPrintf(a.ctx, "Running KeybindUpdatedEvent %v", keybindsList)
	runtime.EventsEmit(a.ctx, "KeybindsUpdate", keybindsList)
}
