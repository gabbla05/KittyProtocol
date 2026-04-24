export type FrameType = "HELLO" | "AUTH" | "DATA" | "MEOW_OK" | "ERROR" | "GET_STATUS" | "STATUS_RES" | "PING" | "BYE";

export interface BaseFrame {
  type: FrameType;
  msg_id: number;
}

export interface HelloFrame extends BaseFrame {
  type: "HELLO";
  version: string;
}

export interface AuthFrame extends BaseFrame {
  type: "AUTH";
  user: string;
  pass: string;
}

export interface DataFrame extends BaseFrame {
  type: "DATA";
  target?: string;
  sender?: string;
  payload: string;
  mac: string;
}

export interface ErrorFrame extends BaseFrame {
  type: "ERROR";
  code: string;
  desc: string;
}