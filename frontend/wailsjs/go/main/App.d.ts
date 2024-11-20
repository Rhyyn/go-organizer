// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {ini} from '../models';
import {main} from '../models';

export function CreateConfigSection(arg1:ini.File,arg2:string):Promise<void>;

export function GetAllKeyBindings():Promise<{[key: string]: any}>;

export function GetDofusWindows():Promise<Array<main.WindowInfo>>;

export function GetKeycodes():Promise<main.UMap>;

export function GetToggleListenerKeybind():Promise<string>;

export function PauseHook():Promise<void>;

export function ResumeHook():Promise<void>;

export function SaveCharacterList(arg1:Array<main.WindowInfo>):Promise<void>;

export function SaveNextCharKeybind(arg1:number,arg2:string):Promise<string>;

export function SavePreviousCharKeybind(arg1:number,arg2:string):Promise<string>;

export function SaveStopOrgaKeyBind(arg1:number,arg2:string):Promise<string>;

export function UpdateDofusWindows():Promise<void>;

export function UpdateDofusWindowsOrder(arg1:Array<main.WindowInfo>):Promise<Array<main.WindowInfo>>;

export function UpdateMainHookState():Promise<void>;

export function UpdateOrder(arg1:Array<string>,arg2:Array<string>):Promise<Array<string>>;
