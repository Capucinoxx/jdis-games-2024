from dataclasses import dataclass
from typing import List, Tuple

import struct
from enum import IntEnum

class ColliderType(IntEnum):
    Wall = 0
    Projectile = 1


@dataclass
class Point:
    x: float
    y: float

    def __init__(self):
        self.x = 0
        self.y = 0
    
    def __str__(self):
        return f'Point(x={self.x}, y={self.y})'
        
@dataclass
class Collider:
    collider_type: ColliderType
    positions: List[Point]

    def __init__(self):
        self.positions = []
        self.collider_type = None

    def __str__(self):
        str_positions = ' '.join(str(pos) for pos in self.positions)

        return f'{self.collider_type.name} {str_positions.__str__(self)}'


@dataclass
class MapState:
    discrete_grid: List[List[int]]
    size: int

    spawns: List[Tuple[Point]]
    walls: List[Collider]

    save: bytearray

    def __init__(self):
        self.discrete_grid = []
        self.size = 0
        self.spawns = []
        self.walls = []
        self.save = bytearray()
