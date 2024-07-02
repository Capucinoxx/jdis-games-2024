import websocket
import json

import struct

def decode(message):
    # little indian byte order
    struct.unpack_from('B', message, 0)

def on_message(ws, message):
    print("Message received from server: ", message)
    ws.send(send_message())

def on_error(ws, error):
    print("Error: ", error)

def on_close(ws, close_status_code, close_msg):
    print("Connection closed")

def on_open(ws):
    print("Connection opened")
    ws.send(send_message())
    print("Message sent")


def send_message():
    json_message = json.dumps({
        'dest': { 'x': 10.0, 'y': 11.34 },
        'shoot': { 'x': 11.2222, 'y': 13.547 }
    })
    # Préfixer le message avec le byte ayant la valeur 3
    prefixed_message = bytearray([3]) + json_message.encode('utf-8')
    return prefixed_message

websocket.enableTrace(True)
ws = websocket.WebSocketApp('ws://localhost:8087/echo',
                            header={'Authorization': 'SECRET-TOKEN'},
                            on_open=on_open,
                            on_message=on_message,
                            on_error=on_error,
                            on_close=on_close)

ws.run_forever()