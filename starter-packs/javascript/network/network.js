import { WebSocket } from 'ws';
import { MyBot } from '../src/bot';

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
            headers: { 'Authorization': this.#secret }
        });
        this.#ws.binaryType = 'arraybuffer';

        this.#ws.on('open', () => this.#on_open());  
        this.#ws.on('on_message', (message) => this.#on_message(message));
        this.#ws.on('close', () => this.#on_close());
        this.#ws.on('error', (error) => this.#on_error(error)); 
    }


    #on_open() {
        console.log(`Connected to ${this.#url}`);
        this.#start_heartbeat();
    }


    #on_message(message) {
        try {
            
        } catch(error) {}
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

    actions.forEach(action => {
        if (action === null)
            return;
        
        switch (action.type) {
            case 'move_to':
                data[action.type] = action.dest;
                break;
            case 'shoot_at':
                data[action.type] = action.pos;
                break;
            case 'store':
                data[action.type] = action.data;
                break;
            case 'switch':
                data[action.type] = action.weapon;
                break;
            default:
                break;
        }
    });

    return JSON.stringify(data);
};