package main

import (
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/gonutz/w32/v2"
	"github.com/lxn/win"
	hook "github.com/robotn/gohook"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// we need two hooks main and pause

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
					a.UpdateToggleKeybind("test")

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

					nextIndex := (currentIndex + 1) % len(a.DofusWindows)
					nextWindow = win.HWND(a.DofusWindows[nextIndex].Hwnd)
					a.winActivate(w32.HWND(nextWindow))
					// for i, window := range a.DofusWindows {
					// 	if window.Title == activeWindowTitle {
					// 		runtime.LogPrintf(a.ctx, "current char : %s", a.DofusWindows[i].CharacterName)
					// 		currentIndex = i
					// 		break
					// 	}
					// }

					// Not using this because might trigger anti cheat ?
					// Leave it here cuz might be helpful one day
					// exeName, _ := GetExecutableName(w32.HWND(activeWindow))
				}
			}
		}
		if ev.Rawcode == 113 {
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

					// Need to separate this logic so we can work with our array
					windowTitleMap := make(map[string]int)
					for i, window := range a.DofusWindows {
						windowTitleMap[window.Title] = i
					}

					currentIndex, found := windowTitleMap[activeWindowTitle]

					if !found {
						runtime.LogPrintf(a.ctx, "Current window not found")
						return
					}

					// Reverse the index logic, decrement and wrap around if less than 0
					nextIndex := (currentIndex - 1 + len(a.DofusWindows)) % len(a.DofusWindows)
					nextWindow = win.HWND(a.DofusWindows[nextIndex].Hwnd)
					a.winActivate(w32.HWND(nextWindow))
				}
			}
		}
	}
}

// TogglePauseHook
func (a *App) StartToggleHook() {
	go func() {
		toggleHook := hook.Start()
		defer hook.End()

		for ev := range toggleHook {
			if ev.Kind == hook.KeyDown || ev.Kind == hook.KeyHold {
				if _, found := blacklist[ev.Rawcode]; found {
					continue
				} else {
					if keyName, found := Keycode[ev.Rawcode]; found {
						// valid
						toggleListenerRawCode = ev.Rawcode
						a.UpdateToggleKeybind(keyName)
						runtime.LogPrintf(a.ctx, "event: %v", ev)
					} else {
						// invalid
						runtime.LogPrintf(a.ctx, "Unknown event with Rawcode: %d", ev.Rawcode)
						a.UpdateToggleKeybind("Invalid.. Try again")
					}

					break // Exit loop after processing event
				}
			} else if ev.Kind == hook.MouseDown || ev.Kind == hook.MouseHold {
				runtime.LogPrintf(a.ctx, "event: %v", ev)
				if _, found := blacklist[ev.Button]; found {
					continue
				} else {
					if keyName, found := Keycode[ev.Button]; found {
						// valid
						toggleListenerRawCode = ev.Rawcode
						a.UpdateToggleKeybind(keyName)
						runtime.LogPrintf(a.ctx, "event: %v", ev)
					} else {
						// invalid
						runtime.LogPrintf(a.ctx, "Unknown event with Rawcode: %d", ev.Rawcode)
						a.UpdateToggleKeybind("Invalid.. Try again")
					}

					break // Exit loop after processing event
				}
			}
		}
	}()
	// toggleHook := hook.Start()
	// defer hook.End()

	// go func() {
	// 	for ev := range toggleHook {
	// 		if ev.Kind == hook.KeyDown || ev.Kind == hook.KeyHold {
	// 			if _, found := blacklist[ev.Rawcode]; found {
	// 				continue
	// 			} else {
	// 				if keyName, found := Keycode[ev.Rawcode]; found {
	// 					// valid
	// 					toggleListenerRawCode = ev.Rawcode
	// 					a.UpdateToggleKeybind(keyName)
	// 					runtime.LogPrintf(a.ctx, "event: %v", ev)
	// 				} else {
	// 					// invalid
	// 					runtime.LogPrintf(a.ctx, "Unknown event with Rawcode: %d", ev.Rawcode)
	// 					a.UpdateToggleKeybind("Invalid.. Try again")
	// 				}

	// 				break // Exit loop after processing event
	// 			}
	// 		} else if ev.Kind == hook.MouseDown || ev.Kind == hook.MouseHold {
	// 			runtime.LogPrintf(a.ctx, "event: %v", ev)
	// 			if _, found := blacklist[ev.Button]; found {
	// 				continue
	// 			} else {
	// 				if keyName, found := Keycode[ev.Button]; found {
	// 					// valid
	// 					toggleListenerRawCode = ev.Rawcode
	// 					a.UpdateToggleKeybind(keyName)
	// 					runtime.LogPrintf(a.ctx, "event: %v", ev)
	// 				} else {
	// 					// invalid
	// 					runtime.LogPrintf(a.ctx, "Unknown event with Rawcode: %d", ev.Rawcode)
	// 					a.UpdateToggleKeybind("Invalid.. Try again")
	// 				}

	// 				break // Exit loop after processing event
	// 			}
	// 		}
	// 	}
	// }()
	// this could get improved a lot but I cba
	// go func() {
	// 	for ev := range toggleHook {
	// 		if ev.Kind == hook.KeyDown || ev.Kind == hook.KeyHold {
	// 			if _, found := blacklist[ev.Rawcode]; found {
	// 				continue
	// 			} else {
	// 				if keyName, found := Keycode[ev.Rawcode]; found {
	// 					// valid
	// 					toggleListenerRawCode = ev.Rawcode
	// 					a.UpdateToggleKeybind(keyName)
	// 					runtime.LogPrintf(a.ctx, "event: %v", ev)
	// 				} else {
	// 					// invalid
	// 					runtime.LogPrintf(a.ctx, "Unknown event with Rawcode: %d", ev.Rawcode)
	// 					a.UpdateToggleKeybind("Invalid.. Try again")
	// 				}
	// 				if wasMainHookActive {
	// 					runtime.LogPrint(a.ctx, "Resuming main hook")
	// 					a.UpdateMainHookActiveState(true)
	// 				}
	// 				break
	// 			}
	// 		} else if ev.Kind == hook.MouseDown || ev.Kind == hook.MouseHold {
	// 			runtime.LogPrintf(a.ctx, "fuck %v", ev)
	// 			if _, found := blacklist[ev.Button]; found {
	// 				continue
	// 			} else {
	// 				if keyName, found := Keycode[ev.Button]; found {
	// 					// valid
	// 					toggleListenerRawCode = ev.Rawcode
	// 					a.UpdateToggleKeybind(keyName)
	// 					runtime.LogPrintf(a.ctx, "event: %v", ev)
	// 				} else {
	// 					// invalid
	// 					runtime.LogPrintf(a.ctx, "Unknown event with Rawcode: %d", ev.Rawcode)
	// 					a.UpdateToggleKeybind("Invalid.. Try again")
	// 				}
	// 				if wasMainHookActive {
	// 					a.UpdateMainHookActiveState(true)
	// 				}
	// 				break
	// 			}
	// 		}
	// 	}
	// }()
}

// func (a *App) AddPauseHook() {
// 	chanHook := hook.Start()
// 	defer hook.End()

// 	for ev := range chanHook {
// 		if ev.Kind == hook.KeyHold || ev.Kind == hook.KeyDown || ev.Kind == hook.KeyUp ||
// 			ev.Kind == hook.MouseDown || ev.Kind == hook.MouseUp {
// 			toggleListenerKeybind = string(ev.Keychar)
// 			toggleListenerRawCode = int(ev.Rawcode)
// 			break
// 		}
// 	}
// }

func (a *App) PauseHook() {
	runtime.LogPrint(a.ctx, "pausing hook")
	isMainHookActive = false
}

func (a *App) ResumeHook() {
	runtime.LogPrint(a.ctx, "resuming hook")
	isMainHookActive = true
}

// Simulate Alt, down+up, used to make dumbass microsoft windows let us use it's SetForeground api
func SimulateAltPress() {
	robotgo.KeyTap("alt")
	time.Sleep(50 * time.Millisecond)
}
