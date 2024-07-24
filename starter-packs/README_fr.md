# Magellan - JDIS Games 2024

## Mise en situation

Ferdinand Magellan était un explorateur portugais connu pour avoir dirigé la première expédition à faire le tour du monde au 16e siècle. Le but principal de l'expédition de Magellan était de trouver une route maritime vers les îles aux épices. À cette époque, les épices étaient extrêmement précieuses en Europe. Pour se diriger, ils utilisaient des astrolabes, des instruments de navigation essentiels pour mesurer l'altitude des étoiles et des planètes au-dessus de l'horizon afin de résoudre des problèmes de navigation.

## Votre objectif

Tout comme Magellan, vous serez des explorateurs devant naviguer des terres inconnues pour s'approprier d'un précieux trésor. Vous connaissez sa coordonnée, mais l'emplacement des obstacles reste floue. Il y a aura deux phases dans votre aventure : découverte et prise.

## Éléments de la carte

Les éléments suivants se trouveront sur la carte :   

|        |                                             |
| ------ | ------------------------------------------- |
| Pièce  | ![image](/starter-packs/docs/coin.png)      |
| Trésor | ![image](/starter-packs/docs/astrolab2.png) |
| Mur    | ...                                         |



### Grille discrète

La grille discrète contient le nombre de murs par 4 cases. Les murs extérieurs détimitant la carte sont aussi comptés.

<div align="center">
  <img width="1000" alt="logo" src="./docs/grille_murs.png">
</div>

### Pointage

| Item       | Points |
| ---------- | ------ |
| Pistolet   | 15     |
| Épée       | 40     |
| Pièce      | 40     |
| Trésor     | 1200   |

## Déroulement du jeu

Dans une partie, il y a deux tours.

### Tour 1

Durée : x temps
Les agents et les pièces sont placées de manière alétoire sur la carte 🗺️

Lorsqu'un agent prend une pièce, celle-ci reaparrait de manière aléatoire sur la carte.

Lors de la fin de la partie, les agents sont enelevés de la carte.

### Tour 2

Durée : x temps
Un trésor est placé au centre de la carte, les agents apparaitront à équidistance de déplacement du trésor. Dans ce tour, il n'y aura pas de pièces sur la carte.

### Mort

Lorsque l'agent perd toute sa vie, ce dernier disparait de la carte et ne reçoit aucune donnée du serveur pendant un temps défini.

## Actions

Plusieurs actions peuvent être envoyées au serveur dans un même cycle de raifraichissement avec certaines contraintes d'utilisation.

### Déplacement

La position envoyée est celle vers laquele l'agent va se déplacer. L'agent ne peut pas traverser les murs. Il n'y a pas d'algorithme de recherche de chemin (pathfinding) d'implémenté, vous devrez l'implémenter pour vous déplacer dans la carte. À noter, il n'y a pas de collisions entre les joueurs.

Cette action n'a pas de contraintes d'utilisation.

###  Attaque

#### Équipement

L'agent doit être équipé d'une arme pour pouvoir l'utiliser.

Cette action ne peut pas être accompagnée de l'utilisation d'une arme dans le même cycle de rafraichissement.

#### Pistolet

Pour utiliser le pistolet, il faut envoyer la position de destination souhaitée d'un projectile. 
Cette action ne peut pas être accompagnée de l'équipement d'une arme dans le même cycle de rafraichissement. 

#### Épée

Lorsque l'épée est équipée, elle apparait à 0 radians comme illustré à l'image suivante.  

<div align="center">
  <img width="200" alt="logo" src="./docs/blade2.png">
</div>


Pour changer la rotation de l'épée, il faut envoyer la nouvelle rotation en radians.

### Sauvegarde

Une quantité limitée d'octects pourra être envoyé au serveur. Ces octets seront sauvegardé le temps d'une partie. Ces données seront reçues par le joueur à chaque fois qu'il va se connecter au serveur.

## Données reçues

### État de la carte

L'état de carte est reçu lors de la connexion au serveur. 

```
MapState {
    DiscretGrid : [][]Integer,      // Grille contenant le nombre de murs par 4 cases
    size: int,                      // La taille de la carte
    walls: List[Collider],          // La liste des murs, cette liste devrait être vide
    save: bytearray                 // La mémoire que vous pouvez sauvegardé dans le jeu
}
```

### État du jeu

L'état du jeu reçu à chaque cycle de rafraichissement.

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
            rotation : Integer                      // Angle en radians
        }
    ],
    Coins : [                                       // Liste des pièces ou trésor
        uid : String,                               // Identifiant de la pièce
        value : Integer,                            // Valeur de la pièce ou du trésor
        position : { x : Float, y : Float }         // Position de la pièce ou du trésor
    ]
}
```

## Constantes

AJOUTER TABLEAU