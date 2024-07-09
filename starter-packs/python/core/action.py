from dataclasses import dataclass
from typing import Tuple

from core.game_state import PlayerWeapon

@dataclass
class MoveAction:

    dest_pos: Tuple

    def __init__(self, dest_pos: Tuple):
        self.dest_pos = dest_pos

    
    def serialize(self) -> dict:
        return {
            'dest': {
                'x': self.dest_pos[0],
                'y': self.dest_pos[1]
            }
        }


@dataclass
class ShootAction:

    target_pos: Tuple

    def __init__(self, target_pos: Tuple):
        self.target_pos = target_pos


    def serialize(self) -> dict:
        return {
            'shoot': {
                'x': self.target_pos[0],
                'y': self.target_pos[1]
            }
        }
    

@dataclass
class SwitchWeaponAction:

    weapon: PlayerWeapon

    def __init__(self, weapon: PlayerWeapon):
        self.weapon = weapon


    def serialize(self) -> dict:
        return {
            'switch': self.weapon.value
        }

