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

    def decode(self, data:bytes):
        self.x, self.y = struct.unpack_from('<dd', data, 0)

        return 16
    
    def __str__(self):
        return f'Point(x={self.x}, y={self.y})'
        
@dataclass
class Collider:
    collider_type: ColliderType
    positions: List[Point]

    def __init__(self):
        self.positions = []
        self.collider_type = None

    def decode(self, pos_size:int, data:bytes):
        # print(f'pos_size {pos_size}')
        for i in range(pos_size):
            p = Point()
            p.decode(data[i * 16:])
            # print(p)
            self.positions.append(p)
        self.collider_type = ColliderType(struct.unpack_from('<B', data, pos_size * 16)[0])
    
    def __str__(self):
        str_positions = ' '.join(str(pos) for pos in self.positions)

        return f'{self.collider_type.name} {str_positions.__str__(self)}'

@dataclass
class MapState:
    discrete_grid: List[List[int]]
    size: int

    spawns: List[Tuple[Point]]
    walls: List[Collider]

    @classmethod
    def decode(cls, data:bytes):
        cls.size = struct.unpack_from('<B', data, 0)[0]
        cls.discrete_grid = []
        cls.spawns = []
        cls.walls = []

        # decode discrete grid
        cls.discrete_grid = [
            list(struct.unpack_from('<' + 'B' * cls.size, data, i * cls.size + 1)) for i in range(cls.size)
        ]

        # decode walls
        offset = cls.size * cls.size + 1
        walls_len = struct.unpack_from('<i', data, offset)[0]
        # cls.walls = [Collider.decode(data[cls.size * cls.size + 5 + i * 16:]) for i in range(walls_len)]
        
        offset += 4
        for i in range(walls_len):
            pos_size = struct.unpack_from('<B', data, offset)[0]
            offset += 1
            collider = Collider()
            collider.decode(pos_size, data[offset:])
            offset += pos_size * 16 + 1
            cls.walls.append(collider)
            # print(wall.__str__(wall) + ' ')

        # print(cls.walls.__str__())
        # print(cls.walls)

        # # decode walls
        # walls_size = struct.unpack_from('<B', data[1:], self.size * self.size + 1 + spawns_size * 2)
        # for i in range(walls_size):
        #     wall = Collider(ColliderType.Wall,
        #                     Point(struct.unpack_from('<f', data[1:], self.size * self.size + 1 + spawns_size * 2 + 1 + i * 3),
        #                           struct.unpack_from('<f', data[1:], self.size * self.size + 1 + spawns_size * 2 + 1 + i * 3 + 1)))
        #     self.walls.append(wall)