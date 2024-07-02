import websocket
import json

from src.bot import MyBot

class Socket:  
    def __init__(self, url: str, token: str):
        self.url = url
        self.token = token
        self.bot = MyBot()

        
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
        
        ws.run_forever()
        
    def on_open(self, ws):
        print("Connection opened")
        self.send_message(ws)
        print("Message sent")
        
    def on_message(self, ws, message):
        print("Message received from server: ", message)
        self.bot.tick()
        self.send_message(ws)
        
    def on_error(self, ws, error):
        print("Error: ", error)
        
    def on_close(self, ws, close_status_code, close_msg):
        print("Connection closed")
        
    def send_message(self, ws):
        json_message = json.dumps({
            'dest': { 'x': 10.0, 'y': 11.34 },
            'shoot': { 'x': 11.2222, 'y': 13.547 }
        })
        # Préfixer le message avec le byte ayant la valeur 3
        prefixed_message = bytearray([3]) + json_message.encode('utf-8')

        ws.send(prefixed_message)