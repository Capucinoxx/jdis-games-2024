from dataclasses import dataclass
from typing import List, Tuple

from core.map_state import Point
from utils.utils import read_string_until_null as read_str

import struct

@dataclass
class Projectile:
    uid: str
    pos: Point
    dest: Point

@dataclass
class Blade:
    start: Point
    end: Point
    rotation: float

@dataclass
class Coin:
    uid: str
    value: int
    pos: Point

@dataclass
class PlayerInfo:
    name: str
    color: int
    health: int
    pos: Point
    dest: Point   
    projectiles: List[Projectile]
    blade: Blade


    def __init__(self):
        self.name = ''
        self.color = 0
        self.health = 0
        self.pos = Point()
        self.dest = Point()

    def decode(self, data:bytes):
        self.name, end_index = read_str(data)
        offset = end_index
        self.color, self.health = struct.unpack_from('<ii', data, offset)
        offset += 8

        offset += self.pos.decode(data[offset:])
        
        has_dest = struct.unpack_from('<?', data, offset)[0]
        offset += 1
    
        if has_dest:
            offset += self.dest.decode(data[offset:])

        projectile_size = struct.unpack_from('<i', data, offset)[0]
        offset += 4

        self.projectiles = []
        for i in range(projectile_size):
            projectile = Projectile()
            projectile.uid, end_index = read_str(data[offset:], 16)
            offset += end_index

            offset += projectile.pos.decode(data[offset:])

            offset += projectile.dest.decode(data[offset:])
            
            self.projectiles.append(projectile)

        self.blade = Blade()
        offset += self.blade.start.decode(data[offset:])
        offset += self.blade.end.decode(data[offset:])
        
        self.blade.rotation = struct.unpack_from('<d', data, offset)[0]
        offset += 8

        return offset

@dataclass
class GameInfo:
    current_tick: int
    current_round: int
    players: List[PlayerInfo]
    coins: List[Coin]

    @classmethod
    def decode(cls, data:bytes):
        cls.current_tick, cls.current_round = struct.unpack_from('<ib', data, 0)
        offset = 5

        player_size = struct.unpack_from('<i', data, offset)[0]
        offset += 4


        cls.players = []
        for i in range(player_size):
            player = PlayerInfo()
            offset = player.decode(data[offset:])

            cls.players.append(player)

        coin_size = struct.unpack_from('<i', data, offset)[0]
        offset += 4

        cls.coins = []
        for i in range(coin_size):
            coin = Coin()
            coin.uid, end_index = read_str(data[offset:], 16)
            offset += end_index

            coin.value = struct.unpack_from('<i', data, offset)[0]
            offset += 4

            offset += coin.pos.decode(data[offset:])

            cls.coins.append(coin)

        print(f'current_tick {cls.current_tick} current_round {cls.current_round}')
