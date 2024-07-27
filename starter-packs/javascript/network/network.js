import './polyfill.js';
import { promises as fs } from 'fs';

import { WebSocket } from 'ws';
import { MyBot } from '../src/bot.js';
import './wasm_exec.js';


const load_wasm = async () => {
    const wasmBuffer = await fs.readFile('network/lib.wasm');
    
    const go = new Go();
    const { instance } = await WebAssembly.instantiate(wasmBuffer, go.importObject);
    
    go.run(instance);
  };

class Socket {
    #url;
    #secret;
    #ws;
    #ping_interval = null;
    #bot = null;

    constructor(url, secret) {
        this.#url = url;
        this.#secret = secret;
        this.#bot = new MyBot();
    }


    run() {
        load_wasm().then(() => {
            console.log('WebAssembly loaded and executed');
            this.#connect();
          }).catch(err => {
            console.error('Error loading WebAssembly:', err);
          });
    }


    #connect() {
        this.#ws = new WebSocket(this.#url, {
            headers: { 'Authorization': this.#secret },
            rejectUnauthorized: false
        });
        this.#ws.binaryType = 'arraybuffer';

        this.#ws.on('open', () => this.#on_open());  
        this.#ws.on('message', (message) => this.#on_message(message));
        this.#ws.on('close', () => this.#on_close());
        this.#ws.on('error', (error) => this.#on_error(error)); 
    }


    #on_open() {
        console.log(`Connected to ${this.#url}`);
        this.#start_heartbeat();
    }


    #on_message(message) {
        try {
            const data = global.getInformations(message);
            if (!('type' in data))
                return;

            switch (data.type) {
                case 4:
                    this.#bot.on_start({ map: data.map, walls: data.walls, size: data.size, save: data.save });
                    break;
                
                case 5:
                    this.#bot.on_end();
                    break;

                case 1:
                    const actions = this.#bot.on_tick({ tick: data.tick, round: data.round, players: data.players, coins: data.coins });
                    const message = encode_actions(actions);
                    console.log(`Sending message: ${message}`);
                    const prefix = new Uint8Array([3]);
                    const buffer = new TextEncoder().encode(message);

                    const prefixed_message = new Uint8Array(prefix.length + buffer.length);
                    prefixed_message.set(prefix, 0);
                    prefixed_message.set(buffer, prefix.length);
                    this.#ws.send(prefixed_message.buffer);
                    break;
            }
        } catch(error) {
            console.error(error);
        }
    }


    #on_error(error) {
        console.log(`Websocket error: ${error}`);
        this.#stop_heartbeat();
    }


    #on_close() {
        console.log('Websocket connection closed');
        this.#stop_heartbeat();
    }


    #start_heartbeat() {
        this.#ping_interval = setInterval(() => {
            if (this.#ws.readyState === WebSocket.OPEN)
                this.#ws.ping();
        }, 1000);
    }

    
    #stop_heartbeat() {
        if (this.#ping_interval) {
            clearInterval(this.#ping_interval);
            this.#ping_interval = null;
        }
    }
};

/**
 * 
 * @param {Model.Actions} actions 
 */
const encode_actions = (actions) => {
    const data = {};

    actions && actions.forEach(action => {
        if (action === null)
            return;
        
        switch (action.type) {
            case 'dest':
                data[action.type] = action.destination;
    
                break;
            case 'shoot':
                data[action.type] = action.pos;
                break;
            case 'save':
                data[action.type] = btoa(String.fromCharCode.apply(null, action.data));
                break;
            case 'switch':
                data[action.type] = action.weapon;
                break;
            case 'rotate_blade':
                data[action.type] = action.rad;
            default:
                break;
        }
    });

    return JSON.stringify(data);
};


export { Socket };