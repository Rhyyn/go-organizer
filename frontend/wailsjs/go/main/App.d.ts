// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {w32} from '../models';
import {ini} from '../models';
import {main} from '../models';

export function ActivateAction(arg1:string):Promise<void>;

export function ActivateNextChar():Promise<void>;

export function ActivatePreviousChar():Promise<void>;

export function CharSelectedEvent(arg1:w32.HWND):Promise<void>;

export function CreateConfigSection(arg1:ini.File,arg2:string):Promise<void>;

export function GetAllKeyBindings():Promise<{[key: number]: main.Keybinds}>;

export function GetDofusWindows():Promise<Array<main.WindowInfo>>;

export function GetKeycodes():Promise<main.UMap>;

export function GetStringFromKey(arg1:number):Promise<string|boolean>;

export function GoHook():Promise<void>;

export function InstallHook():Promise<void>;

export function IsWindowDofus():Promise<boolean|number>;

export function KeybindUpdatedEvent():Promise<void>;

export function PauseHook():Promise<void>;

export function ResumeHook():Promise<void>;

export function SaveCharacterList(arg1:Array<main.WindowInfo>):Promise<void>;

export function SaveKeybind(arg1:number,arg2:string,arg3:string):Promise<string>;

export function SetAlwaysOnTop():Promise<void>;

export function UninstallHook():Promise<void>;

export function UpdateDofusWindows():Promise<void>;

export function UpdateDofusWindowsOrder(arg1:Array<main.WindowInfo>):Promise<Array<main.WindowInfo>>;

export function UpdateOrganizerRunning():Promise<void>;

export function UpdateTemporaryDofusWindows(arg1:Array<main.WindowInfo>):Promise<void>;

export function WinActivate(arg1:w32.HWND):Promise<w32.HWND>;
