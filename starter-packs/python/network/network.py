import websocket
import json
import ssl
import threading
import time
from typing import List, Optional

from src.bot import MyBot
from core.message import MessageType
from core.action import Action
from network.decoder import JDISDecoder


class Socket:  
    def __init__(self, url: str, token: str):
        self.url = url
        self.token = token
        self.bot = MyBot()
        self.ping_interval = 1

        
    def run(self):
        print(f"Starting bot with base URL: {self.url}, token: {self.token}")
        try:
            ws = websocket.WebSocketApp(self.url,
                                        header={'Authorization': self.token},
                                        on_open=self.on_open,
                                        on_message=self.on_message,
                                        on_error=self.on_error,
                                        on_close=self.on_close)
        except Exception as e:
            print("Error: ", e)
            return
        
        ws.run_forever(sslopt={"cert_reqs": ssl.CERT_NONE})


    def handle_message(self, message: bytes) -> Optional[List[Action]]:
        message_type = int(message[0])
        response = None

        decoder = JDISDecoder()

        if message_type == MessageType.GameStart.value:
            map_state = decoder.decode_map_state(message[1:])
            self.bot.on_start(map_state)

        elif message_type == MessageType.GameState.value:
            game_state = decoder.decode_game_state(message[1:])
            response = self.bot.on_tick(game_state)

        elif message_type == MessageType.GameState.GameEnd.value:
            self.bot.on_end()

        else:
            print("Unknown message type")

        return response


    def on_open(self, ws: websocket.WebSocketApp) -> None:
        print("Connection opened")
        self.start_ping_thread(ws)
        

    def on_message(self, ws: websocket.WebSocketApp, message: bytes) -> None:
        response = self.handle_message(message)
        if response:
            self.send_message(ws, response)
        

    def on_error(self, ws: websocket.WebSocketApp, error: str) -> None:
        print("Error: ", error)
        

    def on_close(self, ws: websocket.WebSocketApp, close_status_code, close_msg) -> None:
        print("Connection closed")
        

    def send_message(self, ws: websocket.WebSocketApp, actions: List[Action]) -> None:
        json_reponse = {}
        for action in actions:
            try:
                json_reponse.update(action.serialize())
            except Exception as e:
                print(e)
    
        json_message = json.dumps(json_reponse)

        print(f"Sending message: {json_message}")
        prefixed_message = bytearray([3]) + json_message.encode('utf-8')
        ws.send(prefixed_message)


    def start_ping_thread(self, ws: websocket.WebSocketApp) -> None:
        ping_thread = threading.Thread(target=self.ping, args=(ws,))
        ping_thread.daemon = True
        ping_thread.start()


    def ping(self, ws: websocket.WebSocketApp) -> None:
        while ws.keep_running:
            ws.send('ping')
            time.sleep(self.ping_interval)
    