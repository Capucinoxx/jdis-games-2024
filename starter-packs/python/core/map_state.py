import json
import struct
from dataclasses import dataclass, field
from typing import List, Tuple
from enum import IntEnum


class ColliderType(IntEnum):
    Wall = 0
    Projectile = 1

@dataclass
class Point:
    x: float = 0.0
    y: float = 0.0
    
    def __str__(self):
        return json.dumps(self.__dict__)

@dataclass
class Collider:
    collider_type: ColliderType = ColliderType.Projectile
    positions: List[Point] = field(default_factory=list)

    def __str__(self):
        return json.dumps({
            'collider_type': self.collider_type,
            'positions': [json.loads(str(position)) for position in self.positions]
        })

@dataclass
class MapState:
    size: int                       = 0
    discrete_grid: List[List[int]]  = field(default_factory=list)
    walls: List[Collider]           = field(default_factory=list)
    save: bytearray                 = field(default_factory=bytearray)

    def __str__(self) -> str:
        return json.dumps({
            'size': self.size,
            'discrete_grid': self.discrete_grid,
            'walls': [json.loads(str(wall)) for wall in self.walls],
            'save': ' '.join([f'0x{byte:02x}' for byte in self.save])
        }, indent=4)
