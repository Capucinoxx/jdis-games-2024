interface WasmFunc {
  getInformations: (buf: ArrayBuffer) => Object
}

declare global {
  export interface Window extends WasmFunc {
    Go: any;
  }
}

export { WasmFunc };