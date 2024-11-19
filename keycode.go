// Copyright 2016 The go-vgo Project Developers. See the COPYRIGHT
// file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0> or the MIT license
// <LICENSE-MIT or http://opensource.org/licenses/MIT>, at your
// option. This file may not be copied, modified, or distributed
// except according to those terms.

package main

// UMap type map[string]uint16
type UMap map[uint16]string

func (a *App) GetKeycodes() UMap {
	return Keycode
}

// Keycode robotgo hook key's code map
var Keycode = UMap{
	3:   "Middle Mouse",
	4:   "Mouse button 4",
	5:   "Mouse button 5",
	48:  "0",
	49:  "1",
	50:  "2",
	51:  "3",
	52:  "4",
	53:  "5",
	54:  "6",
	55:  "7",
	56:  "8",
	57:  "9",
	19:  "Pausebreak",
	33:  "pageup",
	34:  "pagedown",
	35:  "end",
	36:  "home",
	37:  "Left arrow",
	38:  "Up arrow",
	39:  "Right arrow",
	40:  "Down arrow",
	45:  "insert",
	46:  "delete",
	65:  "A",
	66:  "B",
	67:  "C",
	68:  "D",
	69:  "E",
	70:  "F",
	71:  "G",
	72:  "H",
	73:  "I",
	74:  "J",
	75:  "K",
	76:  "L",
	77:  "M",
	78:  "N",
	79:  "O",
	80:  "P",
	81:  "Q",
	82:  "R",
	83:  "S",
	84:  "T",
	85:  "U",
	86:  "V",
	87:  "W",
	88:  "X",
	89:  "Y",
	90:  "Z",
	96:  "numpad0",
	97:  "numpad1",
	98:  "numpad2",
	99:  "numpad3",
	100: "numpad4",
	101: "numpad5",
	102: "numpad6",
	103: "numpad7",
	104: "numpad8",
	105: "numpad9",
	112: "F1",
	113: "F2",
	114: "F3",
	115: "F4",
	116: "F5",
	117: "F6",
	118: "F7",
	119: "F8",
	120: "F9",
	121: "F10",
	122: "F11",
	123: "F12",
	186: "semicolon",
	187: "equalsign",
	188: "comma",
	189: "dash",
	190: "period",
	191: "forwardslash",
	192: "graveaccent",
	219: "openbracket",
	220: "backslash",
	221: "closebracket",
	222: "singlequote",
}

var blacklist = map[uint16]string{
	0:  "error",
	8:  "backspace",
	9:  "tab",
	12: "clear",
	13: "enter",
	16: "shift",
	17: "ctrl",
	18: "alt",
	20: "caps lock",
	21: "hangul",
	25: "hanja",
	27: "escape",
	32: "spacebar",
	42: "print",
	44: "Print Screen",
	47: "help",
	1:  "left mouse",
	2:  "right mouse",
	91: "left window key",
	92: "rightwindowkey",
}

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
