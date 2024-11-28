import { WebSocket } from 'ws';
import { MyBot } from '../src/bot.js';
import { JDISDecoder } from './decoder.js';
import { MessageType } from '../core/message.js';

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
    this.#connect();
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
    const message_type = int(message[0]);
    const decoder = new JDISDecoder();
    
    switch (message_type) {
      case MessageType.GameStart:
        this.#bot.on_start(decoder.decode_map_state(message));
        break;
      case MessageType.GameState:
        const response = this.#bot.on_tick(decoder.decode_game_state(message));
        const actions = new TextEncoder().encode(encode_actions(response));

        const message = new Uint8Array(1 + actions.length);
        message.set(new Uint8Array([3]), 0);
        message.set(actions, 1);
        this.#ws.send(message.buffer);
        break;
      case MessageType.GameEnd:
        this.#bot.on_end();
        break;
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
