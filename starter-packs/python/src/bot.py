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
          (fr)
               Cette méthode est appelée à chaque tick de jeu. Vous pouvez y définir 
               le comportement de votre bot. Elle doit retourner une liste d'actions 
               qui sera exécutée par le serveur.
          
               Liste des actions possibles:
               - MoveTo({x, y})    permet de diriger son bot, il ira a vitesse
                                   constante jusqu'à ce point.
          
               - ShootAt({x, y})   Si vous avez le fusil comme arme, cela va tirer
                                   à la coordonnée donnée.
          
               - Store([...])      Permet de storer 100 octets dans le serveur. Lors
                                   de votre reconnection, ces données vous seront
                                   redonnées par le serveur.
          
               - SwitchWeapon(id)  Permet de changer d'arme. Par défaut, votre bot
                                   n'est pas armé, voici vos choix:
                                        PlayerWeapon.PlayerWeaponNone
                                        PlayerWeapon.PlayerWeaponCanon
                                        PlayerWeapon.PlayerWeaponBlade
          
          (en)
               This method is called at each game tick. You can define your bot's behavior here. It must return a list of actions that will be executed by the server.

               Possible actions:

               - MoveTo({x, y})    Directs your bot to move to the specified point at a constant speed.

               - ShootAt({x, y})   If you have the gun equipped, it will shoot at the given coordinates.

               - Store([...])      Allows you to store 100 bytes on the server. When you reconnect, these 
                                   data will be provided to you by the server.

               - SwitchWeapon(id)  Allows you to change your weapon. By default, your bot is unarmed. Here 
                                   are your choices:
                                        PlayerWeapon.PlayerWeaponNone
                                        PlayerWeapon.PlayerWeaponCanon
                                        PlayerWeapon.PlayerWeaponBlade
          
          @param {Model.GameState} message 
          @returns{Model.Actions}
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
        