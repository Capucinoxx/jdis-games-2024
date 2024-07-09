# Pour les tests de la plateforme

À noter, le code du starter pack n'est pas final.

## Lancer votre bot
```
python run_bot.py --token <VOTRE TOKEN>
```

Vous pourrez ajouter la logique de votre code dans [src/bot.py](src/bot.py). C'est le seul fichier que vous avez besoin de modifier.

## Données reçues du serveur

Vous recevez le `MapState` au début de chaque game et le `GameState` à chaque tick. Vous pouvez voir les dataclass respectives dans [core/map_state.py](core/map_state.py) et [core/game_state.py](core/game_state.py).

## Données renvoyées au serveur

Vous pourrez envoyer des actions au serveur à chaque tick. 
Actions disponibles :
- ShootAction(Tuple(x, y))
- MoveAction(Tuple(x, y))
- SwitchWeaponAction(PlayerWeapon)
- SaveAction(bytearray)
