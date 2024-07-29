from dataclasses import dataclass
from typing import Tuple
import struct
import uuid

from core.game_state import PlayerInfo, PlayerWeapon, Projectile, Blade, GameState, Coin
from core.map_state import Point, Collider, ColliderType, MapState


def read_str(byte_array, end_index=None):
    if end_index is None:
        end_index = byte_array.find(b'\0')
        if end_index == -1:
            return None 
    
    string = byte_array[:end_index].decode('utf-8')
    return string, end_index


def read_uuid(byte_array, end_index):
    return str(uuid.UUID(bytes=byte_array[:end_index]))


@dataclass
class JDISDecoder:
    def decode_point(self, data: bytes, offset: int) -> Tuple[Point, int]:
        p = Point()
        p.x, p.y = struct.unpack_from('<dd', data[offset:], 0)
        return p, offset + 16


    def decode_colliders(self, pos_size: int, data: bytes, offset: int) -> Tuple[Collider, int]:
        c = Collider()
        
        for i in range(pos_size):
            p = Point()
            p, _ = self.decode_point(data, i * 16)
            c.positions.append(p)

        offset += pos_size * 16
        c.collider_type = ColliderType(struct.unpack_from('<B', data, offset)[0])
        offset += 1
        
        return c, offset
    

    def decode_map_state(self, data: bytes) -> MapState:
        m = MapState()
        m.size = struct.unpack_from('<B', data, 0)[0]
        m.discrete_grid = []
        m.spawns = []
        m.walls = []

        # decode discrete grid
        m.discrete_grid = [
            list(struct.unpack_from('<' + 'B' * m.size, data, i * m.size + 1)) for i in range(m.size)
        ]

        # decode walls
        offset = m.size * m.size + 1
        walls_len = struct.unpack_from('<i', data, offset)[0]
        
        offset += 4
        for _ in range(walls_len):
            pos_size = struct.unpack_from('<B', data, offset)[0]
            offset += 1

            collider, offset = self.decode_colliders(pos_size, data, offset)

            m.walls.append(collider)

        m.save = bytearray(data[offset: offset + 100])

        return m


    def decode_player_info(self, data: bytes) -> Tuple[PlayerInfo, int]:
        p = PlayerInfo()

        p.name, end_index = read_str(data)
        offset = end_index + 1

        p.color, p.health, p.score = struct.unpack_from('<iiq', data, offset)
        offset += 16

        p.pos, offset = self.decode_point(data, offset)
        
        has_dest = struct.unpack_from('<?', data, offset)[0]
        offset += 1

        if has_dest:
            p.dest, offset = self.decode_point(data, offset)

        p.playerWeapon = PlayerWeapon(struct.unpack_from('<B', data, offset)[0])
        offset += 1

        projectile_size = struct.unpack_from('<i', data, offset)[0]
        offset += 4

        p.projectiles = []
        for _ in range(projectile_size):
            projectile = Projectile()
            projectile.uid = read_uuid(data[offset:], 16)
            offset += 16

            projectile.pos, offset = self.decode_point(data, offset)
            projectile.dest, offset = self.decode_point(data, offset)

            p.projectiles.append(projectile)

        p.blade = Blade()

        p.blade.start, offset = self.decode_point(data, offset)
        p.blade.end, offset = self.decode_point(data, offset)

        p.blade.rotation = struct.unpack_from('<d', data, offset)[0]
        offset += 8

        return p, offset
    

    def decode_game_state(self, data: bytes) -> Tuple[GameState, int]:
        g = GameState()
        g.current_tick, g.current_round = struct.unpack_from('<ib', data, 0)
        offset = 5

        player_size = struct.unpack_from('<i', data, offset)[0]
        offset += 4
        
        g.players = []
        for _ in range(player_size):
            player, n_offset = self.decode_player_info(data[offset:])
            offset += n_offset

            g.players.append(player)

        coin_size = struct.unpack_from('<i', data, offset)[0]
        offset += 4

        g.coins = []
        for _ in range(coin_size):
            coin = Coin()
            coin.uid = read_uuid(data[offset:], 16)
            offset += 16

            coin.pos, offset = self.decode_point(data, offset)

            coin.value = struct.unpack_from('<i', data, offset)[0]
            offset += 4

            g.coins.append(coin)

        return g
