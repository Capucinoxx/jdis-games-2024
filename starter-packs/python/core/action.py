from dataclasses import dataclass
from typing import Tuple, Union

from core.game_state import PlayerWeapon

import base64


@dataclass
class MoveAction:
    """
    (fr) Représente une action de déplacement vers une position spécifiée.
    (en) Represents an action to move to a specified position.

    Attributes:
        dest_pos (Tuple[int, int]) : (fr) La position de destination sous forme de coordonnées (x, y).
                                     (en) The destination position as (x, y) coordinates.
    """

    dest_pos: Tuple[int, int]

    def __init__(self, dest_pos: Tuple[int, int]):
        self.dest_pos = dest_pos

    def serialize(self) -> dict:
        return {"dest": {"x": self.dest_pos[0], "y": self.dest_pos[1]}}


@dataclass
class ShootAction:
    """
    (fr) Représente une action de tir vers la position spécifiée. Cela nécessite d'avoir
         le canon d'équipé (PlayerWeapon.PlayerWeaponCanon).
    (en) Represents an action to shoot at a specified target position. Requires the canon
         (PlayerWeapon.PlayerWeaponCanon) to be equiped.

    Attributes:
        target_pos (Tuple[int, int]) : (fr) La position cible sous forme de coordonnées (x, y).
                                       (en) The target position as (x, y) coordinates.
    """

    target_pos: Tuple

    def __init__(self, target_pos: Tuple):
        self.target_pos = target_pos

    def serialize(self) -> dict:
        return {"shoot": {"x": self.target_pos[0], "y": self.target_pos[1]}}


@dataclass
class SwitchWeaponAction:
    """
    (fr) Représente une action pour changer d'arme. Cette arme reste équipée tant que l'on ne
         la change pas. En utilisant cette action, on ne peut pas utiliser ShootAction et RotateBladeAction dans le même tick.
    (en) Represents an action to switch a specified weapon. This weapon remains equiped as long
         you don't changed. When using this action, ShootAction and RotateBladeAction cannot be used in the same tick.

    Attributes:
        weapon (PlayerWeapon) : (fr) L'arme a équiper.
                                (en) The weapon to switch to.
    """

    weapon: PlayerWeapon

    def __init__(self, weapon: Union[PlayerWeapon, int]):
        if isinstance(weapon, int):
            try:
                self.weapon = PlayerWeapon(weapon)
            except ValueError:
                raise ValueError(f"Invalid value for PlayerWeapon: {weapon}")
        elif isinstance(weapon, PlayerWeapon):
            self.weapon = weapon
        else:
            raise TypeError("weapon must be an instance of PlayerWeapon or int")

    def serialize(self) -> dict:
        return {"switch": self.weapon.value}


@dataclass
class SaveAction:
    """
    (fr) Représente une action pour sauvegarder l'état actuel du jeu. Vous disposez d'un total de 100 octets.
         Personne d'autre n'aura accès à vos données.
    (en) Represents an action to save the current game state. A total of 100 bytes is available. No one else
         will have access to your stored data.

    Attributes:
        save (bytes) : (fr) L'état du jeu à sauvegarder, au format de bytes.
                       (en) The game state to be saved, in bytes format.
    """

    save: bytes

    def serialize(self) -> dict:
        return {"save": base64.b64encode(self.save).decode("utf-8")}


@dataclass
class RotateBladeAction:
    """
    (fr) Représente une action pour faire pivoter la lame d'un angle spécifié en radians. Cela
         nécessite d'avoir la lame d'équipée (PlayerWeapon.PlayerWeaponBlade)
    (en) Represents an action to rotate the blade by a specified angle in radians. Requires the blade
         (PlayerWeapon.PlayerWeaponBlade) to be equiped.

    Attributes:
        rad (bytes) : (fr) L'angle en radians pour faire pivoter la lame.
                      (en) The angle in radians to rotate the blade.
    """

    rad: float

    def serialize(self) -> dict:
        return {"rotate_blade": self.rad}


Action = Union[
    MoveAction, ShootAction, SwitchWeaponAction, SaveAction, RotateBladeAction
]
