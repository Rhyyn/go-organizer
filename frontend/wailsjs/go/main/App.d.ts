// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {w32} from '../models';
import {main} from '../models';

export function Add():Promise<void>;

export function AttemptSetForeground(arg1:w32.HWND,arg2:w32.HWND):Promise<boolean>;

export function GetDofusWindows():Promise<Array<main.WindowInfo>>;

export function SetForegroundWindowEx(arg1:w32.HWND,arg2:w32.HWND):Promise<w32.HWND>;

export function UpdateDofusWindows():Promise<Array<main.WindowInfo>>;

export function UpdateDofusWindowsOrder(arg1:Array<main.WindowInfo>):Promise<void>;

export function WinActivate(arg1:w32.HWND):Promise<w32.HWND>;
