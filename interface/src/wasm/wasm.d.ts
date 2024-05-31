

type ServerMapState = {
  Type: 4;
  map: Array<Array<number>>;
};

type Empty = Record<string, never>;


type ServerMessage = ServerMapState | Empty;

interface WasmFunc {
  getInformations: (buf: ArrayBuffer) => ServerMessage;
}

declare global {
  export interface Window extends WasmFunc {
    Go: any;
  }
}

export { WasmFunc };
