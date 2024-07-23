from dataclasses import dataclass
from typing import List
from enum import IntEnum

from core.map_state import Point

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
    score: int
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
    

    def __str__(self) -> str:
        return f'PlayerInfo(name={self.name}, color={self.color}, health={self.health}, pos={self.pos}, dest={self.dest}, playerWeapon={self.playerWeapon.Name}, projectiles={self.projectiles}, blade={self.blade})'


@dataclass
class GameState:
    current_tick: int
    current_round: int
    players: List[PlayerInfo]
    coins: List[Coin]

    def __init__(self):
        self.current_tick = 0
        self.current_round = 0
        self.players = []
        self.coins = []
