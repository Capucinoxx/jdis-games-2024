from typing import List, Union

from core.action import MoveAction, ShootAction, RotateBladeAction, SwitchWeaponAction, SaveAction
from core.consts import Consts
from core.game_state import GameState, PlayerWeapon
from core.map_state import MapState


class MyBot:
    """
    (fr) Cette classe représente votre bot. Vous pouvez y définir des attributs et des méthodes qui 
         seront conservées entre chaque appel de la méthode `0n_tick`.

    (en) This class represents your bot. You can define attributes and methods in it that will be kept 
         between each call of the `0n_tick` method.
    """

    __map_state: MapState


    def on_tick(self, game_state: GameState) -> List[Union[MoveAction, SwitchWeaponAction, RotateBladeAction, ShootAction, SaveAction]]:
        """
        (fr) Cette méthode est appelée à chaque tick de jeu. Vous pouvez y définir le comportement de
             votre bot. Elle doit retourner une liste d'actions qui sera exécutée par le serveur.

        (en) This method is called every game tick. You can define the behavior of your bot. It must 
            return a list of actions which will be executed by the server.

        Arguments:
            game_state (GameState): (fr) L'état du jeu.
                                    (en) The state of the game.
        """
        print(game_state)

        actions = [
            MoveAction((10.0, 11.34)), 
            ShootAction((11.2222, 13.547)),
            SwitchWeaponAction(PlayerWeapon.PlayerWeaponBlade),
            SaveAction(b'Hello, world!'),
        ]

        return actions
    
    
    def on_start(self, map_state: MapState):
        """
        (fr) Cette méthode est appelée une seule fois au début de la partie. Vous pouvez y définir des
             actions à effectuer au début de la partie.

        (en) This method is called once at the beginning of the game. You can define actions to be 
             performed at the beginning of the game.

        Arguments:
            map_state (MapState): (fr) L'état de la carte.
                                  (en) The state of the map.
        """
        self.__map_state = map_state
        pass


    def on_end(self):
        """
        (fr) Cette méthode est appelée une seule fois à la fin de la partie. Vous pouvez y définir des
             actions à effectuer à la fin de la partie.

        (en) This method is called once at the end of the game. You can define actions to be performed 
             at the end of the game.
        """
        pass
        