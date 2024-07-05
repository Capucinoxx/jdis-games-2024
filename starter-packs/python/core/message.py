from enum import IntEnum

class MessageType(IntEnum):
    """
    (fr)
	GameStart est un message envoyé pour informer les clients que la partie commence.
    position = gamestate
    """
    GameState = 1

    """
    (fr)
	GameEnd est un message envoyé pour informer les clients que la partie est terminée.
    """

    GameStart = 4 # map state meme chose 

    GameEnd = 5

