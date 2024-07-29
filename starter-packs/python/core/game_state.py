import json
from dataclasses import dataclass, field
from typing import List, Optional
from enum import IntEnum

from core.map_state import Point


@dataclass
class Projectile:
    uid: str    = ''
    pos: Point  = field(default_factory=Point)
    dest: Point = field(default_factory=Point)

    def __str__(self) -> str:
        return json.dumps(self.__dict__)


@dataclass
class Blade:
    start: Point    = field(default_factory=Point)
    end: Point      = field(default_factory=Point)
    rotation: float = 0.0

    def __str__(self) -> str:
        return json.dumps(self.__dict__)


@dataclass
class Coin:
    uid: str    = ''
    value: int  = 0
    pos: Point  = field(default_factory=Point)

    def __str__(self) -> str:
        return json.dumps(self.__dict__)


class PlayerWeapon(IntEnum):
    PlayerWeaponNone    = 0
    PlayerWeaponCanon   = 1
    PlayerWeaponBlade   = 2

    def __str__(self) -> str:
        return self.name

@dataclass
class PlayerInfo:
    name: str                       = ''
    color: int                      = 0
    health: int                     = 0
    score: int                      = 0
    pos: Point                      = field(default_factory=Point)
    dest: Point                     = field(default_factory=Point)
    playerWeapon: PlayerWeapon      = PlayerWeapon.PlayerWeaponNone
    projectiles: List[Projectile]   = field(default_factory=list)
    blade: Blade                    = field(default_factory=Blade)

    def isAlive(self) -> bool:
        return self.health > 0

    def __str__(self) -> str:
        return json.dumps(self.__dict__, default=lambda o: o.__dict__, indent=4)


@dataclass
class GameState:
    current_tick: int           = 0
    current_round: int          = 0
    players: List[PlayerInfo]   = field(default_factory=list)
    coins: List[Coin]           = field(default_factory=list)

    def __str__(self) -> str:
        return json.dumps(self.__dict__, default=lambda o: o.__dict__, indent=4)
