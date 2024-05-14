import websocket
import time
import traceback
import struct

# Message format: <rotation: uint32>
def message(rotation: int) -> bytes:
  data = []
  data.append(struct.pack('B', 1))
  data.append(struct.pack('<I', rotation))
  return b''.join(data)

def decode(msg: bytes) -> dict:
  if len(msg) < 6:
    return {'error': 'invalid message'}
  id = struct.unpack('c', msg[0:1])[0]
  message_type = struct.unpack('B', msg[1:2])[0]
  current_time = struct.unpack('<I', msg[2:6])[0]

  def decode_player(msg: bytes) -> dict:
    x = struct.unpack('<f', msg[0:4])[0]
    y = struct.unpack('<f', msg[4:8])[0]
    rotation = struct.unpack('<I', msg[8:12])[0]
    
    return {'x': x, 'y': y, 'rotation': rotation, 'health': int(msg[13])}

  data = {
    'id': id,
    'type': message_type,
    'time': current_time,
    'players': [decode_player(msg[i:]) for i in range(6, len(msg), 14)]
  }

  return data

try:
  ws = websocket.WebSocket()
  headers = {
    "Authorization": "1234"
  }
  ws.connect('ws://localhost:8087/echo', header=headers)

  while True:
    for i in range(0, 360, 10):
      opcode, msg = ws.recv_data()
      print(f'opcode: {opcode}, msg: {decode(msg)}')
      if opcode == websocket.ABNF.OPCODE_CLOSE:
        print('close !')
        break
      ws.send(message(i))
except Exception as e:
  print('close', traceback.format_exc())