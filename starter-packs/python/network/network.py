import websocket
import json

class Socket:  
    def __init__(self, url: str, token: str):
        self.url = url
        self.token = token
        self.ws = websocket.WebSocketApp(url,
                                         header={'Authorization': token},
                                         on_open=self.on_open,
                                         on_message=self.on_message,
                                         on_error=self.on_error,
                                         on_close=self.on_close)
        
    def run(self):
        self.ws.run_forever()
        
    def on_open(self, ws):
        print("Connection opened")
        ws.send(self.send_message())
        print("Message sent")
        
    def on_message(self, ws, message):
        print("Message received from server: ", message)
        ws.send(self.send_message())
        
    def on_error(self, ws, error):
        print("Error: ", error)
        
    def on_close(self, ws, close_status_code, close_msg):
        print("Connection closed")
        
    def send_message(self):
        json_message = json.dumps({
            'dest': { 'x': 10.0, 'y': 11.34 },
            'shoot': { 'x': 11.2222, 'y': 13.547 }
        })
        # Préfixer le message avec le byte ayant la valeur 3
        prefixed_message = bytearray([3]) + json_message.encode('utf-8')

        return prefixed_message