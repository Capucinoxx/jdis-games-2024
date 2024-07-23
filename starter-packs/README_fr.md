# Magellan - JDIS Games 2024

## Mise en situation

Ferdinand Magellan était un explorateur portugais connu pour avoir dirigé la première expédition à faire le tour du monde au 16e siècle. Le but principal de l'expédition de Magellan était de trouver une route maritime vers les îles aux épices. À cette époque, les épices étaient extrêmement précieux en Europe. Pour ce diriger, ils utilisaient des astrolabes, des instruments de navigation essentiels pour mesurer l'altitude des étoiles et des planètes au-dessus de l'horizon pour ainsi résoudre des problèmes de navigation.

## Votre objectif

Tout comme Magellan, vous serez des explorateurs devant naviguer des terres inconnues pour s'approprier d'un précieux trésor. Vous connaissez sa coordonnée, mais l'emplacement des obstacles reste floue. Il y a aura deux phases dans votre aventure : exploration et ...

## Éléments de la carte

Les éléments suivants se trouveront sur la carte : 
- Murs 🧱
- Pièces 🪙
- Trésor 📦

### Grille discrète

La grille discrète contient le nombre de murs par 4 cases.

METTRE IMAGE DE LA MAP AVEC MURS ET SANS MURS

### Pointage

| Item       | Points |
| ---------- | ------ |
| Pistolet   | 15     |
| Épée       | 40     |
| Pièce      | 40     |
| Trésor     | 1200   |

## Déroulement du jeu

### Tour 1

Durée : x temps
Les agents et les pièces sont placées de manière alétoire sur la carte 🗺️

Lorsqu'un agent prend une pièce, celle-ci reaparrait de manière aléatoire sur la carte.

### Tour 2

Durée : x temps
Un trésor est placé au centre de la carte, les agents apparaitront à équidistance de déplacement du trésor. Dans ce tour, il n'y aura pas de pièces sur la carte.

## Actions

### Déplacement

Les déplacements fonctionnent par position, c'est à dire que la position envoyé est celle que l'agent devrait se déplacer vers si aucuns obstacles n'est rencontrée.  

###  Attaque

#### Pistolet

#### Épée

### Sauvegarde

## Données reçues

### État de la carte

L'état de carte est envoyée lors de la connexion au serveur. 

### État du jeu

```
GameState {
    currentTick : Integer,                          // Cycle de rafraichissement courant
    currentRound : Integer,                         // Tour courant
    Players : [                                     // Liste des informations des joueurs 
        name : String,                              // Nom du joueur
        color : String,                             // Couleur du joueur sur la carte
        health : Integer,                           // Quantité de vie du joueur
        score : Integer,                            // Nombre de points du joueur
        position : { x : Float, y : Float },        // Position actuelle
        destination : { x : Float, y : Float },     // Position de destination
        playerWeapon : Integer,                     // Type d'arme équipée, aucune(0), pistolet(1), épée(2) 
        Projectiles: [                              // Informations des projectiles du pistolet
            uid : String,                           // Identifiant du projectile
            position : { x : Float, y : Float },    // Position actuelle
            destination : { x : Float, y : Float }, // Position de destination
        ],
        Blade: {                                    // Informations sur l'épée
            start : { x : Float, y : Float },       // Position du début de l'épée
            end : { x : Float, y : Float },         // Position de fin de l'épée
            rotation : Integer                      // Angle En degrés
        }
    ],
    Coins : [                                       // Liste des pièces ou trésor
        uid : String,                               // Identifiant de la pièce
        value : Integer,                            // Valeur de la pièce ou du trésor
        position : { x : Float, y : Float }         // Position de la pièce ou du trésor
    ]
}
```