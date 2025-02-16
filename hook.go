package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/mouse"
	"github.com/moutend/go-hook/pkg/types"
)

func (a *App) ActivateAction(action string) {
	if len(action) > 0 {
		switch action {
		case "NextChar":
			a.ActivateNextChar()
		case "PreviousChar":
			a.ActivatePreviousChar()
		default:
			a.ActivateCharacter(action)
		}
	}
}

// Used to Pause Organizer
func (a *App) PauseHook() {
	// runtime.LogPrint(a.ctx, "pausing hook")
	isOrganizerRunning = false
}

// Used to Resume Organizer
func (a *App) ResumeHook() {
	// runtime.LogPrint(a.ctx, "resuming hook")
	isOrganizerRunning = true
}

func (a *App) PauseIndividualsHook() {
	IsIndividualKeybindActive = false
}

func (a *App) ResumeIndividualsHook() {
	IsIndividualKeybindActive = true
}

// InstallHook starts the keyboard hook
func (a *App) InstallHook() error {
	if keyboardChan == nil {
		keyboardChan = make(chan types.KeyboardEvent, 100) // Initialize the channel
	}
	if err := keyboard.Install(nil, keyboardChan); err != nil {
		return err
	}

	if mouseChan == nil {
		mouseChan = make(chan types.MouseEvent, 100)
	}

	if err := mouse.Install(nil, mouseChan); err != nil {
		return err
	}

	fmt.Println("Keyboard hook installed.")
	return nil
}

// UninstallHook stops the keyboard hook
func (a *App) UninstallHook() {
	if err := keyboard.Uninstall(); err != nil {
		fmt.Println("Error while uninstalling hook:", err)
	} else {
		fmt.Println("Keyboard hook uninstalled.")
	}

	if err := mouse.Uninstall(); err != nil {
		fmt.Println("Error while uninstalling hook:", err)
	} else {
		fmt.Println("Mouse hook uninstalled.")
	}
}

// Signal handler to stop the hook on interrupt
func (a *App) handleInterrupt() {
	signalChan := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Println("Received interrupt signal. Stopping the keyboard hook...")
		fmt.Printf("Current winEvHook: %v\n", winEvHook)
		a.UninstallHook()
		done <- true
	}()

	<-done
}

func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	a.handleInterrupt()

	fmt.Println("Hook uninstalled.. exiting..")
	return false
}

func (a *App) handleKeyDebounce(eventKey int32) {
	if time.Since(lastInputTime) > (300 * time.Millisecond) {
		lastInputTime = time.Now()
		a.handleKeyDown(eventKey)
	}
}

func (a *App) handleMouseDebounce(eventKey int32) {
	if time.Since(lastInputTime) > (300 * time.Millisecond) {
		lastInputTime = time.Now()
		a.handleMouseDown(eventKey)
	}
}

// lowOrder := uint16(kEvent.MouseData & 0xFFFF)
// TODO: make sure this properly close?
func (a *App) GoHook() error {
	for {
		select {
		// TODO: Maybe introduce timeout?
		case <-time.After(59 * time.Minute):
			fmt.Println("Received timeout signal")
			return nil
		case kEvent := <-keyboardChan:
			switch kEvent.Message {
			case types.WM_KEYDOWN:
				a.handleKeyDebounce(int32(kEvent.VKCode))
			case types.WM_KEYUP:
				a.handleKeyUp(int32(kEvent.VKCode))
			}
		case mEvent := <-mouseChan:
			switch mEvent.Message {
			case WM_XBUTTONDOWN:
				highOrder := int((mEvent.MouseData >> 16) & 0xFFFF)
				a.handleMouseDebounce(int32(highOrder)) // X BUTTON 1
			case WM_XBUTTONUP:
				highOrder := int((mEvent.MouseData >> 16) & 0xFFFF)
				a.handleMouseUp(int32(highOrder)) // X BUTTON 1
			}
		}
	}
}

func (a *App) handleMouseDown(eventKey int32) {
	if keybinds, exists := keybindMap[eventKey]; exists {
		// Action found, perform the corresponding action
		if keybinds.Action == "StopOrganizer" && !isMousePressed[eventKey] {
			isOrganizerRunning = !isOrganizerRunning
			a.UpdateOrganizerRunning()
		} else if isOrganizerRunning && !isKeyPressed[eventKey] {
			a.ActivateAction(keybindMap[eventKey].Action)
		}
		isMousePressed[eventKey] = true
	}
}

func (a *App) handleMouseUp(eventKey int32) {
	if _, exists := keybindMap[eventKey]; exists {
		// Setting key pressed back to false
		if isMousePressed[eventKey] {
			isMousePressed[eventKey] = false
		}
	}
}

func (a *App) handleKeyDown(eventKey int32) {
	if keybinds, exists := keybindMap[eventKey]; exists {
		// Action found, perform the corresponding action
		if keybinds.Action == "StopOrganizer" && !isKeyPressed[eventKey] {
			isOrganizerRunning = !isOrganizerRunning
			a.UpdateOrganizerRunning()
		} else if isOrganizerRunning && !isKeyPressed[eventKey] {
			a.ActivateAction(keybindMap[eventKey].Action)
		}
		isKeyPressed[eventKey] = true
	}
}

func (a *App) handleKeyUp(eventKey int32) {
	if _, exists := keybindMap[eventKey]; exists {
		// Setting key pressed back to false
		if isKeyPressed[eventKey] {
			isKeyPressed[eventKey] = false
		}
	}
}
