package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"github.com/gonutz/w32/v2"
)

// gets the .exe dir and returns it as string
func getExecutableDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exePath), nil
}

// Random bullshit to get .exe name of a window by using it's HWND
func GetExecutableName(hwnd w32.HWND) (string, error) {
	_, pid := w32.GetWindowThreadProcessId(hwnd)

	handle, _, _ := syscall.NewLazyDLL("kernel32.dll").
		NewProc("OpenProcess").
		Call(PROCESS_QUERY_LIMITED_INFORMATION, 0, uintptr(pid))
	if handle == 0 {
		return "", fmt.Errorf("unable to open process for PID %d", pid)
	}
	defer syscall.CloseHandle(syscall.Handle(handle))

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

// Extracts char name and class from Dofus Window Title
func parseTitleComponents(title string) (string, string) {
	parts := strings.Split(title, " - ")
	if len(parts) < 2 {
		return "Unknown", "Unknown"
	}
	return parts[0], parts[1]
}

// later used
// func getVKStringFromVKCode(key types.VKCode) string {
// 	return key.String()
// }
