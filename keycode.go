package main

type UMap map[uint32]string

func (a *App) GetKeycodes() UMap {
	return Keycode
}

func (a *App) GetStringFromKey(key int) (string, bool) {
	for k, name := range Keycode {
		if key == int(k) {
			return name, true
		}
	}
	return "", false
}

var Keycode = UMap{
	// 0x01: "LBUTTON",    // Left mouse button
	// 0x02: "RBUTTON",    // Right mouse button
	// 0x04: "MBUTTON",  // Middle mouse button (three-button mouse)
	0x70: "F1",          // F1 key
	0x71: "F2",          // F2 key
	0x72: "F3",          // F3 key
	0x73: "F4",          // F4 key
	0x74: "F5",          // F5 key
	0x09: "TAB",         // TAB key
	0x30: "0",           // 0 key
	0x31: "1",           // 1 key
	0x32: "2",           // 2 key
	0x33: "3",           // 3 key
	0x34: "4",           // 4 key
	0x35: "5",           // 5 key
	0x36: "6",           // 6 key
	0x37: "7",           // 7 key
	0x38: "8",           // 8 key
	0x39: "9",           // 9 key
	0x41: "A",           // A key
	0x42: "B",           // B key
	0x43: "C",           // C key
	0x44: "D",           // D key
	0x45: "E",           // E key
	0x46: "F",           // F key
	0x47: "G",           // G key
	0x48: "H",           // H key
	0x49: "I",           // I key
	0x4A: "J",           // J key
	0x4B: "K",           // K key
	0x4C: "L",           // L key
	0x4D: "M",           // M key
	0x4E: "N",           // N key
	0x4F: "O",           // O key
	0x50: "P",           // P key
	0x51: "Q",           // Q key
	0x52: "R",           // R key
	0x53: "S",           // S key
	0x54: "T",           // T key
	0x55: "U",           // U key
	0x56: "V",           // V key
	0x57: "W",           // W key
	0x58: "X",           // X key
	0x59: "Y",           // Y key
	0x5A: "Z",           // Z key
	0xA0: "LSHIFT",      // Left SHIFT key
	0x14: "CAPSLOCK",    // CAPS LOCK key
	0x21: "PAGEUP",      // PAGE UP key
	0x22: "PAGEDOWN",    // PAGE DOWN key
	0x23: "END",         // END key
	0x24: "HOME",        // HOME key
	0x01: "XBUTTON1",    // X1 mouse button (05)
	0x02: "XBUTTON2",    // X2 mouse button (06)
	0x75: "F6",          // F6 key
	0x76: "F7",          // F7 key
	0x77: "F8",          // F8 key
	0x78: "F9",          // F9 key
	0x79: "F10",         // F10 key
	0x7A: "F11",         // F11 key
	0x7B: "F12",         // F12 key
	0x26: "ARROW_UP",    // UP ARROW key
	0x28: "ARROW_DOWN",  // DOWN ARROW key
	0x25: "ARROW_LEFT",  // LEFT ARROW key
	0x27: "ARROW_RIGHT", // RIGHT ARROW key
	0x60: "NUMPAD0",     // Numeric keypad 0 key
	0x61: "NUMPAD1",     // Numeric keypad 1 key
	0x62: "NUMPAD2",     // Numeric keypad 2 key
	0x63: "NUMPAD3",     // Numeric keypad 3 key
	0x64: "NUMPAD4",     // Numeric keypad 4 key
	0x65: "NUMPAD5",     // Numeric keypad 5 key
	0x66: "NUMPAD6",     // Numeric keypad 6 key
	0x67: "NUMPAD7",     // Numeric keypad 7 key
	0x68: "NUMPAD8",     // Numeric keypad 8 key
	0x69: "NUMPAD9",     // Numeric keypad 9 key
	0x2E: "DELETE",      // DEL key
	0x7C: "F13",         // F13 key
	0x7D: "F14",         // F14 key
	0x7E: "F15",         // F15 key
	0x7F: "F16",         // F16 key
	0x80: "F17",         // F17 key
	0x81: "F18",         // F18 key
	0x82: "F19",         // F19 key
	0x83: "F20",         // F20 key
	0x84: "F21",         // F21 key
	0x85: "F22",         // F22 key
	0x86: "F23",         // F23 key
	0x87: "F4",          // F24 key
	0xBA: ";",           // Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the ';:' key
	0xBB: "+",           // For any country/region, the '+' key
	0xBC: ",",           // For any country/region, the ',' key
	0xBD: "-",           // For any country/region, the '-' key
	0xBE: ".",           // For any country/region, the '.' key
	0xBF: "?",           // Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the '/?' key
	0xC0: "~",           // Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the '`~' key
	0xDB: "[",           // Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the '[{' key
	0xDC: "|",           // Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the '\|' key
	0xDD: "]",           // Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the ']}' key
	0xDE: "'",           // Used for miscellaneous characters; it can vary by keyboard. For the US standard keyboard, the 'single-quote/double-quote'
}

// Keycode robotgo hook key's code map
// var Keycode = UMap{
// 	3:   "Middle Mouse",
// 	4:   "Mouse button 4",
// 	5:   "Mouse button 5",
// 	48:  "0",
// 	49:  "1",
// 	50:  "2",
// 	51:  "3",
// 	52:  "4",
// 	53:  "5",
// 	54:  "6",
// 	55:  "7",
// 	56:  "8",
// 	57:  "9",
// 	19:  "Pausebreak",
// 	33:  "pageup",
// 	34:  "pagedown",
// 	35:  "end",
// 	36:  "home",
// 	37:  "Left arrow",
// 	38:  "Up arrow",
// 	39:  "Right arrow",
// 	40:  "Down arrow",
// 	45:  "insert",
// 	46:  "delete",
// 	65:  "A",
// 	66:  "B",
// 	67:  "C",
// 	68:  "D",
// 	69:  "E",
// 	70:  "F",
// 	71:  "G",
// 	72:  "H",
// 	73:  "I",
// 	74:  "J",
// 	75:  "K",
// 	76:  "L",
// 	77:  "M",
// 	78:  "N",
// 	79:  "O",
// 	80:  "P",
// 	81:  "Q",
// 	82:  "R",
// 	83:  "S",
// 	84:  "T",
// 	85:  "U",
// 	86:  "V",
// 	87:  "W",
// 	88:  "X",
// 	89:  "Y",
// 	90:  "Z",
// 	96:  "numpad0",
// 	97:  "numpad1",
// 	98:  "numpad2",
// 	99:  "numpad3",
// 	100: "numpad4",
// 	101: "numpad5",
// 	102: "numpad6",
// 	103: "numpad7",
// 	104: "numpad8",
// 	105: "numpad9",
// 	112: "F1",
// 	113: "F2",
// 	114: "F3",
// 	115: "F4",
// 	116: "F5",
// 	117: "F6",
// 	118: "F7",
// 	119: "F8",
// 	120: "F9",
// 	121: "F10",
// 	122: "F11",
// 	123: "F12",
// 	186: "semicolon",
// 	187: "equalsign",
// 	188: "comma",
// 	189: "dash",
// 	190: "period",
// 	191: "forwardslash",
// 	192: "graveaccent",
// 	219: "openbracket",
// 	220: "backslash",
// 	221: "closebracket",
// 	222: "singlequote",
// }

// Special is the special key map
// var Special = map[string]string{
// 	"~": "`",
// 	"!": "1",
// 	"@": "2",
// 	"#": "3",
// 	"$": "4",
// 	"%": "5",
// 	"^": "6",
// 	"&": "7",
// 	"*": "8",
// 	"(": "9",
// 	")": "0",
// 	"_": "-",
// 	"+": "=",
// 	"{": "[",
// 	"}": "]",
// 	"|": "\\",
// 	":": ";",
// 	`"`: "'",
// 	"<": ",",
// 	">": ".",
// 	"?": "/",
// }
