import './wasm.d.ts';
import './wasm_exec.js';
import { URL_BASE } from '../config';


const load_wasm = async (): Promise<void> => {
  const wasm = new window.Go();
  

  const result = await WebAssembly.instantiateStreaming(fetch(`https://${URL_BASE}/lib.wasm`), wasm.importObject);
  wasm.run(result.instance);
}

export { load_wasm };