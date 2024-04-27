// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {main} from '../models';

export function GetPorts():Promise<main.ComPorts>;

export function GetReceiveData():Promise<Array<main.ReceiveData>>;

export function OpenSerialPort():Promise<boolean>;

export function SendData(arg1:string):Promise<main.SendData>;

export function SerialClose():Promise<void>;

export function SetDelimiter(arg1:main.DelimiterSetting):Promise<void>;

export function SetPortSetting(arg1:main.PortSetting):Promise<void>;