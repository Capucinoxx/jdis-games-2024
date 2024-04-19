# forlorn

Édition 2024 des JDIS Games organisé par le groupe JDIS de l'Université de Sherbrooke. Cette compétition réunie des particiant.e.s le temps d'une 
journée durant laquelle ils devront programmer un agent devant jouer à un jeu. 

Cette année, votre agent se fait téléporter dans un labyrinthe inconnu. En essaynt de scruter le labyrinthe, en essayant de scruter chaque recoin du labyrinthe, vous tomberez nez à nez avec d'autres

## Pour démarrer le projet

Actuellement, le projet fonctionne avec de la cache pour l'authentification et peut donc fonctionner sasn connecter externe.
```
go run main.go
```

Un exemple de bot pouvant être utiliser pour tester la communication:
```py
import websocket
import time
import traceback
import struct

# Message format: <rotation: uint32>
def message(rotation: int) -> bytes:
  data = []
  data.append(struct.pack('B', 1))
  data.append(struct.pack('<I', rotation))
  return b''.join(data)

def decode(msg: bytes) -> dict:
  if len(msg) < 6:
    return {'error': 'invalid message'}
  id = struct.unpack('c', msg[0:1])[0]
  message_type = struct.unpack('B', msg[1:2])[0]
  current_time = struct.unpack('<I', msg[2:6])[0]

  def decode_player(msg: bytes) -> dict:
    x = struct.unpack('<f', msg[0:4])[0]
    y = struct.unpack('<f', msg[4:8])[0]
    rotation = struct.unpack('<I', msg[8:12])[0]
    
    return {'x': x, 'y': y, 'rotation': rotation, 'health': int(msg[13])}

  data = {
    'id': id,
    'type': message_type,
    'time': current_time,
    'players': [decode_player(msg[i:]) for i in range(6, len(msg), 14)]
  }

  return data

try:
  ws = websocket.WebSocket()
  ws.connect('ws://localhost:8087/echo')

  while True:
    for i in range(0, 360, 10):
      opcode, msg = ws.recv_data()
      print(f'opcode: {opcode}, msg: {decode(msg)}')
      if opcode == websocket.ABNF.OPCODE_CLOSE:
        print('close !')
        break
      ws.send(message(i))
except Exception as e:
  print('close', traceback.format_exc())
```

## Structure du serveur Golang

| Répertoire  | Objectif                                                                 |
| :---------- | :----------------------------------------------------------------------- |
| `main.go`   | Point d'entrée du serveur.                                               |
| `internal/` | Contient le code propre à l'édition. Devra être modifié selon l'édition. |
| `pkg/`      | Contient les différents bout de code réutilisable d'année en année.      |

### Liste des différents package pouvant être réutilisé

| Répertoire       | Objectif                                                                                                                                                                                      |
| :--------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `pkg/connector/` | Point d'entrée pour la communication avec des services externes à l'application. Actuellement il y a un connecteur pour `influxDB` qui sera utilisé la persistance des pointages.             |
| `pkg/manager/`   | Gestionnaires des actions, que ce soit au niveau de la communication réseau (arrivée et départ d'une connexion), la connexion par authentification ou encore la gestion avec l'environnement. |
| `pkg/model/`     | Représentation des différents objets du jeu.                                                                                                                                                  |
| `pkg/network/`   | Repreésente la couche réseau du serveur. Les headers HTTP ainsi que la mise en place du websocket.                                                                                            |
| `pkg/protocol/`  | Fonctions d'encodage et de décodage lors de la communication avec les autres instances du jeu.                                                                                                |
| `pkg/utils`      | Fonctions utilitaires pour l'affichage des journaux.                                                                                                                                          |

## TODOS:

- [ ] Fragmenter le `game_manager` en deux parties. Une restant dans le dosser `pkg/` et l'autre allant dans le dossier `internal/`
- [ ] Ajout gestionnaire collisions pour actions spécialer
- [ ] Ajouter des colliders de type projectile, trigger dans collision_manager
- [ ] Implémentation formule calcul de point
- [ ] Utilisation de `InfluxService.Write()` pour insertion des score à chaque fois qu'il augmente
- ...