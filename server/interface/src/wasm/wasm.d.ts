import '../types/index.d.ts';

interface WasmFunc {
  getInformations: (buf: ArrayBuffer) => ServerMessage;
}

declare global {
  export interface Window extends WasmFunc {
    Go: any;
  }
}

export { WasmFunc };
