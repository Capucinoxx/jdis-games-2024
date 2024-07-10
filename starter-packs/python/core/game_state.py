from dataclasses import dataclass
from typing import List
from enum import IntEnum

from core.map_state import Point
from utils.utils import read_string_until_null as read_str, read_uuid

import struct

@dataclass
class Projectile:
    uid: str
    pos: Point
    dest: Point

    def __init__(self):
        self.uid = ''
        self.pos = Point()
        self.dest = Point()

    def __str__(self) -> str:
        return f'Projectile(uid={self.uid}, pos={self.pos}, dest={self.dest})'

@dataclass
class Blade:
    start: Point
    end: Point
    rotation: float

    def __init__(self):
        self.start = Point()
        self.end = Point()
        self.rotation = 0

    def __str__(self) -> str:
        return f'Blade(start={self.start}, end={self.end}, rotation={self.rotation})'

@dataclass
class Coin:
    uid: str
    value: int
    pos: Point

    def __init__(self):
        self.uid = ''
        self.value = 0
        self.pos = Point()

    def __str__(self) -> str:
        return f'Coin(uid={self.uid}, value={self.value}, pos={self.pos})'


class PlayerWeapon(IntEnum):
    PlayerWeaponNone = 0
    PlayerWeaponCanon = 1
    PlayerWeaponBlade = 2


@dataclass
class PlayerInfo:
    name: str
    color: int
    health: int
    pos: Point
    dest: Point   
    playerWeapon : PlayerWeapon
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
        offset = end_index + 1

        self.color, self.health = struct.unpack_from('<ii', data, offset)
        offset += 8

        offset += self.pos.decode(data[offset:])
        
        has_dest = struct.unpack_from('<?', data, offset)[0]
        offset += 1

        if has_dest:
            offset += self.dest.decode(data[offset:])

        self.playerWeapon = PlayerWeapon(struct.unpack_from('<B', data, offset)[0])
        offset += 1

        projectile_size = struct.unpack_from('<i', data, offset)[0]
        offset += 4

        self.projectiles = []
        for i in range(projectile_size):
            projectile = Projectile()
            projectile.uid = read_uuid(data[offset:], 16)
            offset += 16

            offset += projectile.pos.decode(data[offset:])
            offset += projectile.dest.decode(data[offset:])

            self.projectiles.append(projectile)

        self.blade = Blade()
        offset += self.blade.start.decode(data[offset:])
        offset += self.blade.end.decode(data[offset:])

        self.blade.rotation = struct.unpack_from('<d', data, offset)[0]
        offset += 8

        return offset
    
    def __str__(self) -> str:
        return f'PlayerInfo(name={self.name}, color={self.color}, health={self.health}, pos={self.pos}, dest={self.dest}, playerWeapon={self.playerWeapon.Name}, projectiles={self.projectiles}, blade={self.blade})'


@dataclass
class GameState:
    current_tick: int
    current_round: int
    players: List[PlayerInfo]
    coins: List[Coin]

    @classmethod
    def print(cls) -> str:
        print(f"""
GameState(
    current_tick={cls.current_tick}
    current_round={cls.current_round}
    Players=[{cls.players}]
    Coins=[{cls.coins}]
)
        """)

    @classmethod
    def decode(cls, data:bytes):
        cls.current_tick, cls.current_round = struct.unpack_from('<ib', data, 0)
        offset = 5

        player_size = struct.unpack_from('<i', data, offset)[0]
        offset += 4
        
        cls.players = []
        for i in range(player_size):
            player = PlayerInfo()
            offset += player.decode(data[offset:])

            cls.players.append(player)

        coin_size = struct.unpack_from('<i', data, offset)[0]
        offset += 4

        cls.coins = []
        for i in range(coin_size):
            coin = Coin()
            coin.uid = read_uuid(data[offset:], 16)
            offset += 16

            offset += coin.pos.decode(data[offset:])

            coin.value = struct.unpack_from('<i', data, offset)[0]
            offset += 4

            cls.coins.append(coin)

        return cls